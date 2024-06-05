package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/f3rcho/rest-posts/handlers"
	"github.com/f3rcho/rest-posts/middleware"
	"github.com/f3rcho/rest-posts/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		JWTSecret:   JWT_SECRET,
		Port:        PORT,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal("Error creating new server", err)
	}

	s.Start(BindRoutes)
}

const (
	POST_ROUTE = "/posts/{id}"
)

func BindRoutes(s server.Server, r *mux.Router) {
	r.Use(middleware.CheckAuth(s))
	r.Use(middleware.JSONMiddleware)

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.InserPost(s)).Methods(http.MethodPost)
	r.HandleFunc(POST_ROUTE, handlers.GetPostById(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts", handlers.ListPosts(s)).Methods(http.MethodGet)
	r.HandleFunc(POST_ROUTE, handlers.UpdatePost(s)).Methods(http.MethodPut)
	r.HandleFunc(POST_ROUTE, handlers.DeletePostById(s)).Methods(http.MethodDelete)

	r.HandleFunc("/ws", s.Hub().HandleWebSocket)
}
