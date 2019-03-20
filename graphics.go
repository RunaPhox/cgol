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
		time.Sleep(start.Sub(time.Now()) + 16*time.Millisecond)
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

func plotLineLow(x0, y0, x1, y1 int32, r *sdl.Renderer, game *stage) {
	dx := x1 - x0
	dy := y1 - y0
	yi := int32(1)
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	D := 2*dy - dx
	y := y0

	for x := x0; x <= x1; x++ {
		plotSq(x, y, r, game)
		if D > 0 {
			y = y + yi
			D = D - 2*dx
		}
		D = D + 2*dy
	}
}

func plotLineHigh(x0, y0, x1, y1 int32, r *sdl.Renderer, game *stage) {
	dx := x1 - x0
	dy := y1 - y0
	xi := int32(1)
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	D := 2*dx - dy
	x := x0

	for y := y0; y <= y1; y++ {
		plotSq(x, y, r, game)
		if D > 0 {
			x = x + xi
			D = D - 2*dy
		}
		D = D + 2*dx
	}
}

func plotLine(x0, y0, x1, y1 int32, r *sdl.Renderer, game *stage) {
	if math.Abs(float64(y1-y0)) < math.Abs(float64(x1-x0)) {
		if x0 > x1 {
			plotLineLow(x1, y1, x0, y0, r, game)
		} else {
			plotLineLow(x0, y0, x1, y1, r, game)
		}
	} else {
		if y0 > y1 {
			plotLineHigh(x1, y1, x0, y0, r, game)
		} else {
			plotLineHigh(x0, y0, x1, y1, r, game)
		}
	}
}

func plotSq(x, y int32, r *sdl.Renderer, game *stage) {
	r.DrawRect(&sdl.Rect{
		X: x * game.cellSize,
		Y: y * game.cellSize,
		W: game.cellSize,
		H: game.cellSize,
	})
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

		plotLine(edit.lastP.x, edit.lastP.y, edit.ctrlP.x, edit.ctrlP.y, r, game)
		/* Straight lines vertical and horizontal
		x1, y1, x2, y2 := sqrPoints(
			edit.lastP.x, edit.lastP.y, edit.ctrlP.x, edit.ctrlP.y)
		if x2-x1 <= y2-y1 {
			r.DrawRect(&sdl.Rect{
				X: edit.ctrlP.x * game.cellSize,
				Y: y1 * game.cellSize,
				W: game.cellSize,
				H: game.cellSize * (y2 - y1 + 1),
			})
		} else {
			r.DrawRect(&sdl.Rect{
				X: x1 * game.cellSize,
				Y: edit.ctrlP.y * game.cellSize,
				W: game.cellSize * (x2 - x1 + 1),
				H: game.cellSize,
			})
		}
		*/
	} else if edit.shift {
		x1, y1, x2, y2 := sqrPoints(
			edit.lastP.x, edit.lastP.y, edit.shiftP.x, edit.shiftP.y)
		r.DrawRect(&sdl.Rect{
			X: x1 * game.cellSize,
			Y: y1 * game.cellSize,
			W: game.cellSize * (x2 - x1 + 1),
			H: game.cellSize * (y2 - y1 + 1),
		})
	} else {
		r.DrawRect(&sdl.Rect{
			X: edit.lastP.x * game.cellSize,
			Y: edit.lastP.y * game.cellSize,
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
