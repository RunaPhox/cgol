package main

import (
	"time"
	"github.com/runaphox/cgol/conway"
	"github.com/veandco/go-sdl2/sdl"
	"math"
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
}

func randColor(rgb chan uint8) {
	for v := 0.0; ; v += 0.01 {
		sin, cos, sin2 := math.Sin(v), math.Cos(v), math.Sin(2*v)
		rgb <- uint8(sin * sin * 255)
		rgb <- uint8(sin2 * sin2 * 255)
		rgb <- uint8(cos * cos * 255)
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

func mouseButtonHandling(m *sdl.MouseButtonEvent, tab *[][]byte,
	edit *edit) {
	x, y := tabIndex(m.X, m.Y)
	edit.lastX, edit.lastY = x, y
	if  k := sdl.GetKeyboardState();
	    m.State == sdl.PRESSED {
		if k[sdl.SCANCODE_LSHIFT] == 1 {
		   	edit.shiftX, edit.shiftY = x, y
		   	edit.shift = true
		} else {
		   	toggleCell(tab, x, y)
		   	edit.shift = false
		}
	} else if m.State == sdl.RELEASED {
		if k[sdl.SCANCODE_LSHIFT] == 1 {
			minY := math.Min(float64(y), float64(edit.shiftY))
			maxY := math.Max(float64(y), float64(edit.shiftY))
			minX := math.Min(float64(x), float64(edit.shiftX))
			maxX := math.Max(float64(x), float64(edit.shiftX))
		   	for i := minY; i <= maxY; i++ {
		   		for j := minX; j <= maxX; j++ {
		   			toggleCell(tab, int32(j), int32(i))
		   		}
		   	}
		}
	}
}

func mouseMotionHandling(m *sdl.MouseMotionEvent, tab *[][]byte,
	edit *edit) {
	if m.State == sdl.BUTTON_LEFT {
		x, y := tabIndex(m.X, m.Y)

		if (x != edit.lastX || y != edit.lastY) && !edit.shift {
			toggleCell(tab, x, y)
			edit.lastX, edit.lastY = x, y
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
				if e.Type == sdl.KEYUP {
					switch e.Keysym.Sym {
					case sdl.K_SPACE:
						game.pause = !game.pause
					case sdl.K_c:
						game.tab = newTab(rows, columns)
					case sdl.K_f:
						toggleFullscreen(w)
					case sdl.K_q:
						game.quit = true
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

		draw(wrap.r, game.tab, rgb)
		if !game.pause {
			game.tab = conway.Update(game.tab)
		}

		time.Sleep(start.Sub(time.Now()) + 32*time.Millisecond)
	}
}
