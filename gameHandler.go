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

type gamestate int

const (
	pregame   gamestate = iota
	initRound gamestate = iota
	running   gamestate = iota
	roundDone gamestate = iota
	gameDone  gamestate = iota
)

func (gs gamestate) MarshalText() (text []byte, err error) {
	switch gs {
	case pregame:
		return []byte("pregame"), nil
	case initRound:
		return []byte("init_round"), nil
	case running:
		return []byte("running"), nil
	case roundDone:
		return []byte("round_done"), nil
	case gameDone:
		return []byte("game_done"), nil
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

	man *Managers

	gameView *gameViewer

	// Channels
	timerDeadline *time.Timer
	timerMinline  *time.Timer
	timeStamp     time.Time
	register      chan *Player
	unregister    chan *Player

	adminChan chan string
	turnDone  chan bool

	// Game info
	GameNumber  int
	RoundNumber int
	GameMap     GameMap
	baseMap     GameMap // A copy for safe keeping
	CurrentTick int

	mapLock     sync.Mutex
	playersLock sync.Mutex
}

func newGameHandler() *GameHandler {
	logger.Info("New GameHandler")
	gh := &GameHandler{
		players:     make(map[*Player]bool),
		Status:      pregame,
		register:    make(chan *Player, 5),
		unregister:  make(chan *Player, 5),
		adminChan:   make(chan string, 5),
		turnDone:    make(chan bool, 2),
		GameNumber:  0,
		RoundNumber: 0,
		config:      NewConfigHolder(),
	}

	gh.timerDeadline = time.NewTimer(gh.config.turnTimeMin)
	gh.timerMinline = time.NewTimer(gh.config.turnTimeMax)

	return gh
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
			g.playersLock.Unlock()
			go player.run()

			if g.man.am == nil {
				continue
			}

			g.man.am.pushPlayers()
		case player := <-g.unregister:
			logger.Infof("Unregistering %v", player)

			delete(g.players, player)
			err := player.conn.Close()
			g.man.am.pushPlayers()

			if err != nil {
				logger.Info("Problems closing websocket")
			}
		}
	}
}

// Runs as a go func
func (g *GameHandler) gameState() {
	for {
		switch g.Status {
		case pregame:
			g.pregame()
		case running:
			g.running()
		}
	}
}

func (g *GameHandler) pregame() {
	time.Sleep(time.Nanosecond)
}

func (g *GameHandler) running() {
	logger.Info("Starting the running subsystem")
	g.timerMinline = time.NewTimer(g.config.turnTimeMin)
	g.timerDeadline = time.NewTimer(g.config.turnTimeMax)
	g.timeStamp = time.Now()

	for {
		select {
		case <-g.timerMinline.C:
			logger.Info("Min deadline hit")
			// Check if all is done and exec if true
			if g.checkPlayersDone() {
				g.turnDone <- true
			}
		case <-g.timerDeadline.C:
			logger.Info("Max deadline hit")
			// Exec regardless if all is done
			g.turnDone <- false
		case all := <-g.turnDone:
			logger.Infof("Turn is done. All: %v", all)
			g.execTurn()
			g.newTurn()
		default:
			timeMark := time.Now().After(g.timeStamp.Add(g.config.turnTimeMin))
			// If all is done and we are after the timestamp
			if g.checkPlayersDone() && timeMark {
				g.turnDone <- true
			}
		}
	}
}

func (g *GameHandler) checkPlayersDone() bool {
	for p := range g.players {
		if p.status != CommandSent {
			return false
		}
	}
	logger.Info("Players done")

	return true
}

func (g *GameHandler) roundDone() {
	time.Sleep(time.Nanosecond)
}

func (g *GameHandler) gameDone() {
	time.Sleep(time.Nanosecond)
}

// newTurn reset states for the clients and readies everything for a new round
// It also resets the timers we need
func (g *GameHandler) newTurn() {
	logger.Info("Init new turn")

	g.pushToPlayers()

	for p := range g.players {
		if p.status == CommandWait {
			p.ticksLost += 1
		}
		p.status = CommandWait
	}

	// Sets the new timers
	g.timerDeadline.Stop()
	g.timerMinline.Stop()
	g.timerDeadline.Reset(g.config.turnTimeMax)
	g.timerMinline.Reset(g.config.turnTimeMin)
}

