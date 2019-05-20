package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/google/logger"
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
	gameDone  chan bool

	// Game info
	RoundNumber int
	TotalRounds int
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
		gameDone:    make(chan bool, 2),
		RoundNumber: 0,
		config:      NewConfigHolder(),
	}

	gh.timerDeadline = time.NewTimer(gh.config.turnTimeMin)
	gh.timerMinline = time.NewTimer(gh.config.turnTimeMax)
	gh.TotalRounds = gh.config.GameRounds
	return gh
}

func (g *GameHandler) run() {
	logger.Info("GameHandler started")
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
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

			g.man.am.pushState()
		case player := <-g.unregister:
			g.playersLock.Lock()

			logger.Infof("Unregistering %v", player)

			delete(g.players, player)
			err := player.conn.Close()
			g.playersLock.Unlock()
			g.man.am.pushState()

			if err != nil {
				logger.Info("Problems closing websocket")
			}
		case <-ticker.C:
			if g.man.am == nil {
				continue
			}
			g.man.am.pushState()
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

	defer g.timerMinline.Stop()
	defer g.timerDeadline.Stop()

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
			g.man.am.pushState()
			g.gameView.qSend <- g.generateStatusJson()
			g.newTurn()
			g.timeStamp = time.Now()
		case <-g.gameDone:
			logger.Info("Game done!")
			g.Status = gameDone
			g.man.am.gameDone()
			return
		default:
			timeMark := time.Now().After(g.timeStamp.Add(g.config.turnTimeMin))
			// If all is done and we are after the timestamp
			if g.checkPlayersDone() && timeMark {
				logger.Info("Default turn is done triggered")
				g.turnDone <- true
			}
		}
	}
}

func (g *GameHandler) checkPlayersDone() bool {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	for p := range g.players {
		if p.command == "" && p.status != Dead {
			return false
		}
	}
	return true
}

// newTurn reset states for the clients and readies everything for a new round
// It also resets the timers we need
func (g *GameHandler) newTurn() {
	logger.Info("Init new turn")

	g.pushToPlayers()

	for p := range g.players {
		if p.status == Dead {
			p.Head = coord{-1, -1}
			p.Tail = coord{-1, -1}
			p.next = coord{-1, -1}
			continue
		}

		p.status = CommandWait
		p.command = ""
	}

	// Sets the new timers
	g.timerDeadline.Stop()
	g.timerMinline.Stop()
	g.timerDeadline.Reset(g.config.turnTimeMax)
	g.timerMinline.Reset(g.config.turnTimeMin)
}

func (g *GameHandler) setRandomMove() string {
	moves := []string{"left", "right", "top", "bottom"}
	move := moves[rand.Intn(len(moves))]
	logger.Infof("Moving %v", move)
	return move
}

// Run trough all the clients in an order
// All moves happens on the same time, so two snakes can kill each other with their head
func (g *GameHandler) execTurn() {
	logger.Info("Exec turn is running")

	g.CurrentTick += 1

	g.mapLock.Lock()
	// TODO: Update map targets

	for p := range g.players {
		if p.status == Dead {
			continue
		}

		if p.command == "" {
			logger.Infof("%v no move set. Settings random move", p.Username)
			p.command = g.setRandomMove()
			p.ticksLost += 1
		}
		p.status = CommandSent
		p.setMove()
	}

	var nextCoords []coord // Kills you
	var tailCoords []coord // Will not kill you
	for p := range g.players {
		if p.status == Dead {
			continue
		}
		nextCoords = append(nextCoords, p.next)
		tailCoords = append(tailCoords, p.Tail)
	}

	// update Locations to the next node, remove tail if size says so
	for p := range g.players {
		logger.Infof("U: %v, n: %v, s: %v", p.Username, p.next, p.status)
		if p.status != CommandSent {
			continue
		}

		block := p.checkSnakeCollision(nextCoords, tailCoords, g)
		r, _ := block.toRune()
		logger.Infof("Collision with %v", strconv.QuoteRuneToASCII(r))

		if block == blockSnake || block == blockWall {
			p.die()

			// Grant everyone else that is alive a point
			bonus := 0
			switch g.playersLeft() {
			case 3:
				bonus = 1
			case 2:
				bonus = 1
			case 1:
				bonus = 3
			}
			g.grantBonusPoints(bonus)

			continue
		} else if block == blockFood {
			logger.Infof("%v hit food", p.Username)
			g.GameMap.removeFood(p)
			g.GameMap.spreadFood(g.config.targetFood)
		}

		p.makeMove(block, g.GameMap)
	}

	// Removes all snakes from the maps
	posX, posY, _ := g.GameMap.getAllEmpty(blockSnake)

	for i := range posX {
		_ = g.GameMap.setTile(posX[i], posY[i], blockClear)
	}

	// Regenerates snakes on map
	for p := range g.players {
		for i := 0; i < p.Size; i++ {
			_ = g.GameMap.setTile(p.PosX[i], p.PosY[i], blockSnake)
		}
	}

	g.mapLock.Unlock()

	if g.isRoundDone() {
		//TODO: What to do when round is done
		logger.Info("Round is done trigger")
		for p := range g.players {
			p.TotalScore += p.RoundScore
			p.RoundScore = 0
		}

		if g.RoundNumber >= g.config.GameRounds {
			logger.Info("Game over! Swapping to winner screen")
			g.Status = pregame
			g.gameDone <- true
		}

		g.initRound()
	}
}

