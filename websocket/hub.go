package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not open websocket connection, %v\n", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	client := NewClient(hub, socket)
	hub.register <- client

	go client.Write()
}

func (hub *Hub) onConnect(client *Client) {
	log.Printf("CLient connected: %s\n", client.socket.RemoteAddr())

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	client.id = client.socket.RemoteAddr().String()
	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Printf("CLient disconnected: %s\n", client.socket.RemoteAddr())
	client.socket.Close()

	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for index, c := range hub.clients {
		if c.id == client.id {
			hub.clients = append(hub.clients[:index], hub.clients[index+1:]...)
			break
		}
	}
}

func (hub *Hub) BroadCast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}