// Run trough all the clients in an order
// All moves happens on the same time, so two snakes can kill each other with their head
func (g *GameHandler) execTurn() {
	logger.Info("Exec turn is running")

	g.CurrentTick += 1

	g.mapLock.Lock()
	// TODO: Update map targets

	// TODO: Check for collisions
	g.mapLock.Unlock()

	if g.isRoundDone() {
		//TODO: What to do when round is done
	}
}

func (g *GameHandler) pushToPlayers() {
	//TODO: Push data to players

	for p := range g.players {
		go p.pushGameState(g)
	}

}

// Checks if ticks has been reached or there is less than two players left
func (g *GameHandler) isRoundDone() bool {
	return false
}

func (g *GameHandler) startGame() {
	// Kick all players without names
	logger.Info("Starting game")
	if g.Status != pregame {
		logger.Error("Was not in pregame state? aborting")
		return
	}

	g.playersLock.Lock()

	logger.Info("Kicking all players without name")
	for player := range g.players {
		if player.status == NoUsername {
			g.man.gm.unregister <- player
		}
	}
	g.playersLock.Unlock()

	g.initGame()
	g.initRound()

	logger.Info("Starting the rest of the system")
	g.Status = running
}

// This function will at some point pause the game
func (g *GameHandler) pauseGame() {

}

// This function will practically restart the current set of rounds
func (g *GameHandler) restartGame() {

}

// Things we need:
//	Players
//	Scores
//	TurnNewConfigHolder
//	Round

// This function sets all the values that should last for an entire game
// Should we kick all players?
func (g *GameHandler) initGame() {
	g.RoundNumber = 0
	g.Status = pregame

	if g.config.mapSizeY != 0 && g.config.mapSizeX != 0 {
		g.GameMap = baseGameMap(g.config.mapSizeX, g.config.mapSizeY, g.config.OuterWalls)
		g.baseMap = g.GameMap
		return
	}
	size := baseGameMapSize(len(g.players))
	g.GameMap = baseGameMap(size, size, g.config.OuterWalls)
	g.baseMap = g.GameMap
}

// This function sets all the values that should last for an entire round
func (g *GameHandler) initRound() {
	// TODO: finish initRound
	logger.Info("Initializing new round")
	// Reset ticks
	g.CurrentTick = 0
	g.GameMap = g.baseMap
	// TODO: Make this work properly, make it so that we can pass a proper map, not just reinint this one

	// Set pos of players
	logger.Info("Init players")
	for player := range g.players {
		logger.Infof("p: %s", player.Username)
		x, y, err := g.GameMap.findEmptySpot(false)
		if err != nil {
			panic("Could not init round, not enough space to start new round")
		}
		logger.Infof("pos: %v %v", x, y)

		player.PosX = []int{x, x, x}
		player.PosY = []int{y, y, y}
		player.Size = 3
		player.RoundScore = 0

	}

	// Set food
	logger.Infof("Setting %v foods", g.config.targetFood)
	food := 0
	for food < g.config.targetFood {
		x, y, err := g.GameMap.findEmptySpot(false)
		logger.Infof("pos: %v %v", x, y)

		if err != nil {
			panic("Could not init round, not enough space to start new round")
		}
		food += 1

		err = g.GameMap.setTile(x, y, blockFood)
		if err != nil {
			panic("Could not init round, not enough space to start new round")
		}
	}

	logger.Infof("Foods: %v", g.GameMap.Foods)

}

// Creates the status object used by the game frontend
func (g *GameHandler) generateStatusJson() []byte {

	tmpPlayers := make(map[string]Player)

	for k := range g.players {
		tmpPlayers[k.Username] = *k
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
		MapSize:       "40x40",
		OuterWalls:    1,

		turnTimeMin: 400 * time.Millisecond,
		turnTimeMax: 800 * time.Millisecond,

		GameRounds:     5,
		RoundTicks:     1000,
		targetFood:     2,    // The number of food we are trying to have on the map at once
		contWithWinner: true, // Should end game when winners is clear (for example 3/5 wins already)

	}
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

	// GameData
	// X,Y is two lists which when zipped creates the coordinates of the snake
	PosX       []int
	PosY       []int
	HeadX      int
	HeadY      int
	Size       int
	TotalScore int
	RoundScore int
}
