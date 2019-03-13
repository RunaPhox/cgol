package main

import (
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func renderSimulation(r *sdl.Renderer, game *stage, rgb chan uint8,
	edit *edit) {
	for {
		start := time.Now()
		draw(r, game, rgb, edit)
		time.Sleep(start.Sub(time.Now()) + 32*time.Millisecond)
	}
}

func randColor(rgb chan uint8) {
	for v := 0.0; ; v += 0.002 {
		sin, cos, sin2 := math.Sin(v), math.Cos(v), math.Sin(2*v)
		rgb <- uint8(sin * sin * 255)
		rgb <- uint8(sin2 * sin2 * 255)
		rgb <- uint8(cos * cos * 255)
	}
}

func draw(r *sdl.Renderer, game *stage, rgb chan uint8,
	edit *edit) {
	r.SetDrawColor(0x00, 0x00, 0x00, 0xff)
	r.Clear()
	drawPop(r, game, rgb)
	drawGrid(r, edit, game)
	r.Present()
}

func drawGrid(r *sdl.Renderer, edit *edit, game *stage) {
	r.SetDrawColor(0x66, 0x66, 0x66, 0xff)
	for i := int32(game.cellSize); i <= game.height-game.cellSize; i += game.cellSize {
		r.DrawLine(0, i, game.width, i)
	}

	for i := int32(game.cellSize); i <= game.width-game.cellSize; i += game.cellSize {
		r.DrawLine(i, 0, i, game.height)
	}

	r.SetDrawColor(0xf4, 0xdf, 0x42, 0xFF)

	if edit.ctrl {
		x1, y1, x2, y2 := sqrPoints(
			edit.lastX, edit.lastY, edit.ctrlX, edit.ctrlY)
		if x2-x1 <= y2-y1 {
			r.DrawRect(&sdl.Rect{
				X: edit.ctrlX * game.cellSize,
				Y: y1 * game.cellSize,
				W: game.cellSize,
				H: game.cellSize * (y2 - y1 + 1),
			})
		} else {
			r.DrawRect(&sdl.Rect{
				X: x1 * game.cellSize,
				Y: edit.ctrlY * game.cellSize,
				W: game.cellSize * (x2 - x1 + 1),
				H: game.cellSize,
			})
		}
	} else if edit.shift {
		x1, y1, x2, y2 := sqrPoints(
			edit.lastX, edit.lastY, edit.shiftX, edit.shiftY)
		r.DrawRect(&sdl.Rect{
			X: x1 * game.cellSize,
			Y: y1 * game.cellSize,
			W: game.cellSize * (x2 - x1 + 1),
			H: game.cellSize * (y2 - y1 + 1),
		})
	} else {
		r.DrawRect(&sdl.Rect{
			X: edit.lastX * game.cellSize,
			Y: edit.lastY * game.cellSize,
			W: game.cellSize,
			H: game.cellSize,
		})
	}
}

func drawPop(r *sdl.Renderer, game *stage, rgb chan uint8) {
	re, gr, bl := <-rgb, <-rgb, <-rgb
	r.SetDrawColor(re, gr, bl, 0xff)
	for i := int32(0); i < game.rows; i++ {
		for j := int32(0); j < game.columns; j++ {
			if game.tab[i][j] == 1 {
				r.FillRect(&sdl.Rect{
					X: j * game.cellSize,
					Y: i * game.cellSize,
					W: game.cellSize,
					H: game.cellSize,
				})
			}
		}
	}
}
