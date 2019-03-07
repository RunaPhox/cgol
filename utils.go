package main

import "math"

func tabIndex(x, y int32) (int32, int32) {
	yInd := y / cellSize
	xInd := x / cellSize
	return xInd, yInd
}

func sqrPoints(x1, y1, x2, y2 int32) (int32, int32, int32, int32) {
	minX := int32(math.Min(float64(x1), float64(x2)))
	minY := int32(math.Min(float64(y1), float64(y2)))
	maxX := int32(math.Max(float64(x1), float64(x2)))
	maxY := int32(math.Max(float64(y1), float64(y2)))
	return minX, minY, maxX, maxY
}
