package main

import (
	"fmt"
	"hagerbot.com/rog"
	"io/ioutil"
	"strings"
	"time"
)

type Point struct {
	x int
	y int
}

var (
	dirLeft     = Point{-1, 0}
	dirRight    = Point{1, 0}
	dirUp       = Point{0, -1}
	dirDown     = Point{0, 1}
	dirSameSpot = Point{0, 0}
)

func (level *Level) doneMovingBoulders() bool {
	for x := 0; x < level.width; x++ {
		for y := level.height - 1; y >= 0; y-- {
			tile := level.tiles[y][x]
			switch tile {
			case Boulder:
				if level.boulderMoved(x, y) {
					return false
				}
			}
		}
	}
	return true
}

func (level *Level) boulderMoved(x int, y int) bool {
	south := &level.tiles[y+1][x]
	east := &level.tiles[y][x+1]
	west := &level.tiles[y][x-1]
	se := &level.tiles[y+1][x+1]
	sw := &level.tiles[y+1][x-1]
	north := &level.tiles[y-1][x]
	boulderPtr := &level.tiles[y][x]

	if *south == Space {
		*south = Boulder
		*boulderPtr = Space
		return true
	} else if (*south == Boulder || *south == Diamond) &&
		*east == Space &&
		*se == Space {
		*boulderPtr = Space
		*se = Boulder
		return true
	} else if (*south == Boulder || *south == Diamond) &&
		*west == Space &&
		*sw == Space {
		*boulderPtr = Space
		*sw = Boulder
		return true
	} else if *south == Player && *north == Space {
		level.gameOver = true
		return false
	}
	return false
}

func (level *Level) movePlayer(d Point) {
	tilePtr := &level.tiles[level.player.y+d.y][level.player.x+d.x]
	playerPtr := &level.tiles[level.player.y][level.player.x]
	switch *tilePtr {
	case Ground, Space:
		*playerPtr = Space
		*tilePtr = Player
		level.player.x += d.x
		level.player.y += d.y
	case Diamond:
		*playerPtr = Space
		*tilePtr = Player
		level.player.x += d.x
		level.player.y += d.y
		level.points++
	}
}

type Level struct {
	tiles       [][]rune
	width       int
	height      int
	player      Point
	pointsToGet int
	points      int
	gameOver    bool
}

const (
	Player  = 'X'
	Wall    = '#'
	Ground  = '-'
	Boulder = 'O'
	Space   = ' '
	Diamond = '@'
)

func parseLevel(levelStr string) (level Level) {
	lines := strings.Split(levelStr, "\n")
	var foundPlayer = false
	for y, line := range lines {
		level.tiles = append(level.tiles, []rune(line))
		for x, tile := range line {
			switch tile {
			case Player:
				level.player.x = x
				level.player.y = y
				foundPlayer = true
			case Diamond:
				level.pointsToGet += 1
			}
		}
	}
	if !foundPlayer {
		panic("No player in this level!")
	}
	level.width = len(level.tiles[0])
	level.height = len(level.tiles)
	level.gameOver = false
	return level
}

func (level *Level) getString() (output string) {
	for _, line := range level.tiles {
		output += string(line) + "\n"
	}
	return output
}

func getDirection(key int) Point {
	switch key {
	case 63235:
		return dirRight
	case 63234:
		return dirLeft
	case 63233:
		return dirDown
	case 63232:
		return dirUp
	default:
		return dirSameSpot
	}
	return dirSameSpot
}

func main() {
	levelBytes, err := ioutil.ReadFile("level")
	if err != nil {
		panic("error loading level")
	}
	level := parseLevel(string(levelBytes))
	rog.Open(level.width+1, level.height+1, 1, false, "HackDash", nil)
	for rog.Running() {
		levelStr := level.getString()
		rog.Set(0, 0, nil, nil, levelStr)
		if level.gameOver {
			time.Sleep(time.Second)
			fmt.Println("Game over, man! GAME OVER!")
			level = parseLevel(string(levelBytes))
		}
		key := rog.Key()
		if key == rog.Esc {
			rog.Close()
		}
		if !level.doneMovingBoulders() {
			time.Sleep(25 * time.Millisecond)
		} else {
			d := getDirection(key)
			if d != dirSameSpot {
				level.movePlayer(d)
			}
			time.Sleep(50 * time.Millisecond)
			if level.pointsToGet == level.points {
				fmt.Printf("You won with %d points!\n", level.points)
				rog.Close()
			}
		}
		rog.Flush()
	}
}
