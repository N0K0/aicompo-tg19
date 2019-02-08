package main

import (
	"log"
	"net/http"
	"os"
	_ "runtime/pprof"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

// Managers is a general structure to keep the other managers
type Managers struct {
	gm *GameHandler
	am *adminHandler
	vc *gameViewer
}

var upgrader = websocket.Upgrader{}

func wsConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("Error setting up new socket connection")
		return
	}

	if len(manager.gm.players) >= 16 {
		logger.Info("Unable to add more players")

		err = newSocket.Close()
		if err != nil {
			logger.Info("Unable to close new socket")
		}
		return
	}

	logger.Infof("New socket: %v", &newSocket)

	manager.gm.register <- Player{
		conn:         newSocket,
		Username:     "No Username",
		status:       NoUsername,
		Command:      "",
		qSend:        make(chan []byte, 10),
		qRecv:        make(chan []byte, 10),
		gmUnregister: manager.gm.unregister,
		turnsLost:    0,
		PosX:         make([]int, 0),
		PosY:         make([]int, 0),
		size:         0,
	}

	logger.Infof("New socket from %v", newSocket.RemoteAddr())
}

// wsAdminConnector is used by admins to control the game
func wsAdminConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("Error setting up new socket connection")
		return
	}

	// We don't remake the admin connection just because the browser got lost
	if manager.am == nil {
		logger.Info("Creating new adminmanager")
		manager.am = &adminHandler{
			gm:    manager.gm,
			conn:  newSocket,
			qRecv: make(chan []byte, 10),
			qSend: make(chan []byte, 10),
		}
		go manager.am.run()

	} else if manager.am.conn == nil {
		logger.Info("Using old manager, giving it new socket")
		manager.am.conn = newSocket
		// We only need to rerun the sockets
		go manager.am.readSocket()
		go manager.am.writeSocket()
	}
}

// wsViewConnector is an connection to get view data
func wsViewConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	newSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("Error setting up new socket connection")
		return
	}

	// We don't remake the admin connection just because the browser got lost
	if manager.vc == nil {
		logger.Info("Creating new gameViewr")
		manager.vc = &gameViewer{
			gm:    manager.gm,
			conn:  newSocket,
			qRecv: make(chan []byte, 10),
			qSend: make(chan []byte, 10),
		}

		manager.gm.gameView = manager.vc
		go manager.vc.run()

	} else if manager.vc.conn == nil {
		logger.Info("Using old manager, giving it new socket")
		manager.vc.conn = newSocket
		// We only need to rerun the sockets
		go manager.vc.readSocket()
		go manager.vc.writeSocket()
	}
}

func main() {
	logger.Init("aiCompo", true, false, os.Stdout)
	gameHandler := newGameHandler()
	go gameHandler.run()

	m := &Managers{
		gameHandler,
		nil,
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
	http.HandleFunc("/view",
		func(w http.ResponseWriter, r *http.Request) {
			wsViewConnector(m, w, r)
		},
	)

	http.Handle("/", http.FileServer(http.Dir("frontend/")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
