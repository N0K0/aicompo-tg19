package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/google/logger"
)

type block int

const (
	blockClear     = 0    // _
	blockWall      = iota // X
	blockSnake     = iota // *
	blockSnakeHead = iota // H
	blockFood      = iota // ^
)

const (
	blockClearChar     = '_'
	blockWallChar      = 'X'
	blockSnakeChar     = '*'
	blockSnakeHeadChar = 'H'
	blockFoodChar      = '^'
)

/*
GameMap is a small struct for maps passed as a "2d string"
*/
type GameMap struct {

	// The "simple map"
	SizeX   int
	SizeY   int
	Content [][]block

	// Some lists if you would rather use them for the parsing
	Heads []Head
	Walls []Wall // This contains also the rest of the snake blocks
	Foods []Food
}

func (gm *GameMap) getTile(x int, y int) (block, error) {
	if x < 0 || x >= gm.SizeX {
		return blockClear, errors.New("X out of bounds")
	}

	if y < 0 || y >= gm.SizeY {
		return blockClear, errors.New("Y out of bounds")
	}

	return gm.Content[y][x], nil
}

func (gm *GameMap) setTileLine(startX int, startY int, deltaX int, deltaY int, value block, iterations int) {
	x := startX
	y := startY
	for iter := 0; iter < iterations; iter++ {
		err := gm.setTile(x, y, value)
		if err != nil {
			logger.Info("Hit outside bound, stopping line")
			return
		}
		x += deltaX
		y += deltaY
	}
}

func (gm *GameMap) setTile(x int, y int, value block) error {
	if x < 0 || x >= gm.SizeX {
		return errors.New("X out of bounds")
	}

	if y < 0 || y >= gm.SizeY {
		return errors.New("Y out of bounds")
	}

	switch value {
	case blockWall:
		fallthrough
	case blockSnake:
		gm.Walls = append(gm.Walls, Wall{x, y})
	case blockSnakeHead:
		gm.Heads = append(gm.Heads, Head{x, y})
	case blockFood:
		gm.Foods = append(gm.Foods, Food{x, y})
	}

	gm.Content[y][x] = value
	return nil
}

func (b block) MarshalText() (text []byte, err error) {
	val, err := b.toRune()

	if err != nil {
		return []byte(""), err
	}
	return []byte(strconv.QuoteRuneToASCII(val)), nil
}

// Returns an rune from an block. return space and an error if not possible
func (b *block) toRune() (rune, error) {
	switch *b {
	case blockClear:
		return blockClearChar, nil
	case blockWall:
		return blockWallChar, nil
	case blockSnake:
		return blockSnakeChar, nil
	case blockSnakeHead:
		return blockSnakeHeadChar, nil
	case blockFood:
		return blockFoodChar, nil
	}
	return rune(-1), errors.New("unable to convert block to rune")
}

// Creates a block from a rune. return an error and -1 if not possible
func toBlock(r rune) (block, error) {
	switch r {
	case blockClearChar:
		return blockClear, nil
	case blockWallChar:
		return blockWall, nil
	case blockSnakeChar:
		return blockSnake, nil
	case blockSnakeHeadChar:
		return blockSnakeHead, nil
	case blockFoodChar:
		return blockFood, nil
	}
	return -1, errors.New("unable to convert rune to block")
}

func (gm *GameMap) getAllEmpty() ([]int, []int, error) {
	listX := make([]int, 0)
	listY := make([]int, 0)

	for indexY := range gm.Content {
		yLine := gm.Content[indexY]
		for indexX, xBlock := range yLine {
			if xBlock != blockClear {
				continue
			}
			listX = append(listX, indexX)
			listY = append(listY, indexY)
		}
	}

	if len(listX) == 0 {
		return nil, nil, errors.New("no empty tiles left")
	}

	return listX, listY, nil
}

//
// Finds an empty spot which can be used for placing down objects
// The fair bool value exists for trying to place the objects some part away from the snakes head
// Note that fair does not do anything yet!
// Returns the x,y coordinates for a empty spot
func (gm *GameMap) findEmptySpot(fair bool) (int, int, error) {
	if fair {
		logger.Warning("the param fair is not implemented yet!")
	}

	listX, listY, err := gm.getAllEmpty()

	if err != nil {
		return -1, -1, err
	}

	if len(listX) != len(listY) {
		return -1, -1, errors.New("invalid length on the two lists")
	}

	element := rand.Intn(len(listX))
	return listX[element], listY[element], nil
}

func baseGameMap(sizeX int, sizeY int, walls int) GameMap {
	logger.Infof("Generating map with size: %v,%v", sizeX, sizeY)

	blankMap := fmt.Sprintf("%v,%v\n", sizeX, sizeY)
	for i := 1; i < sizeY; i++ {
		blankMap = blankMap + strings.Repeat("_", sizeX) + "\n"
	}

	gm := mapFromString(blankMap)

	if walls == 0 {
		return gm
	}

	// Top
	gm.setTileLine(0, 0, 1, 0, blockWall, sizeX)

	// Bottom
	gm.setTileLine(0, sizeY-1, 1, 0, blockWall, sizeX)

	// Left
	gm.setTileLine(0, 0, 0, 1, blockWall, sizeY)

	// Right
	gm.setTileLine(sizeX-1, 0, 0, 1, blockWall, sizeY)

	return gm
}

func baseGameMapSize(numberPlayers int) int {
	return numberPlayers*2 + 10
}

/* mapFromString takes in an map in the form of x,y and then y lines with x length denoting the map.
The map is denoted with the chars:

_ : May walk on
X : Wall. Blocked
* : Fuel
^ : Bullet

Note that it is also implicit that out of bound are walls.

Example map:

---Star---
6,6
X____X
__XX__
_XXXX_
______
XX_XX_
XX____
---End---

*/
func mapFromString(mapInput string) GameMap {
	logger.Info("Parsing map")
	defer logger.Info("Map parsed")
	lines := strings.Split(mapInput, "\n")
	size := strings.Split(lines[0], ",")
	sizeX, err := strconv.Atoi(size[0])

	if err != nil {
		log.Fatal("Got invalid map for the X")
	}

	sizeY, err := strconv.Atoi(size[1])
	if err != nil {
		log.Fatal("Got invalid map for the Y")
	}

	if sizeY != len(lines)-1 {
		log.Fatalf("Mismatch between size Y of the map and the number given. SizeY: %v len(lines): %v", sizeY, len(lines))
	}

	if len(lines[1]) != sizeX {
		log.Fatal("Mismatch between size of X and the size of the first line of the map")
	}

	gm := GameMap{
		Content: nil,
	}

	content := make([][]block, sizeY)

	for index, line := range lines[1:] {
		contentLine := make([]block, sizeX)

		for index, char := range line {
			switch char {
			case blockClearChar:
				contentLine[index] = blockClear
			case blockWallChar:
				contentLine[index] = blockWall
			case blockSnakeChar:
				contentLine[index] = blockSnake
			case blockSnakeHead:
				contentLine[index] = blockSnakeHead
			case blockFoodChar:
				contentLine[index] = blockFood
			default:
				log.Panicf("Found invalid char: '%c'", char)
			}
		}

		content[index] = contentLine
	}
	gm.Content = content
	gm.SizeX = sizeX
	gm.SizeY = sizeY
	return gm
}

func (m *GameMap) spreadFood(targetAmount int) {

}
