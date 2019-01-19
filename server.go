package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/gorilla/websocket"
)

// Managers is a general structure to keep the other managers
type Managers struct {
	gm *GameHandler
	cm *Manager
	am *adminHandler
}

var upgrader = websocket.Upgrader{}
var index = template.Must(template.ParseFiles("frontend/index.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s", r.RemoteAddr, r.URL)

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "frontend/index.html")
}

func wsConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	log.Print("WS started")
	defer log.Print("WS done")

	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error setting up new socket connection")
		return
	}

	log.Printf("New socket: %v", &newSocket)

	conn := &Connection{
		man:      manager.cm,
		conn:     newSocket,
		username: "No Username",
		status:   NoUsername,
		command:  "No Command",
		qSend:    make(chan []byte, 10),
		qRecv:    make(chan []byte, 10),
	}

	log.Print("Register con")
	manager.cm.register <- conn
	log.Print("Register con done")

	log.Print("Register Player")
	manager.gm.register <- &Player{
		conn:      conn,
		gm:        manager.gm,
		man:       manager.cm,
		turnsLost: 0,
		posX:      make([]int, 0),
		posY:      make([]int, 0),
		size:      0,
	}
	log.Print("Register player done")

	log.Printf("New socket from %v", newSocket.RemoteAddr())
}

// wsAdminConnector is used by admins to control the game
func wsAdminConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	log.Printf("Admin connecting...")
	defer log.Printf("Admin connected")

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
			cm:    manager.cm,
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

	connectionManager := newManager()
	gameHandler := newGameHandler()

	go connectionManager.run()
	go gameHandler.run()

	m := &Managers{
		gm: gameHandler,
		cm: connectionManager,
		am: nil,
	}

	fs := http.FileServer(http.Dir("frontend/"))
	http.Handle("/", fs)

	http.HandleFunc("/ws",
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("/ws")
			wsConnector(m, w, r)
		},
	)
	http.HandleFunc("/admin",
		func(w http.ResponseWriter, r *http.Request) {
			log.Print("/admin")
			defer log.Print("DONE /admin")
			wsAdminConnector(m, w, r)
		},
	)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
