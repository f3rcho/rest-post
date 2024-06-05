package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub      *Hub
	id       string
	socket   *websocket.Conn
	outbound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

// func (c *Client) Write() {
// 	for {
// 		select {
// 		case message, ok := <-c.outbound:
// 			if !ok {
// 				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
// 			}
// 			c.socket.WriteMessage(websocket.TextMessage, message)
// 		}
// 	}
// }

func (c *Client) Write() {
	for message := range c.outbound {
		c.socket.WriteMessage(websocket.TextMessage, message)
	}
	// If the channel is closed, send a close message
	c.socket.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c Client) Close() {
	c.socket.Close()
	close(c.outbound)
}
