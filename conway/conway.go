package conway

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
	var indY, indX int
	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			indY = i
			indX = j

			if indY < 0 {
				indY = len(tab) + indY
			} else {
				indY = indY % len(tab)
			}

			if indX < 0 {
				indX = len(tab[0]) + indX
			} else {
				indX = indX % len(tab[0])
			}
			n += int(tab[indY][indX])
		}
	}

	n -= int(tab[y][x])

	return
}
