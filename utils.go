package main

import (
	"math"
	"time"
)

func eval(x float32, l line) float32 {
	n := l.a*x + l.b
	return n
}

func lineFunc(x1, y1, x2, y2 float32) (float32, float32) {
	a := (y2 - y1) / (x2 - x1)
	b := y1 - a*x1
	return a, b
}

func sqrPoints(x1, y1, x2, y2 int32) (int32, int32, int32, int32) {
	minX := int32(math.Min(float64(x1), float64(x2)))
	minY := int32(math.Min(float64(y1), float64(y2)))
	maxX := int32(math.Max(float64(x1), float64(x2)))
	maxY := int32(math.Max(float64(y1), float64(y2)))
	return minX, minY, maxX, maxY
}

func ipow(base, exp int32) time.Duration {
	var result int32 = 1
	for {
		if exp%2 != 0 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}
	return time.Duration(result)
}
