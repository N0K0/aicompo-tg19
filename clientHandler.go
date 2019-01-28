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

func (p *Player) writeSocket() {
	log.Print("Write Socket")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Print("Player dead")
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-p.qSend:
			err := p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("client.conn.SetWriteDeadline")
			}
			if !ok {
				// Client closed connection.
				err := p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Print("Client closed connection")
				}
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Print("client.conn.NextWriter")
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Print("w.Write")
			}

			// Add queued chat messages to the current websocket message.
			n := len(p.qSend)
			for i := 0; i < n; i++ {
				_, err := w.Write(<-p.qSend)
				if err != nil {
					log.Print("w.Write")
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("client.conn.SetWriteDeadline")
			}
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (p *Player) readSocket() {
	log.Print("Read Socket")
	p.conn.SetReadLimit(maxMessageSize)
	err := p.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Print("client.conn.SetReadDeadline")
	}
	p.conn.SetPongHandler(
		func(string) error {
			err := p.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				log.Print("client.conn.SetReadDeadline")
			}
			return nil
		},
	)

	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("error: %v", err)
			} else {
				log.Printf("Client '%s' closed socket at %v ", p.Username, p.conn.RemoteAddr())
			}
			p.gm.unregister <- *p
			break
		}
		p.qRecv <- message
		log.Printf("Got message %s (queue: %v)", message, len(p.qRecv))
	}
}

func (p *Player) parseCommand() {
	for {
		select {
		case incoming, ok := <-p.qRecv:
			log.Printf("Queue: %v", len(p.qRecv))

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
				p.setUsername(&command)
				break
			case "move":
				break
			default:
				log.Print("Player sent invalid command")
				p.sendError("Invalid command type!")
				break
			}
		}
	}
}

func (p *Player) sendError(message string) {
	msg := ClientInfo{
		Type:    "error",
		Message: message,
	}

	jsonString, err := json.Marshal(msg)
	if err != nil {
		log.Print("Problems with creating error message")
	}

	p.qSend <- jsonString
}

func (p *Player) setUsername(cmd *Command) {
	username := cmd.Value

	if len(username) > 14 {
		p.sendError("Username is too long! Max length 14")
		return
	}

	if p.status != NoUsername {
		p.sendError("Username already set!")
		return
	}

	p.Username = username
	p.status = ReadyToPlay
	log.Printf("Player given name: %v", username)
}

func (p *Player) run() {

	// Routines used to interact with the WebSocket
	// I'm keeping is separated from the logic below to make everything a bit more clean
	go p.writeSocket()
	go p.readSocket()
	go p.parseCommand()
}
