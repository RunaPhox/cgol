package main

type stage struct {
	quit  bool
	pause bool
	wrap bool
	tab   [][]byte
}

func newTab(row, col int) [][]byte {
	tab := make([][]byte, rows)
	for i := range tab {
		tab[i] = make([]byte, columns)
	}

	return tab
}
