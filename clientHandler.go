package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/google/logger"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 3 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 3 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Status for enumeration
type Status int

// Status enumeration types
const (
	NoUsername Status = iota
	ReadyToPlay
	CommandWait
	CommandSent
	Dead
	Disconnected
)

func (p *Player) writePump() {
	logger.Info("Write Socket")
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Infof("Player %v not answering pings. Removing from game", p.Username)
		pingTicker.Stop()
		p.gmUnregister <- p
	}()
	for {
		select {
		case message, ok := <-p.qSend:
			err := p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Info("client.conn.SetWriteDeadline")
			}
			if !ok {
				// Client closed connection.
				p.connLock.Lock()
				err := p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.Info("Client closed connection")
				}
				p.connLock.Unlock()
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				logger.Info("w.Write")
			}

			// Add queued chat messages to the current websocket message.
			n := len(p.qSend)
			for i := 0; i < n; i++ {
				_, err := w.Write(<-p.qSend)
				if err != nil {
					logger.Info("w.Write")
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-pingTicker.C:
			//logger.Info("ping ticker")
			err := p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Info("Error set deadline Ticker")
			}
			p.connLock.Lock()
			err = p.conn.WriteMessage(websocket.PingMessage, nil)
			p.connLock.Unlock()

			if err != nil {
				return
			}
		}
	}
}

func (p *Player) readPump() {
	logger.Info("Read Socket")
	err := p.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logger.Info("client.conn.SetReadDeadline")
	}
	p.conn.SetPongHandler(
		func(string) error {
			err := p.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				logger.Info("client.conn.SetReadDeadline")
			}
			return nil
		},
	)

	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				logger.Infof("error: %v", err)
			} else {
				logger.Infof("Client '%s' closed socket at %v ", p.Username, p.conn.RemoteAddr())
			}
			p.gmUnregister <- p
			break
		}
		//logger.Infof("Got message from player: %v", message)
		p.qRecv <- message
	}
}

// This function pushes the gamestatus object to the players
// Contains about the same as the one we push to the frontend, but with an extra field for that player
func (p *Player) pushGameState(g *GameHandler) {

	pso := PlayerStatusObject{
		StatusObject: g.generateStatusObject(),
		Self:         p,
	}

	bytes, err := json.Marshal(pso)
	if err != nil {
		logger.Infof("Unable to marshal json")
		panic("Unable to marshal json")
	}

	p.qSend <- bytes
}

func (p *Player) parseCommand() {
	for {
		select {
		case incoming, ok := <-p.qRecv:
			if !ok {
				logger.Infof("Closed socket")
				return
			}

			command := Command{}
			err := json.Unmarshal(incoming, &command)
			if err != nil {
				logger.Infof("Invalid json: %v", err)
			}
			//logger.Infof("Json: %v", command)

			switch command.Type {

			case "username":
				p.setUsername(&command)
			case "color":
				p.setColor(&command)
			case "move":
				p.parseMoveCommand(&command)
			default:
				logger.Info("Player sent invalid command")
				p.sendError("Invalid command type!")
				break
			}
		}
	}
}
func (p *Player) parseMoveCommand(c *Command) {
	if p.status == Dead {
		return
	}

	logger.Infof("Move %v: %v", p.Username, c.Value)

	switch c.Value {
	case "left":
		p.command = "left"
	case "right":
		p.command = "right"
	case "up":
		p.command = "up"
	case "down":
		p.command = "down"
	default:
		p.command = ""
		logger.Info("Move command invalid!")
		p.sendError("Move command invalid!")
	}
}

func (p *Player) sendError(message string) {
	msg := Envelope{
		Type:    "error",
		Message: message,
	}

	jsonString, err := json.Marshal(msg)
	if err != nil {
		logger.Info("Problems with creating error message")
	}

	p.qSend <- jsonString
}

func (p *Player) setUsername(cmd *Command) {
	username := cmd.Value
	logger.Infof("Setting username %v", username)

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
	logger.Infof("Player given name: %v", username)
}

// Takes in an arbitrary string. Is passed to the frontend to be used with the ctx.fillStyle property
func (p *Player) setColor(cmd *Command) {
	color := cmd.Value
	p.Color = color
	logger.Infof("Setting color for %v '%v'", p.Username, p.Color)
}

