package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/logger"
	"github.com/gorilla/websocket"
)

const (
	// TODO: This is getting moved into GameConfigHolder
	// For which when a turn executes regardless of if all has sent their commands
	turnTimeMax = 7500 * time.Millisecond
	turnTimeMin = 2000 * time.Millisecond

	// Other config. TODO: Look into flags
	gameRounds     = 5
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
	players map[*Player]bool
	Status  gamestate
	config  *GameConfigHolder

	gameView *gameViewer

	// Channels
	timerDeadline *time.Timer
	timerMinline  *time.Timer
	register      chan *Player
	unregister    chan *Player

	adminChan chan string

	// Game info
	GameNumber  int
	RoundNumber int
	GameMap     GameMap

	mapLock     sync.Mutex
	playersLock sync.Mutex
}

func newGameHandler() *GameHandler {
	logger.Info("New GameHandler")
	return &GameHandler{
		players:       make(map[*Player]bool),
		Status:        pregame,
		timerDeadline: time.NewTimer(turnTimeMax),
		timerMinline:  time.NewTimer(turnTimeMin),
		register:      make(chan *Player, 1),
		unregister:    make(chan *Player, 1),
		adminChan:     make(chan string, 1),
		GameNumber:    0,
		RoundNumber:   0,
		GameMap:       baseGameMap(),
		config:        NewConfigHolder(),
	}
}

func (g *GameHandler) run() {
	logger.Info("GameHandler started")
	defer log.Panic("GameHandler Stopped")
	go g.gameState()
	for {
		select {
		case player := <-g.register:
			g.playersLock.Lock()
			g.players[player] = true
			logger.Infof("Players: %v", len(g.players))
			g.playersLock.Unlock()
			go player.run()
		case player := <-g.unregister:
			logger.Infof("Unregistering %v", player)
			delete(g.players, player)
			close(player.qRecv)
			close(player.qSend)
			err := player.conn.Close()
			if err != nil {
				logger.Info("Problems closing websocket")
			}
		}
	}
}

func (g *GameHandler) gameState() {
	for {
		time.Sleep(1 * time.Nanosecond)
		switch g.Status {
		case pregame:
			break
		case running:
			break
		}
	}
}

func (g *GameHandler) pregame() {
	time.Sleep(time.Nanosecond)
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
			logger.Info("GameHandler")
			time.Sleep(1 * time.Second)
		}
	}
}

func (g *GameHandler) gameDone() {
	time.Sleep(time.Nanosecond)
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
//	TurnNewConfigHolder
//	Round

func (g *GameHandler) generateStatusJson() []byte {

	tmpPlayers := make(map[string]bool)

	for k := range g.players {
		tmpPlayers[k.Username] = true
	}

	status := StatusObject{
		NumPlayers: len(tmpPlayers),
		Players:    tmpPlayers,
		GameStatus: *g,
	}

	bytes, err := json.Marshal(status)
	if err != nil {
		logger.Infof("Unable to marshal json")
		panic("Unable to marshal json")
	}
	return bytes
}

func NewConfigHolder() *GameConfigHolder {
	return &GameConfigHolder{
		MinTurnUpdate: 400,
		MaxTurnUpdate: 800,
		MapSize:       "0x0",
		OuterWalls:    1,
	}
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
	turnsLost    int
	gmUnregister chan *Player

	// GameData
	// X,Y is two lists which when zipped creates the coordinates of the snake
	PosX []int
	PosY []int
	size int
}