func (g *GameHandler) pushToPlayers() {

	for p := range g.players {
		go p.pushGameState(g)
	}
}

func (g *GameHandler) playersLeft() int {
	left := 0

	for p := range g.players {
		if p.status != Dead {
			left++
		}
	}
	return left
}

func (g *GameHandler) grantBonusPoints(points int) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	for p := range g.players {
		if p.status != Dead {
			p.RoundScore += points
		}
	}
}

// Checks if ticks has been reached or there is less than two players left
func (g *GameHandler) isRoundDone() bool {
	logger.Info("Checking if round is done")
	totalPlayers := len(g.players)
	livePlayers := totalPlayers
	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	for p := range g.players {
		if p.status == Dead || p.status == Disconnected {
			livePlayers -= 1
		}
	}

	if livePlayers <= 1 {
		logger.Info("One player left standing")
		return true
	}

	if g.CurrentTick >= g.config.RoundTicks {
		logger.Info("Round out of time")
		return true
	}

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

	for p := range g.players {
		p.status = ReadyToPlay
		p.RoundScore = 0
		p.TotalScore = 0

	}

}

// This function sets all the values that should last for an entire round
func (g *GameHandler) initRound() {
	// TODO: finish initRound
	logger.Info("Initializing new round")
	// Reset ticks
	g.CurrentTick = 0
	g.RoundNumber += 1
	g.GameMap = *g.setupGameMap()
	// TODO: Make this work properly, make it so that we can pass a proper map, not just reinint this one

	// Set pos of players
	logger.Info("Init players")
	for player := range g.players {
		logger.Infof("p: %s", player.Username)
		x, y, err := g.GameMap.findEmptySpot()
		if err != nil {
			panic("Could not init round, not enough space to start new round")
		}
		logger.Infof("pos: %v %v", x, y)
		player.status = ReadyToPlay

		player.PosX = []int{x, x, x}
		player.PosY = []int{y, y, y}
		player.Head = coord{x, y}
		player.Tail = coord{x, y}

		player.Size = 3
		player.RoundScore = 0

	}

	// Set food
	logger.Infof("Setting %v foods", g.config.targetFood)
	g.GameMap.spreadFood(g.config.targetFood)
	logger.Infof("Foods: %v", g.GameMap.Foods)

}

func (g *GameHandler) setupGameMap() *GameMap {
	if g.config.mapSizeY != 0 && g.config.mapSizeX != 0 {
		gm := baseGameMap(g.config.mapSizeX, g.config.mapSizeY, g.config.OuterWalls)
		return &gm
	}
	size := baseGameMapSize(len(g.players))
	gm := baseGameMap(size, size, g.config.OuterWalls)
	return &gm
}

func (g *GameHandler) generateStatusObject() *StatusObject {
	tmpPlayers := make(map[string]Player)

	g.playersLock.Lock()
	defer g.playersLock.Unlock()
	for k := range g.players {
		tmpPlayers[k.Username] = *k
	}

	return &StatusObject{
		NumPlayers: len(tmpPlayers),
		Players:    tmpPlayers,
		GameStatus: *g,
	}
}

// Creates the status object used by the game frontend
func (g *GameHandler) generateStatusJson() []byte {

	status := g.generateStatusObject()

	bytes, err := json.Marshal(status)
	if err != nil {
		logger.Infof("Unable to marshal json")
		panic("Unable to marshal json")
	}
	bytes = append(bytes, byte('\n'))
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

		GameRounds:     10,
		RoundTicks:     1000,
		targetFood:     2,    // The number of food we are trying to have on the map at once
		contWithWinner: true, // Should end game when winners is clear (for example 3/5 wins already)

	}
}
