package main

import "log"

// Manager is the struct for keeping order in all our connections
type Manager struct {
	clients map[*Connection]Connection

	// Registering and removing connections
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
			man.clients[client] = *client
			log.Print("Starting new client")
			go client.run()

		case client := <-man.unregister:
			log.Print("Removing client")
			close(client.qRecv)
			close(client.qSend)
			delete(man.clients, client)
		default:
			//log.Print("Manager loop")
			//time.Sleep(1 * time.Second)
		}
	}
}
