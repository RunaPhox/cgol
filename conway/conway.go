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
	n += int(tab[(y-1)%len(tab)][x])
	n += int(tab[(y+1)%len(tab)][x])
	n += int(tab[(y-1)%len(tab)][(x-1)%len(tab)])
	n += int(tab[(y+1)%len(tab)][(x-1)%len(tab)])
	n += int(tab[y][(x-1)%len(tab)])
	n += int(tab[y][(x+1)%len(tab)])
	n += int(tab[(y-1)%len(tab)][(x+1)%len(tab)])
	n += int(tab[(y+1)%len(tab)][(x+1)%len(tab)])

	return
}
