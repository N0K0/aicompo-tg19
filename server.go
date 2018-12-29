package main

import (
	"flag"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Managers is a general structure to keep the other managers
type Managers struct {
	gm *GameHandler
	cm *Manager
}

var listenAddr = flag.String("listen_addr", "0.0.0.0:8080", "Add listening socket and port")
var upgrader = websocket.Upgrader{}
var index = template.Must(template.ParseFiles("frontend/index.html"))

func generatePassword() string {
	passwordLen := 5

	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	str := make([]rune, passwordLen)
	for i := range str {
		str[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(str)
}

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
	conn := &Connection{
		man:      manager.cm,
		conn:     newSocket,
		username: "No Username",
		status:   NoUsername,
		command:  "No Command",
		qSend:    make(chan []byte, 10000),
		qRecv:    make(chan []byte, 10000),
	}

	manager.cm.register <- conn

	manager.gm.register <- &Player{
		conn:      conn,
		gm:        manager.gm,
		man:       manager.cm,
		turnsLost: 0,
		posX:      0,
		posY:      0,
		score:     0,
	}

	log.Printf("New socket from %v", newSocket.RemoteAddr())
}

// wsAdminConnector is used by admins to control the game
func wsAdminConnector(manager *Managers, password string, w http.ResponseWriter, r *http.Request) {
	log.Printf("Admin connecting...")
	defer log.Printf("Admin connected")

	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error setting up new socket connection")
		return
	}

	adm := &adminHandler{
		gm:   manager.gm,
		cm:   manager.cm,
		conn: newSocket,

		password: password,
	}

	go adm.run()
}

func main() {
	flag.Parse()
	log.Printf("Starting server at %s", *listenAddr)

	adminPassword := generatePassword()
	log.Print("******************************************************")
	log.Print("******************************************************")
	log.Printf("ADMIN PASSWORD: %s", adminPassword)
	log.Print("******************************************************")
	log.Print("******************************************************")

	connectionManager := newManager()
	go connectionManager.run()

	gameHandler := newGameHandler()
	go gameHandler.run()

	m := &Managers{
		gm: gameHandler,
		cm: connectionManager,
	}

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws",
		func(w http.ResponseWriter, r *http.Request) {
			wsConnector(m, w, r)
		},
	)
	http.HandleFunc("/admin",
		func(w http.ResponseWriter, r *http.Request) {
			wsConnector(m, w, r)
		},
	)

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
