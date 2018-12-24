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

// GameHandler takes care of the general logic connected to running the game.
// It looks for all the commands that has been sent in from the different players.
type GameHandler struct {
	gameNumber  int
	roundNumber int
	gameMap     GameMap
}

func newGameHandler() *GameHandler {
	log.Print("New GameHandler")
	return &GameHandler{
		gameNumber:  0,
		roundNumber: 0,
		gameMap:     *baseGameMap(),
	}
}

func (g *GameHandler) run() {
	log.Print("GameHandler started")
	defer log.Panic("GameHandler Stopped")

}

// Player struct is strongly connected to the player struct.
// There should be an 1:1 ration between those entities
type Player struct {
	// Logic data
	conn *Connection
	gm   *GameHandler
	man  *Manager

	// GameData
	posX int
	posY int

	score int
}
