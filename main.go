package main

import "fmt"

const size = 10

var oldG [size][size]byte
var newG [size][size]byte

func setup() {
	oldG = [size][size]byte{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func nextGen() {
	newG = oldG
	for i, row := range oldG {
		if i == 0 || i == len(oldG)-1 {
			continue
		}

		for j := range row {
			if j == 0 || j == len(row)-1 {
				continue
			}

			cont := byte(0)
			for ix := i - 1; ix <= i+1; ix++ {
				for jx := j - 1; jx <= j+1; jx++ {
					cont += oldG[ix][jx]
				}
			}
			cont -= oldG[i][j]

			/*
				if cont != 0 {
					fmt.Printf("vecinos de (%d, %d) = %d\n", j, i, cont)
				}
			*/

			if oldG[i][j] == 0 && cont == 3 {
				newG[i][j] = 1
			}
			if oldG[i][j] > 0 && (cont < 2 || cont > 3) {
				newG[i][j] = 0
			}
		}
	}

	oldG = newG
}

func main() {
	setup()

	for {
		for i := 0; i < size; i++ {
			fmt.Println(oldG[i])
		}
		nextGen()

		x := ""
		fmt.Scanf("%s", &x)
	}
}
