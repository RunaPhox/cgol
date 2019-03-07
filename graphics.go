package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func randColor(rgb chan uint8) {
	for v := 0.0; ; v += 0.002 {
		sin, cos, sin2 := math.Sin(v), math.Cos(v), math.Sin(2*v)
		rgb <- uint8(sin * sin * 255)
		rgb <- uint8(sin2 * sin2 * 255)
		rgb <- uint8(cos * cos * 255)
	}
}

func draw(r *sdl.Renderer, tab [][]byte, rgb chan uint8,
	edit *edit) {
	r.SetDrawColor(0x00, 0x00, 0x00, 0xff)
	r.Clear()
	drawPop(r, tab, rgb)
	drawGrid(r, edit)
	r.Present()
}

func drawGrid(r *sdl.Renderer, edit *edit) {
	r.SetDrawColor(0x66, 0x66, 0x66, 0xff)
	for i := int32(cellSize); i <= height-cellSize; i += cellSize {
		r.DrawLine(0, i, width, i)
	}

	for i := int32(cellSize); i <= width-cellSize; i += cellSize {
		r.DrawLine(i, 0, i, height)
	}

	r.SetDrawColor(0xf4, 0xdf, 0x42, 0xFF)

	if edit.shift {
		x1, y1, x2, y2 := sqrPoints(
			edit.lastX, edit.lastY, edit.shiftX, edit.shiftY)
		r.DrawRect(&sdl.Rect{
			X: x1 * cellSize,
			Y: y1 * cellSize,
			W: cellSize * (x2 - x1 + 1),
			H: cellSize * (y2 - y1 + 1),
		})
	} else {
		r.DrawRect(&sdl.Rect{
			X: edit.lastX * cellSize,
			Y: edit.lastY * cellSize,
			W: cellSize,
			H: cellSize,
		})
	}
}

func drawPop(r *sdl.Renderer, tab [][]byte, rgb chan uint8) {
	re, gr, bl := <-rgb, <-rgb, <-rgb
	r.SetDrawColor(re, gr, bl, 0xff)
	for i := int32(0); i < rows; i++ {
		for j := int32(0); j < columns; j++ {
			if tab[i][j] == 1 {
				r.FillRect(&sdl.Rect{
					X: j * cellSize,
					Y: i * cellSize,
					W: cellSize,
					H: cellSize,
				})
			}
		}
	}
}
