package main

import (
	"fmt"
	"math/rand"
	"time"

	term "github.com/nsf/termbox-go"
)

var gameMap [20][20]string
var snake Snake

// symbols
const (
	right = "ᐳ"
	left  = "ᐸ"
	up    = "ᐱ"
	down  = "ᐯ"
	coin  = "C"
)

// colors for symbols
const (
	green  = "\033[32m"
	yellow = "\033[33m"
	reset  = "\033[97m"
)

type Coord struct {
	x, y int
}

type Snake struct {
	pos    []Coord
	size   int
	symbol string
}

func (s Snake) head() Coord {
	return s.pos[len(s.pos)-1]
}

func (s Snake) move(c Coord) {
	gameMap[c.x][c.y] = green + s.symbol + reset
	snake.pos = append(snake.pos, c)

	if len(snake.pos) > snake.size {
		gameMap[snake.pos[0].x][snake.pos[0].y] = " "
		snake.pos = snake.pos[1:]
	}
}

func (s Snake) collision(c Coord) bool {
	for i := 0; i < len(snake.pos); i++ {
		if snake.pos[i] == c {
			return true
		}
	}
	return false
}

func resetGame() {
	snake = Snake{[]Coord{{0, 0}, {0, 1}}, 2, right}

	for i := 0; i < len(gameMap); i++ {
		for j := 0; j < len(gameMap[i]); j++ {
			temp := " "
			coord := Coord{i, j}
			for _, pos := range snake.pos {
				if pos == coord {
					temp = green + right + reset
					break
				}
			}
			gameMap[i][j] = temp
		}
	}

}

func showMap() {
	for i := 0; i < len(gameMap); i++ {
		fmt.Println(gameMap[i])
	}
}

func generateCoin() Coord {
	return Coord{rand.Intn(len(gameMap)), rand.Intn(len(gameMap[0]))}
}

func putCoin(c Coord) {
	gameMap[c.x][c.y] = yellow + coin + reset
}

func main() {
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	resetGame()

	key := "right"
	move := Coord{0, 1}

	currentCoin := generateCoin()
	putCoin(currentCoin)

Game:
	for range time.Tick(time.Second / 4) {

		go func(key *string) {
			switch ev := term.PollEvent(); ev.Type {
			case term.EventKey:
				switch ev.Key {
				case term.KeyEsc:
					*key = "exit"
				case term.KeyArrowUp:
					if *key != "down" {
						*key = "up"
					}
				case term.KeyArrowDown:
					if *key != "up" {
						*key = "down"
					}
				case term.KeyArrowLeft:
					if *key != "right" {
						*key = "left"
					}
				case term.KeyArrowRight:
					if *key != "left" {
						*key = "right"
					}
				}
			case term.EventError:
				panic(ev.Err)
			}
		}(&key)

		if key == "exit" {
			break Game
		}

		move = snake.head()

		switch key {
		case "right":
			move.y++
			if move.y >= len(gameMap) {
				move.y = 0
			}
			snake.symbol = right
		case "down":
			move.x++
			if move.x >= len(gameMap) {
				move.x = 0
			}
			snake.symbol = down
		case "left":
			move.y--
			if move.y < 0 {
				move.y = len(gameMap) - 1
			}
			snake.symbol = left
		case "up":
			move.x--
			if move.x < 0 {
				move.x = len(gameMap) - 1
			}
			snake.symbol = up
		}

		snake.move(move)

		if snake.head() == currentCoin {
			currentCoin = Coord{-1, -1}
			snake.size++

			for {
				tempCoin := generateCoin()

				if !snake.collision(tempCoin) {
					currentCoin = tempCoin
					putCoin(currentCoin)
					break
				}
			}
		}

		term.Sync()
		showMap()
	}
}
