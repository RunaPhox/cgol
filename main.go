package main

import (
	"time"

	"github.com/runaphox/cgol/conway"
)

func main() {
	var wrap sdlContext
	var err error
	var w, h, c int32 = 1920, 1080, 20

	game := stage{
		pause:    true,
		wrap:     true,
		width:    w,
		height:   h,
		cellSize: c,
		columns:  w / c,
		rows:     h / c,
		timeEx:   60,
	}
	game.tab = game.newTab()

	wrap.w, wrap.r, err = initSdl(&game)
	if err != nil {
		panic(err)
	}
	defer closeSdl(wrap)

	edit := edit{}

	rgb := make(chan uint8, 3)
	go randColor(rgb)
	go handleEvents(wrap.w, &game, &edit)
	go renderSimulation(wrap.r, &game, rgb, &edit)

	for !game.quit {
		start := time.Now()
		if !game.pause {
			game.tab = conway.Update(game.tab, game.wrap)
		}
		time.Sleep(start.Sub(time.Now()) + time.Duration(game.timeEx)*time.Millisecond)
	}
}
