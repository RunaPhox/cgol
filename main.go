package main

import (
	"github.com/runaphox/cgol/conway"
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

const (
	width    = 1920
	height   = 1080
	cellSize = 20
	columns  = width / cellSize
	rows     = height / cellSize
	popul    = columns * rows
)

type sdlContext struct {
	w *sdl.Window
	r *sdl.Renderer
}

type stage struct {
	quit  bool
	pause bool
	tab   [][]byte
}

type edit struct {
	lastX  int32
	lastY  int32
	shiftX int32
	shiftY int32
	shift  bool
	toggle bool
}

func randColor(rgb chan uint8) {
	for v := 0.0; ; v += 0.01 {
		sin, cos, sin2 := math.Sin(v), math.Cos(v), math.Sin(2*v)
		rgb <- uint8(sin * sin * 255)
		rgb <- uint8(sin2 * sin2 * 255)
		rgb <- uint8(cos * cos * 255)
	}
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

func draw(r *sdl.Renderer, tab [][]byte, rgb chan uint8,
	edit *edit) {
	r.SetDrawColor(0x00, 0x00, 0x00, 0xff)
	r.Clear()
	drawPop(r, tab, rgb)
	drawGrid(r, edit)
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

type cell func(*[][]byte, int32, int32)

func toggleCell(tab *[][]byte, x, y int32) {
	if (*tab)[y][x] == 0 {
		(*tab)[y][x] = 1
	} else {
		(*tab)[y][x] = 0
	}
}

func reviveCell(tab *[][]byte, x, y int32) { (*tab)[y][x] = 1 }

func killCell(tab *[][]byte, x, y int32) { (*tab)[y][x] = 0 }

func sqrPoints(x1, y1, x2, y2 int32) (int32, int32, int32, int32) {
	minX := int32(math.Min(float64(x1), float64(x2)))
	minY := int32(math.Min(float64(y1), float64(y2)))
	maxX := int32(math.Max(float64(x1), float64(x2)))
	maxY := int32(math.Max(float64(y1), float64(y2)))
	return minX, minY, maxX, maxY
}

func cellSqr(tab *[][]byte, edit *edit, f cell) {
	minX, minY, maxX, maxY := sqrPoints(
		edit.lastX, edit.lastY, edit.shiftX, edit.shiftY)
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			f(tab, int32(j), int32(i))
		}
	}
}

func mouseButtonHandling(m *sdl.MouseButtonEvent, tab *[][]byte,
	edit *edit) {
	edit.lastX, edit.lastY = tabIndex(m.X, m.Y)
	if m.State == sdl.PRESSED {
		if !edit.shift {
			if edit.toggle {
				toggleCell(tab, edit.lastX, edit.lastY)
			} else if m.Button == sdl.BUTTON_LEFT {
				reviveCell(tab, edit.lastX, edit.lastY)
			} else if m.Button == sdl.BUTTON_RIGHT {
				killCell(tab, edit.lastX, edit.lastY)
			}
		}
	} else if m.State == sdl.RELEASED {
		if edit.shift {
			if edit.toggle {
				cellSqr(tab, edit, toggleCell)
			} else if m.Button == sdl.BUTTON_LEFT {
				cellSqr(tab, edit, reviveCell)
			} else if m.Button == sdl.BUTTON_RIGHT {
				cellSqr(tab, edit, killCell)
			}
		}
	}
}

func mouseMotionHandling(m *sdl.MouseMotionEvent, tab *[][]byte,
	edit *edit) {
	x, y := tabIndex(m.X, m.Y)
	if m.State == sdl.BUTTON_LEFT || m.State == /*sdl.BUTTON_RIGHT*/ 4 {
		if x != edit.lastX || y != edit.lastY {
			edit.lastX, edit.lastY = x, y
			if !edit.shift {
				if edit.toggle {
					toggleCell(tab, edit.lastX, edit.lastY)
				} else if m.State == sdl.BUTTON_LEFT {
					reviveCell(tab, edit.lastX, edit.lastY)
				} else if m.State == /*sdl.BUTTON_RIGHT*/ 4 {
					killCell(tab, edit.lastX, edit.lastY)
				}
			}
		}
	} else {
		edit.lastX, edit.lastY = x, y
		edit.shiftX, edit.shiftY = x, y
	}
}

func keyboardHandling(k *sdl.KeyboardEvent, w *sdl.Window,
	game *stage, edit *edit) {
	if k.Type == sdl.KEYUP {
		switch k.Keysym.Sym {
		case sdl.K_SPACE:
			game.pause = !game.pause
		case sdl.K_c:
			game.tab = newTab(rows, columns)
		case sdl.K_f:
			toggleFullscreen(w)
		case sdl.K_q:
			game.quit = true
		case sdl.K_t:
			edit.toggle = !edit.toggle
		}
	}

	switch k.Keysym.Sym {
	case sdl.K_LSHIFT:
		switch k.State {
		case sdl.PRESSED:
			edit.shift = true
		case sdl.RELEASED:
			edit.shift = false
			x, y, _ := sdl.GetMouseState()
			edit.shiftX, edit.shiftY = tabIndex(x, y)
		}
	}
}

func handleEvents(w *sdl.Window, game *stage, edit *edit) {
	for !game.quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				game.quit = true
			case *sdl.MouseButtonEvent:
				mouseButtonHandling(e, &game.tab, edit)
			case *sdl.MouseMotionEvent:
				mouseMotionHandling(e, &game.tab, edit)
			case *sdl.KeyboardEvent:
				keyboardHandling(e, w, game, edit)
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

func closeSdl(wrap sdlContext) {
	if wrap.r != nil {
		wrap.r.Destroy()
	}

	if wrap.w != nil {
		wrap.w.Destroy()
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
