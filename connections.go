package main

import (
	"log"
)

// Manager is the struct for keeping order in all our connections
type Manager struct {
	clients map[*Connection]Connection

	// Registering and removeing connections
	register   chan *Connection
	unregister chan *Connection
}

func newManager() *Manager {
	log.Print("New manager")
	return &Manager{
		clients:    make(map[*Connection]Connection),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
	}
}

func (man *Manager) run() {
	log.Print("Manager started")
	defer log.Print("Manager died")
	for {
		select {
		case client := <-man.register:
			log.Printf("Registering new client %v", client)
			// TODO: Add code for keeping control of new clients
			man.clients[client] = *client

			log.Print("Starting new client")
			client.run()

		case client := <-man.unregister:
			log.Printf("Removing client %v", client)

			// TODO: Add code to remove client
			delete(man.clients, client)
		default:
			continue
		}
	}
}
