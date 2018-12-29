package main

import (
	"log"
	"strconv"
	"strings"
)

type blockstatus int

const (
	blockClear = iota // _
	blockWall  = iota // X
	blockFuel  = iota // *
	blockAmmo  = iota // ^
)

const (
	blockClearChar = '_'
	blockWallChar  = 'X'
	blockFuelChar  = '*'
	blockAmmoChar  = '^'
)

/*
GameMap is a small struct for maps
*/
type GameMap struct {
	sizeX   int
	sizeY   int
	content [][]int
}

func baseGameMap() *GameMap {
	return mapFromString(baseMap)
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
func mapFromString(mapInput string) *GameMap {
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

	log.Printf("sizeY: %v, mapSize: %v", sizeY, len(lines))
	if sizeY != len(lines)-1 {
		log.Fatal("Mismatch betwee size Y of the map and the number given")
	}

	log.Printf("sizeX: %v, mapSize: %v", sizeX, len(lines[1]))
	if len(lines[1]) != sizeX {
		log.Fatal("Mistmatch between size of X and the size of the first line of the map")
	}

	gm := &GameMap{
		sizeX:   0,
		sizeY:   0,
		content: nil,
	}

	content := make([][]int, sizeY)

	for index, line := range lines[1:] {
		log.Print(line)
		contentLine := make([]int, sizeX)

		for index, char := range line {
			switch char {
			case blockClearChar:
				contentLine[index] = blockClear
				break
			case blockWallChar:
				contentLine[index] = blockWall
				break
			case blockAmmoChar:
				contentLine[index] = blockAmmo
				break
			case blockFuelChar:
				contentLine[index] = blockFuel
				break
			default:
				log.Panicf("Found invalid char: '%c'", char)
			}
		}
		log.Print(contentLine)

		content[index] = contentLine
	}
	return gm

}

const baseMap = `6,6
X____X
__XX__
_XXXX_
______
XX_XX_
XX____`
