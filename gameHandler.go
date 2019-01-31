package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// For which when a turn executes regardless of if all has sent their commands
	turnTimeMax = 750 * time.Millisecond
	turnTimeMin = 200 * time.Millisecond

	// Other config. TODO: Look into flags
	gameRounds = 5

	contWithWinner = true // Should end game when winners is clear (for example 3/5 wins already)

)

type gamestate int

const (
	pregame gamestate = iota
	running gamestate = iota
	done    gamestate = iota
)

func (gs gamestate) MarshalText() (text []byte, err error) {
	switch gs {
	case pregame:
		return []byte("pregame"), nil
	case running:
		return []byte("running"), nil
	case done:
		return []byte("done"), nil
	default:
		return []byte(""), errors.New("game state invalid")
	}
}

// GameHandler takes care of the general logic connected to running the game.
// It looks for all the commands that has been sent in from the different players.
type GameHandler struct {
	// Meta info
	players map[*Player]Player
	Status  gamestate

	gameView *gameViewer

	// Channels
	timerDeadline *time.Timer
	timerMinline  *time.Timer
	register      chan Player
	unregister    chan Player

	adminChan chan string

	// Game info
	GameNumber  int
	RoundNumber int
	GameMap     GameMap

	mapLock sync.Mutex
}

func newGameHandler() *GameHandler {
	log.Print("New GameHandler")
	return &GameHandler{
		players:       make(map[*Player]Player),
		Status:        pregame,
		timerDeadline: time.NewTimer(turnTimeMax),
		timerMinline:  time.NewTimer(turnTimeMin),
		register:      make(chan Player, 10),
		unregister:    make(chan Player, 10),
		adminChan:     make(chan string, 10),
		GameNumber:    0,
		RoundNumber:   0,
		GameMap:       baseGameMap(),
	}
}

func (g *GameHandler) run() {
	log.Print("GameHandler started")
	defer log.Panic("GameHandler Stopped")
	go g.gameState()
	for {
		select {
		case player := <-g.register:
			g.players[&player] = player
			log.Printf("Players: %v", len(g.players))
			go player.run()
		case player := <-g.unregister:
			log.Printf("Unregistering %v", player)
			delete(g.players, &player)
			close(player.qRecv)
			close(player.qSend)
			err := player.conn.Close()
			if err != nil {
				log.Print("Problems closing websocket")
			}
		}
	}
}

func (g *GameHandler) gameState() {
	for {
		switch g.Status {
		case pregame:
			break
		case running:
			break
		case done:
			break
		}
	}
}

func (g *GameHandler) pregame() {

}

func (g *GameHandler) running() {
	for {
		select {
		case <-g.timerMinline.C:
			// Check if all is done and exec if true
			break
		case <-g.timerDeadline.C:
			// Exec regardless if all is done
			break
		default:
			log.Print("GameHandler")
			time.Sleep(1 * time.Second)
		}
	}
}

func (g *GameHandler) gameDone() {

}

// newTurn reset states for the clients and readies everything for a new round
// It also resets the timers we need
func (g *GameHandler) newTurn() {
	// Sets the new timers
	g.timerDeadline.Reset(turnTimeMax)
	g.timerMinline.Reset(turnTimeMin)
}

// execTurn resets
func (g *GameHandler) execTurn() {
	// Run trough all the clients in an order (which?)
	g.mapLock.Lock()
	g.mapLock.Unlock()
}

// Things we need:
//	Players
//	Scores
//	Turn
//	Round

func (g *GameHandler) generateStatusJson() []byte {
	log.Printf("Generating state")
	defer log.Printf("Generating Done")
	tmpPlayers := make(map[string]Player)

	for k := range g.players {
		tmpPlayers[k.Username] = *k
	}

	status := StatusObject{
		NumPlayers: len(tmpPlayers),
		Players:    tmpPlayers,
		GameStatus: *g,
	}

	log.Printf("Marshal")
	bytes, err := json.Marshal(status)
	if err != nil {
		log.Printf("Unable to marshal json")
		panic("Unable to marshal json")
	}
	log.Printf("Marshal Done")

	return bytes

}

// Player struct is strongly connected to the player struct.
// There should be an 1:1 ration between those entities
type Player struct {
	// The websocket to the client
	conn *websocket.Conn

	// The Username of the client
	Username string `json:"username"`

	// The status of the connection
	status Status

	// Last command read
	Command string `json:"command"`

	// Channels for caching data
	qSend chan []byte
	qRecv chan []byte
	// Logic data
	gm        *GameHandler
	turnsLost int

	// GameData
	// X,Y is two lists which when zipped creates the coordinates of the snake
	PosX []int
	PosY []int
	size int
}
