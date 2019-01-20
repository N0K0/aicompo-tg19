package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Managers is a general structure to keep the other managers
type Managers struct {
	gm *GameHandler
	am *adminHandler
}

var upgrader = websocket.Upgrader{}

func wsConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	log.Print("WS started")
	defer log.Print("WS done")

	if len(manager.gm.players) >= 16 {
		text := []byte(`{error: "No more spots"}`)
		w.Write(text)
		return
	}

	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error setting up new socket connection")
		return
	}

	log.Printf("New socket: %v", &newSocket)
	log.Print("Register Player")

	manager.gm.register <- &Player{
		conn:      newSocket,
		username:  "No username",
		status:    0,
		command:   "",
		qSend:     make(chan []byte, 10),
		qRecv:     make(chan []byte, 10),
		gm:        manager.gm,
		turnsLost: 0,
		posX:      make([]int, 0),
		posY:      make([]int, 0),
		size:      0,
	}

	log.Printf("New socket from %v", newSocket.RemoteAddr())
}

// wsAdminConnector is used by admins to control the game
func wsAdminConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error setting up new socket connection")
		return
	}

	// We don't remake the admin connection just because the browser got lost
	if manager.am == nil {
		log.Print("Creating new adminmanager")
		manager.am = &adminHandler{
			gm:    manager.gm,
			conn:  newSocket,
			qRecv: make(chan []byte, 10),
			qSend: make(chan []byte, 10),
		}
		go manager.am.run()

	} else if manager.am.conn == nil {
		log.Print("Using old manager, giving it new socket")
		manager.am.conn = newSocket
		// We only need to rerun the sockets
		go manager.am.readSocket()
		go manager.am.writeSocket()
	}
}

func main() {

	gameHandler := newGameHandler()
	go gameHandler.run()

	m := &Managers{
		gameHandler,
		nil,
	}

	http.HandleFunc("/ws",
		func(w http.ResponseWriter, r *http.Request) {
			wsConnector(m, w, r)
		},
	)
	http.HandleFunc("/admin",
		func(w http.ResponseWriter, r *http.Request) {
			wsAdminConnector(m, w, r)
		},
	)

	http.Handle("/", http.FileServer(http.Dir("frontend/")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
