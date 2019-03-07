package main

import (
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
