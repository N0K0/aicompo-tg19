package main

import (
	"log"
	"time"
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

	status gamestate

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
		players: make(map[*Player]Player),
		status:  pregame,

		timerDeadline: time.NewTimer(turnTimeMax),
		timerMinline:  time.NewTimer(turnTimeMin),
		register:      make(chan *Player, 10000),
		unregister:    make(chan *Player, 10000),
		adminChan:     make(chan string, 10000),
		gameNumber:    0,
		roundNumber:   0,
		gameMap:       *baseGameMap(),
	}
}

func (g *GameHandler) run() {
	log.Print("GameHandler started")
	defer log.Panic("GameHandler Stopped")
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
	// Logic data
	conn      *Connection
	gm        *GameHandler
	man       *Manager
	turnsLost int

	// GameData
	posX  int
	posY  int
	score int
}
