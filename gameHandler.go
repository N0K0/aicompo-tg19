package main

import (
	"log"
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

// GameHandler takes care of the general logic connected to running the game.
// It looks for all the commands that has been sent in from the different players.
type GameHandler struct {
	// Meta info
	players map[*Player]Player
	status  gamestate

	// Channels
	timerDeadline *time.Timer
	timerMinline  *time.Timer
	register      chan *Player
	unregister    chan *Player

	adminChan chan string

	// Game info
	gameNumber  int
	roundNumber int
	gameMap     GameMap
}

func newGameHandler() *GameHandler {
	log.Print("New GameHandler")
	return &GameHandler{
		players:       make(map[*Player]Player),
		status:        pregame,
		timerDeadline: time.NewTimer(turnTimeMax),
		timerMinline:  time.NewTimer(turnTimeMin),
		register:      make(chan *Player, 10),
		unregister:    make(chan *Player, 10),
		adminChan:     make(chan string, 10),
		gameNumber:    0,
		roundNumber:   0,
		gameMap:       *baseGameMap(),
	}
}

func (g *GameHandler) run() {
	log.Print("GameHandler started")
	defer log.Panic("GameHandler Stopped")
	go g.gameState()
	for {
		select {
		case player := <-g.register:
			g.players[player] = *player
			log.Printf("Players: %v", len(g.players))
			go player.run()
		case player := <-g.unregister:
			log.Printf("Unregistering %v", player)
			delete(g.players, player)
			close(player.qRecv)
			close(player.qSend)
		}
	}
}

func (g *GameHandler) gameState() {
	for {
		switch g.status {
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

}

// Player struct is strongly connected to the player struct.
// There should be an 1:1 ration between those entities
type Player struct {
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
	// Logic data
	gm        *GameHandler
	turnsLost int

	// GameData
	// X,Y is two lists which when zipped creates the coordinates of the snake
	posX []int
	posY []int
	size int
}
