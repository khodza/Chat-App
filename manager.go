package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	webSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	Clients ClientList
	sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Clients: make(ClientList),
	}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")
	con, err := webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}
	client := NewClient(con, m)

	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.Clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, exists := m.Clients[client]; exists {
		client.connection.Close()
		delete(m.Clients, client)
	}
}
