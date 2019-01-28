package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type blockstatus int

const (
	blockClear = iota // _
	blockWall  = iota // X
	blockSnake = iota // *
	blockFood  = iota // ^
)

const (
	blockClearChar = '_'
	blockWallChar  = 'X'
	blockSnakeChar = '*'
	blockFoodChar  = '^'
)

/*
GameMap is a small struct for maps
*/
type GameMap struct {
	SizeX   int
	SizeY   int
	Content [][]int
}

func baseGameMap() GameMap {
	defaultSize := 80
	blankMap := fmt.Sprintf("%v,%v\n", defaultSize, defaultSize)
	for i := 1; i < defaultSize; i++ {
		blankMap = blankMap + strings.Repeat("_", defaultSize) + "\n"
	}
	return mapFromString(blankMap)
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
	log.Print("Parsing map")
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
		SizeX:   0,
		SizeY:   0,
		Content: nil,
	}

	content := make([][]int, sizeY)

	for index, line := range lines[1:] {
		contentLine := make([]int, sizeX)

		for index, char := range line {
			switch char {
			case blockClearChar:
				contentLine[index] = blockClear
				break
			case blockWallChar:
				contentLine[index] = blockWall
				break
			case blockSnakeChar:
				contentLine[index] = blockSnake
				break
			case blockFoodChar:
				contentLine[index] = blockFood
				break
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
