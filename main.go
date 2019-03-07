package main

import (
	"math"
	"time"

	"github.com/runaphox/cgol/conway"
)

const (
	width    = 1920
	height   = 1080
	cellSize = 20
	columns  = width / cellSize
	rows     = height / cellSize
	popul    = columns * rows
)

type stage struct {
	quit  bool
	pause bool
	tab   [][]byte
}

func tabIndex(x, y int32) (int32, int32) {
	yInd := y / cellSize
	xInd := x / cellSize
	return xInd, yInd
}

func sqrPoints(x1, y1, x2, y2 int32) (int32, int32, int32, int32) {
	minX := int32(math.Min(float64(x1), float64(x2)))
	minY := int32(math.Min(float64(y1), float64(y2)))
	maxX := int32(math.Max(float64(x1), float64(x2)))
	maxY := int32(math.Max(float64(y1), float64(y2)))
	return minX, minY, maxX, maxY
}

func newTab(row, col int) [][]byte {
	tab := make([][]byte, rows)
	for i := range tab {
		tab[i] = make([]byte, columns)
	}

	return tab
}

func main() {
	var wrap sdlContext
	var err error

	wrap.w, wrap.r, err = initSdl()
	if err != nil {
		panic(err)
	}
	defer closeSdl(wrap)

	game := stage{pause: true, tab: newTab(rows, columns)}
	edit := edit{}

	rgb := make(chan uint8, 3)
	go randColor(rgb)
	go handleEvents(wrap.w, &game, &edit)

	for !game.quit {
		start := time.Now()

		draw(wrap.r, game.tab, rgb, &edit)
		if !game.pause {
			game.tab = conway.Update(game.tab)
		}

		time.Sleep(start.Sub(time.Now()) + 32*time.Millisecond)
	}
}
