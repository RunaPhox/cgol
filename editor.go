package main

import "math"

type point struct {
	x int32
	y int32
}

type line struct {
	a float32
	b float32
}

type edit struct {
	lastP    point
	shiftP   point
	ctrlP    point
	linectrl line
	shift    bool
	ctrl     bool
	toggle   bool
}

type cell func(*[][]byte, int32, int32)

func toggleCell(tab *[][]byte, x, y int32) {
	if (*tab)[y][x] == 0 {
		(*tab)[y][x] = 1
	} else {
		(*tab)[y][x] = 0
	}
}

func reviveCell(tab *[][]byte, x, y int32) {
	(*tab)[y][x] = 1
}

func killCell(tab *[][]byte, x, y int32) {
	(*tab)[y][x] = 0
}

func editLineLow(tab *[][]byte, x0, y0, x1, y1 int32, f cell) {
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
		f(tab, x, y)
		if D > 0 {
			y = y + yi
			D = D - 2*dx
		}
		D = D + 2*dy
	}
}

func editLineHigh(tab *[][]byte, x0, y0, x1, y1 int32, f cell) {
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
		f(tab, x, y)
		if D > 0 {
			x = x + xi
			D = D - 2*dy
		}
		D = D + 2*dx
	}
}

func editPlotLine(tab *[][]byte, x0, y0, x1, y1 int32, f cell) {
	if math.Abs(float64(y1-y0)) < math.Abs(float64(x1-x0)) {
		if x0 > x1 {
			editLineLow(tab, x1, y1, x0, y0, f)
		} else {
			editLineLow(tab, x0, y0, x1, y1, f)
		}
	} else {
		if y0 > y1 {
			editLineHigh(tab, x1, y1, x0, y0, f)
		} else {
			editLineHigh(tab, x0, y0, x1, y1, f)
		}
	}
}

func editRect(tab *[][]byte, edit *edit, f cell) {
	minX, minY, maxX, maxY := sqrPoints(
		edit.lastP.x, edit.lastP.y, edit.shiftP.x, edit.shiftP.y)
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			f(tab, int32(j), int32(i))
		}
	}
}
