package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/f3rcho/rest-posts/database"
	"github.com/f3rcho/rest-posts/repository"
	"github.com/f3rcho/rest-posts/websocket"
	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWTSecret is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("DatabaseUrl is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	binder(b, b.router)
	address := ":" + b.Config().Port

	repo, err := database.NewPostGresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	repository.SetRespository(repo)
	go b.hub.Run()
	log.Println("Starting server on port:", b.config.Port)
	if err := http.ListenAndServe(address, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