func (p *Player) setMove() {
	p.next = coord{p.Head.X, p.Head.Y}
	switch p.command {
	case "left":
		p.next.X -= 1
	case "right":
		p.next.X += 1
	case "up":
		p.next.Y += 1
	case "down":
		p.next.Y -= 1
	}
}

func (p *Player) makeMove(b block, gm GameMap) {
	logger.Infof("Player makeMove: %v", p.Username)

	// Found fruit in move
	if b == blockFood {
		logger.Info("Got food")
		p.Size += 1
		p.RoundScore += 1
	}

	p.Head = coord{p.next.X, p.next.Y}

	p.PosX = append([]int{p.Head.X}, p.PosX...)
	p.PosY = append([]int{p.Head.Y}, p.PosY...)

	if s := len(p.PosX); s > p.Size {
		p.PosX = p.PosX[:s-1]
		p.PosY = p.PosY[:s-1]
		_ = gm.setTile(p.Tail.X, p.Tail.Y, blockClear)
	}

	if p.Size == 0 {
		p.Tail = coord{-1, -1}
	} else {
		p.Tail = coord{p.PosX[p.Size-1], p.PosY[p.Size-1]}
	}

	logger.Infof("Head: %v", p.Head)
	logger.Infof("Tail: %v", p.Tail)
	logger.Infof("Body: %v", p.PosX)
	logger.Infof("Body: %v", p.PosY)

}

// The player has dun goofed. The player dies
func (p *Player) die() {
	logger.Infof("Killing player %v", p.Username)
	p.status = Dead
	p.PosX = []int{}
	p.PosY = []int{}
	p.Head = coord{-1, -1}
	p.Tail = coord{-1, -1}
	p.next = coord{-1, -1}

	p.Size = 0
}

// Checks all head placements
// Returns blockSnake or blockWall if it should die. Else it returns blockFood or blockClear
func (p *Player) checkSnakeCollision(nextCoords []coord, tailCoords []coord, g *GameHandler) block {
	// its a fatal collision if one of the following is hit:
	// The next head of a snake
	// A wall
	// A part of the snake that is not a tail

	nX := p.next.X
	nY := p.next.Y
	nextBlock, err := g.GameMap.getTile(nX, nY)
	if err != nil {
		logger.Fatal("Found coord outside of map! Wrapping is not implemented yet!")
	}

	nextHits := 0
	for player := range g.players {
		if p.next.X == player.next.X && p.next.Y == player.next.Y {
			nextHits += 1
		}
	}

	// We need to check the next blocks first, since multiple player may hit the same food at the same time
	index, err := g.GameMap.findCoordInList(&p.next, nextCoords)
	if index >= 0 && nextHits > 1 { // This is the same as an headon collision next round
		logger.Infof("%v headon collision", p.Username)
		return blockSnake
	}

	// Kills you
	if nextBlock == blockWall {
		logger.Infof("%v wall collision", p.Username)
		return blockWall
	} else if nextBlock == blockSnake {
		_, err = g.GameMap.findCoordInList(&p.next, tailCoords)
		if index != -1 { // Hit the wall of an snake, and not the tail
			logger.Infof("%v snake collision", p.Username)
			return blockSnake
		}
		logger.Infof("%v taking over over tail spot", p.Username)

	} else if nextBlock == blockFood {
		return blockFood
	} else if nextBlock == blockClear {
		return blockClear
	}

	logger.Fatal("Its not possible to get here")
	return blockClear
}

func (p *Player) run() {

	// Routines used to interact with the WebSocket
	// I'm keeping is separated from the logic below to make everything a bit more clean
	go p.writePump()
	go p.readPump()
	go p.parseCommand()
}

// Player struct is strongly connected to the player struct.
// There should be an 1:1 ration between those entities
type Player struct {
	// The websocket to the client
	conn     *websocket.Conn
	connLock sync.Mutex

	// The Username of the client
	Username string `json:"username"`

	// Color in a format that
	Color string

	// The status of the connection
	status Status

	// Last command read
	command string

	// Channels for caching data
	qSend chan []byte
	qRecv chan []byte
	// Logic data
	ticksLost    int
	gmUnregister chan *Player
	sync.Mutex
	// GameData
	// X,Y is two lists which when zipped creates the coordinates of the snake
	PosX       []int // List of all nodes that is in the snake
	PosY       []int
	Head       coord // Head of the snake
	Tail       coord // Backmost block of the snake
	next       coord // Next movetarget of the snake
	Size       int   // Current size
	TotalScore int   // Gamescore
	RoundScore int   //Score this round
}
