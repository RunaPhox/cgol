package main

import (
	"fmt"
	"time"

	"github.com/runaphox/cgol/conway"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width    = 1920
	height   = 1080
	cellSize = 40
	columns  = width / cellSize
	rows     = height / cellSize
	popul    = columns * rows
)

func drawGrid(r *sdl.Renderer) {
	r.SetDrawColor(0x66, 0x66, 0x66, 0xff)
	for i := int32(cellSize); i <= height-cellSize; i += cellSize {
		r.DrawLine(0, i, width, i)
	}

	for i := int32(cellSize); i <= width-cellSize; i += cellSize {
		r.DrawLine(i, 0, i, height)
	}
}

func drawPop(r *sdl.Renderer, tab [][]byte) {
	r.SetDrawColor(0x00, 0x35, 0xdb, 0xff)
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

func handleEvents(quit *bool) {
	for !*quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch ev.(type) {
			case *sdl.QuitEvent:
				*quit = true
			}
		}
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	w, err := sdl.CreateWindow("title", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, width, height,
		sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN)

	if err != nil {
		panic(err)
	}
	defer w.Destroy()

	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer r.Destroy()

	tab := make([][]byte, rows)
	for i := range tab {
		tab[i] = make([]byte, columns)
	}

	tab[5][3] = 1
	tab[4][9] = 1
	tab[4][10] = 1
	tab[4][10] = 1
	tab[4][11] = 1

	quit := false
	go handleEvents(&quit)
	for i := 0; !quit; i++ {
		start := time.Now()

		r.SetDrawColor(0x00, 0x00, 0x00, 0xff)
		r.Clear()
		drawPop(r, tab)
		drawGrid(r)
		r.Present()

		tab = conway.Update(tab)

		time.Sleep(start.Sub(time.Now()) + 128*time.Millisecond)
	}
}

func printSlice(s []int) {
	fmt.Printf("len: %d, cap: %d %v\n", len(s), cap(s), s)
}
