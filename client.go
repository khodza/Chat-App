package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	Manager    *Manager
	egress     chan []byte
}

func NewClient(connection *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: connection,
		Manager:    manager,
		egress:     make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.Manager.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error reading Message: ", err.Error())
			}
			break
		}
		for wsclient := range c.Manager.Clients {
			wsclient.egress <- payload
		}
		log.Println(messageType)
		log.Println(string(payload))
	}

}

func (c *Client) writeMessages() {
	defer func() {
		c.Manager.removeClient(c)
	}()
	for {
		select {
		case messages, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Error closing connection: ", err.Error())
				}
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, messages); err != nil {
				log.Println("Error writing message: ", err.Error())
			}
		}
	}
}
