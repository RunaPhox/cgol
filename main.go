package main

import (
	"time"

	"github.com/runaphox/cgol/conway"

	"github.com/veandco/go-sdl2/sdl"

	"math/rand"
)

const (
	width    = 1920
	height   = 1080
	cellSize = 20
	columns  = width / cellSize
	rows     = height / cellSize
	popul    = columns * rows
)

func randColor(rgb chan uint8) {
	for {
		rgb <- uint8(rand.Intn(256))
	}
}

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

func drawPop(r *sdl.Renderer, tab [][]byte, rgb chan uint8) {
	re, gr, bl := <- rgb, <- rgb, <- rgb
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

func draw(r *sdl.Renderer, tab [][]byte, rgb chan uint8) {
	r.SetDrawColor(0x00, 0x00, 0x00, 0xff)
	r.Clear()
	drawPop(r, tab, rgb)
	drawGrid(r)
	r.Present()
}

func toggleFullscreen(w *sdl.Window) {
	fs := w.GetFlags()&sdl.WINDOW_FULLSCREEN != 0
	if fs {
		w.SetFullscreen(0)
	} else {
		w.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	}
}

func tabIndex(x, y int32) (int32, int32) {
	yInd := y / cellSize
	xInd := x / cellSize
	return xInd, yInd
}

func toggleCell(tab *[][]byte, x, y int32) {
	if (*tab)[y][x] == 0 {
		(*tab)[y][x] = 1
	} else {
		(*tab)[y][x] = 0
	}
}

func mouseButtonHandling(m *sdl.MouseButtonEvent, tab *[][]byte) {
	if m.State == sdl.PRESSED {
		x, y := tabIndex(m.X, m.Y)
		toggleCell(tab, x, y)
	}
}

func mouseMotionHandling(m *sdl.MouseMotionEvent, tab *[][]byte,
                         lastX, lastY *int32) {
	if _, _, t := sdl.GetMouseState(); t == sdl.PRESSED {
		x, y := tabIndex(m.X, m.Y)
		if x != *lastX || y != *lastY {
			toggleCell(tab, x, y)
			*lastX, *lastY = x, y
		}
	}
}

func handleEvents(w *sdl.Window, quit, pause *bool, tab *[][]byte) {
	var lastX, lastY int32 = -1, -1
	for !*quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				*quit = true
			case *sdl.MouseButtonEvent:
				mouseButtonHandling(e, tab)
			case *sdl.MouseMotionEvent:
				mouseMotionHandling(e, tab, &lastX, &lastY)
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYUP {
					switch e.Keysym.Sym {
					case sdl.K_SPACE:
						*pause = !*pause
					case sdl.K_c:
						*tab = newTab(rows, columns)
					case sdl.K_f:
						toggleFullscreen(w)
					case sdl.K_q:
						*quit = true
					}
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

func main() {
	w, r, err := initSdl()
	if err != nil {
		panic(err)
	}
	defer closeSdl(w, r)

	tab := newTab(rows, columns)

	quit, pause := false, true
	rgb := make(chan uint8, 3)
	go randColor(rgb)
	go handleEvents(w, &quit, &pause, &tab)

	for !quit {
		start := time.Now()

		draw(r, tab, rgb)
		if !pause {
			tab = conway.Update(tab)
		}

		time.Sleep(start.Sub(time.Now()) + 32*time.Millisecond)
	}
}
