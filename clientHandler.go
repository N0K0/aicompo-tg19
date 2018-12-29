package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Status for enumeration
type Status int

// Status enumeration types
const (
	Waiting        Status = iota
	AckCommand     Status = iota
	InvalidCommand Status = iota
	TimedOut       Status = iota
	NoUsername     Status = iota
)

// Connection is a struct for managing connections
type Connection struct {
	// The connection manager
	man *Manager

	// The websocket to the client
	conn *websocket.Conn

	// The username of the client
	username string

	// The status of the connection
	status Status

	// Last command read
	command string

	// Channels for caching data
	qSend chan []byte
	qRecv chan []byte
}

func (client *Connection) writeSocket() {
	log.Print("Write Socket")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-client.qSend:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print("Error with NextWriter")
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(client.qSend)
			for i := 0; i < n; i++ {
				w.Write(<-client.qSend)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Connection) readSocket() {
	log.Print("Read Socket")
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(
		func(string) error {
			client.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		},
	)

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			} else {
				log.Printf("Client '%s' closed socket at %v ", client.username, client.conn.RemoteAddr())
			}
			break
		}
		client.qRecv <- message
		log.Printf("Got message %s (queue: %v)", message, len(client.qRecv))
	}
}

func (client *Connection) run() {
	log.Print("Init new client loop")

	// Routines used to interact with the WebSocket
	// I'm keeping is separated from the logic below to make everything a bit more clean
	go client.writeSocket()
	go client.readSocket()
}
