package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var listenAddr = flag.String("listen_addr", "0.0.0.0:8080", "Add listening socket and port")
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

func wsConnector(manager *Manager, w http.ResponseWriter, r *http.Request) {
	log.Print("WS started")
	defer log.Print("WS done")

	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error setting up new socket connection")
		return
	}

	manager.register <- &Connection{
		man:      manager,
		conn:     newSocket,
		username: "No Username",
		status:   NoUsername,
		command:  "No Command",
		qSend:    make(chan []byte, 10000),
		qRecv:    make(chan []byte, 10000),
	}

	log.Printf("New socket from %v", newSocket.RemoteAddr())
}

func main() {
	flag.Parse()
	log.Printf("Starting server at %s", *listenAddr)

	connectionManager := newManager()
	go connectionManager.run()

	gameHandler := newGameHandler()
	go gameHandler.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws",
		func(w http.ResponseWriter, r *http.Request) {
			wsConnector(connectionManager, w, r)
		},
	)

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
