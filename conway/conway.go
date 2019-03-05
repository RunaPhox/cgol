package conway

import "math"

// Update takes a existent byte matrix (the board) and updates it according
// to the rules of Conway's Game of Life
func Update(tab [][]byte) [][]byte {
	buf := make([][]byte, len(tab))
	for i := range buf {
		buf[i] = make([]byte, len(tab[i]))
		copy(buf[i], tab[i])
	}

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

func countNeighbor(tab [][]byte, x, y int) (n int) {
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
