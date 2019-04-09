package main

import (
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

var wsUpgrader = websocket.Upgrader{}
var adminUpgrader = websocket.Upgrader{}
var viewUpgrader = websocket.Upgrader{}

func wsConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("Error setting up new socket connection")
		return
	}

	/*	if len(manager.gm.players) >= 16 {
		logger.Info("Unable to add more players")

		err := ws.Close()
		if err != nil {
			logger.Info("Unable to close new socket")
		}
		return
	}*/

	manager.gm.register <- &Player{
		conn:         ws,
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

	logger.Infof("New socket from %v", ws.RemoteAddr())
}

// wsAdminConnector is used by admins to control the game
func wsAdminConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {

	ws, err := adminUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("Error setting up new socket connection")
		return
	}

	// We don't remake the admin connection just because the browser got lost
	if manager.am == nil {
		logger.Info("Creating new adminmanager")
		manager.am = &adminHandler{
			gm:        manager.gm,
			man:       manager,
			conn:      ws,
			closeChan: make(chan bool, 1),
			qRecv:     make(chan []byte, 10),
			qSend:     make(chan []byte, 10),
		}
		go manager.am.run()

	} else if manager.am.conn == nil {
		// Got an working manager, just need a new connection to the frontend
		logger.Info("Using old admin-manager, giving it new socket")
		manager.am.conn = ws
		// We only need to rerun the sockets
		go manager.am.readSocket()
		go manager.am.writeSocket()
	} else {
		logger.Info("Something wierd with admin socket")
	}
}

// wsViewConnector is an connection to get view data
func wsViewConnector(manager *Managers, w http.ResponseWriter, r *http.Request) {
	ws, err := viewUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("Error setting up new socket connection")
		return
	}

	// We don't remake the admin connection just because the browser got lost
	if manager.vc == nil {
		logger.Info("Creating new gameViewer")
		manager.vc = &gameViewer{
			conn:         ws,
			gm:           manager.gm,
			qSend:        make(chan []byte, 10),
			qRecv:        make(chan []byte, 10),
			statusTicker: nil,
			pingTicker:   nil,
		}

		manager.gm.gameView = manager.vc
		go manager.vc.run()

	} else if manager.vc.conn == nil {
		logger.Info("Using old view-manager, giving it new socket")
		manager.vc.conn = ws
		// We only need to rerun the sockets
		go manager.vc.readPump()
		go manager.vc.writePump()
		go manager.vc.statusUpdater()
	}
}

func main() {
	// newMem := profile.MemProfileRate(512)
	// defer profile.Start(newMem).Stop()

	logger.Init("aiCompo", false, false, os.Stderr)
	gameHandler := newGameHandler()

	m := &Managers{
		gameHandler,
		nil,
		nil,
	}

	gameHandler.man = m

	go gameHandler.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsConnector(m, w, r)
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		wsAdminConnector(m, w, r)
	})

	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		wsViewConnector(m, w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("frontend/")))

	logger.Fatal(http.ListenAndServe(":8080", nil))
}
