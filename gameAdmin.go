package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// adminHandler takes care
type adminHandler struct {
	gm   *GameHandler
	cm   *Manager
	conn *websocket.Conn

	qSend chan []byte
	qRecv chan []byte

	password string
}

func (admin *adminHandler) run() {
	log.Printf("Admin loop started")
	log.Printf("Game state: %v", admin.gm.status)

	for {
		select {
		case incoming := <-admin.qRecv:
			log.Printf("New admin message: %s", incoming)
			break
		}
	}
}

func (admin *adminHandler) writeSocket() {
	log.Print("Write Socket")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
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
		case <-ticker.C:
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			} else {
			}
			break
		}
		admin.qRecv <- message
		log.Printf("Got admin message %s (queue: %v)", message, len(admin.qRecv))
	}
}
