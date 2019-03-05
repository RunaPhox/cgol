package main

import (
	"fmt"
	"math"
)

const n = 40

func printTab(tab [n][n]byte) {
	for i := 0; i < n; i++ {
		fmt.Println(tab[i])
	}
}

func update(tab [n][n]byte) [n][n]byte {
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

func countNeighbor(tab [n][n]byte, x, y int) (n int) {
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
	tab := [n][n]byte{
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

	for {
		tab = update(tab)
		printTab(tab)

		fmt.Scanln()
	}
}

func printSlice(s []int) {
	fmt.Printf("len: %d, cap: %d %v\n", len(s), cap(s), s)
}
