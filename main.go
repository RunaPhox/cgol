package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	size   = 40
	width  = 1920
	height = 1080
)

func newTab() [size][size]byte {
	return [size][size]byte{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func printTab(tab [size][size]byte) {
	for i := 0; i < size; i++ {
		fmt.Println(tab[i])
	}
}

func update(tab [size][size]byte) [size][size]byte {
	buf := tab

	for i, row := range tab {
		for j := range row {
			cnt := countNeighbor(tab, j, i)

			if tab[i][j] == 0 && cnt == 3 {
				buf[i][j] = 1
			}
			if tab[i][j] == 1 && (cnt < 2 || cnt > 3) {
				buf[i][j] = 0
			}

		}
	}

	return buf
}

func countNeighbor(tab [size][size]byte, x, y int) (n int) {
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

	/*
		tab := newTab()
		tab = update(tab)
		printTab(tab)
	*/

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
		r.SetDrawColor(0x50, 0x00, 0x20, 0xff)
		r.DrawRect(&sdl.Rect{X: 200, Y: 200, W: 500, H: 500})
		r.SetDrawColor(0x50, 0x70, 0x00, 0xff)
		r.FillRect(&sdl.Rect{X: 900, Y: 800, W: 200, H: 100})
		r.Present()
	}
}

func printSlice(s []int) {
	fmt.Printf("len: %d, cap: %d %v\n", len(s), cap(s), s)
}
