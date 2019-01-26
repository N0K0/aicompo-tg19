package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 20 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 120 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 5120
)

// Status for enumeration
type Status int

// Status enumeration types
const (
	NoUsername  Status = iota
	ReadyToPlay Status = iota
	Waiting     Status = iota
)

func (player *Player) writeSocket() {
	log.Print("Write Socket")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Print("Player dead")
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-player.qSend:
			err := player.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("client.conn.SetWriteDeadline")
			}
			if !ok {
				// Client closed connection.
				err := player.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Print("Client closed connection")
				}
				return
			}

			w, err := player.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print("client.conn.NextWriter")
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Print("w.Write")
			}

			// Add queued chat messages to the current websocket message.
			n := len(player.qSend)
			for i := 0; i < n; i++ {
				_, err := w.Write(<-player.qSend)
				if err != nil {
					log.Print("w.Write")
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := player.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("client.conn.SetWriteDeadline")
			}
			if err := player.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (player *Player) readSocket() {
	log.Print("Read Socket")
	player.conn.SetReadLimit(maxMessageSize)
	err := player.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Print("client.conn.SetReadDeadline")
	}
	player.conn.SetPongHandler(
		func(string) error {
			err := player.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				log.Print("client.conn.SetReadDeadline")
			}
			return nil
		},
	)

	for {
		_, message, err := player.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("error: %v", err)
			} else {
				log.Printf("Client '%s' closed socket at %v ", player.username, player.conn.RemoteAddr())
			}
			player.gm.unregister <- player
			break
		}
		player.qRecv <- message
		log.Printf("Got message %s (queue: %v)", message, len(player.qRecv))
	}
}

func (player *Player) parseCommand() {
	for {
		select {
		case incoming, ok := <-player.qRecv:
			log.Printf("Queue: %v", len(player.qRecv))

			if !ok {
				log.Printf("Closed socket")
				return
			}

			command := Command{}
			err := json.Unmarshal(incoming, &command)
			if err != nil {
				log.Printf("Invalid json: %v", err)
			}
			log.Printf("Json: %v", command)

			switch command.Type {

			case "username":
				player.setUsername(&command)
				break
			case "move":
				break
			default:
				log.Print("Player sent invalid command")
				player.sendError("Invalid command type!")
				break
			}
		}
	}
}

func (player *Player) sendError(message string) {
	msg := ClientInfo{
		Type:    "error",
		Message: message,
	}

	jsonString, err := json.Marshal(msg)
	if err != nil {
		log.Print("Problems with creating error message")
	}

	player.qSend <- jsonString
}

func (player *Player) setUsername(cmd *Command) {
	username := cmd.Value

	if len(username) > 14 {
		player.sendError("Username is too long! Max length 14")
		return
	}

	if player.status != NoUsername {
		player.sendError("Username already set!")
		return
	}

	player.username = username
	player.status = ReadyToPlay
	log.Printf("Player given name: %v", username)
}

func (player *Player) run() {

	// Routines used to interact with the WebSocket
	// I'm keeping is separated from the logic below to make everything a bit more clean
	go player.writeSocket()
	go player.readSocket()
	go player.parseCommand()
}
