package main

import (
	"time"

	"github.com/runaphox/cgol/conway"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width    = 1920
	height   = 1080
	cellSize = 20
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

	r.SetDrawColor(0xf4, 0xdf, 0x42, 0xFF)
	x, y, _ := sdl.GetMouseState()
	r.DrawRect(&sdl.Rect{
		X: x / cellSize * cellSize,
		Y: y / cellSize * cellSize,
		W: cellSize,
		H: cellSize,
	})
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

func mouseHandling(m *sdl.MouseButtonEvent, tab *[][]byte) {
	if m.Type == sdl.MOUSEBUTTONDOWN {
		yInd := m.Y / cellSize
		xInd := m.X / cellSize
		if (*tab)[yInd][xInd] == 0 {
			(*tab)[yInd][xInd] = 1
		} else {
			(*tab)[yInd][xInd] = 0
		}
	}
}

func handleEvents(quit, pause *bool, tab *[][]byte) {
	for !*quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				*quit = true
			case *sdl.MouseButtonEvent:
				mouseHandling(e, tab)
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_SPACE &&
					e.Type == sdl.KEYUP {
					*pause = !*pause
				}
			}
		}
	}
}

func initSdl() (*sdl.Window, *sdl.Renderer, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, nil, err
	}
	w, err := sdl.CreateWindow("title", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, width, height,
		sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN)

	if err != nil {
		return w, nil, err
	}

	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return w, r, err
	}

	return w, r, nil
}

func closeSdl(w *sdl.Window, r *sdl.Renderer) {
	if r != nil {
		r.Destroy()
	}

	if w != nil {
		w.Destroy()
	}

	sdl.Quit()
}

func newTab(row, col int) [][]byte {
	tab := make([][]byte, rows)
	for i := range tab {
		tab[i] = make([]byte, columns)
	}

	return tab
}

func draw(r *sdl.Renderer, tab [][]byte) {
	r.SetDrawColor(0x00, 0x00, 0x00, 0xff)
	r.Clear()
	drawPop(r, tab)
	drawGrid(r)
	r.Present()
}

func main() {
	w, r, err := initSdl()
	if err != nil {
		panic(err)
	}
	defer closeSdl(w, r)

	tab := newTab(rows, columns)

	quit, pause := false, true
	go handleEvents(&quit, &pause, &tab)

	for !quit {
		start := time.Now()

		draw(r, tab)
		if !pause {
			tab = conway.Update(tab)
		}

		time.Sleep(start.Sub(time.Now()) + 32*time.Millisecond)
	}
}
