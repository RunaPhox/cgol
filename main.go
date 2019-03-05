package main

import (
	"fmt"
	"math"

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

func update(tab [rows][columns]byte) [rows][columns]byte {
	buf := tab

	for i, row := range tab {
		for j := range row {
			cnt := countNeighbor(tab, j, i)

			if tab[i][j] == 0 && cnt == 3 {
				buf[i][j] = 1
			} else if tab[i][j] == 1 && (cnt < 2 || cnt > 3) {
				buf[i][j] = 0
			}

		}
	}

	return buf
}

func drawGrid(r *sdl.Renderer) {
	r.SetDrawColor(0x66, 0x66, 0x66, 0xff)
	for i := int32(cellSize); i <= height-cellSize; i += cellSize {
		r.DrawLine(0, i, width, i)
	}

	for i := int32(cellSize); i <= width-cellSize; i += cellSize {
		r.DrawLine(i, 0, i, height)
	}
}

func drawPop(r *sdl.Renderer, tab [rows][columns]byte) {
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

func countNeighbor(tab [rows][columns]byte, x, y int) (n int) {
	infY := int(math.Max(0, float64(y-1)))
	supY := int(math.Min(float64(len(tab)-1), float64(y+1)))
	infX := int(math.Max(0, float64(x-1)))
	supX := int(math.Min(float64(len(tab[0])-1), float64(x+1)))

	for i := infY; i <= supY; i++ {
		for j := infX; j <= supX; j++ {
			n += int(tab[i][j])
		}
	}

	n -= int(tab[y][x])

	return
}

func main() {

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	w, err := sdl.CreateWindow("title", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	defer w.Destroy()

	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer r.Destroy()

	var tab [rows][columns]byte

	tab[4][9] = 1
	tab[4][10] = 1
	tab[4][10] = 1
	tab[4][11] = 1

	quit := false
	for !quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch ev.(type) {
			case *sdl.QuitEvent:
				quit = true
			}
		}

		r.SetDrawColor(0x00, 0x00, 0x00, 0xff)
		r.Clear()
		drawPop(r, tab)
		drawGrid(r)
		r.Present()

		tab = update(tab)
	}
}

func printSlice(s []int) {
	fmt.Printf("len: %d, cap: %d %v\n", len(s), cap(s), s)
}
