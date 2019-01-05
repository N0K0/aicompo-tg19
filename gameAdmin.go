package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// adminHandler takes care of starting the game and is used for spectating the events. Used by the webinterface
type adminHandler struct {
	gm   *GameHandler
	cm   *Manager
	conn *websocket.Conn

	qSend chan []byte
	qRecv chan []byte

	ticker *time.Ticker
}

func (admin *adminHandler) run() {
	log.Printf("Admin loop started")
	log.Printf("Game state: %v", admin.gm.status)

	go admin.readSocket()
	go admin.writeSocket()

	logger := time.NewTicker(3 * time.Second)
	for {
		select {
		case incoming := <-admin.qRecv:
			log.Printf("New admin message: %s", incoming)
			break
		case <-logger.C:
			admin.logStatus()
			break
		default:
			log.Print("Admin loop")
			time.Sleep(1 * time.Second)
		}
	}
}

func (admin *adminHandler) logStatus() {
	statusString := `Adminstatus:
		Connected Players: %v
		Connections:       %v
	`

	numPlayers := len(admin.gm.players)
	numConn := len(admin.cm.clients)

	log.Printf(statusString, numPlayers, numConn)
}

func (admin *adminHandler) writeSocket() {
	log.Print("Write Socket")
	admin.ticker = time.NewTicker(pingPeriod)
	defer func() {
		admin.ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-admin.qSend:
			admin.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				admin.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := admin.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print("Error with NextWriter")
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(admin.qSend)
			for i := 0; i < n; i++ {
				w.Write(<-admin.qSend)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-admin.ticker.C:
			admin.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := admin.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (admin *adminHandler) readSocket() {
	log.Print("Admin Read Socket")
	admin.conn.SetReadLimit(maxMessageSize)
	admin.conn.SetReadDeadline(time.Now().Add(pongWait))
	admin.conn.SetPongHandler(
		func(string) error {
			admin.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)

	for {
		_, message, err := admin.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("error: %v", err)
			} else {
				// Admin disconnected
				log.Printf("Admin closed socket at %v ", admin.conn.RemoteAddr())
				admin.conn = nil
				admin.ticker.Stop()
			}
			break
		}
		log.Print("Waiting for admin message")
		admin.qRecv <- message
		log.Printf("Got admin message %s (queue: %v)", message, len(admin.qRecv))
	}
}
