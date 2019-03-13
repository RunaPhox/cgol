package main

type stage struct {
	quit     bool
	pause    bool
	wrap     bool
	width    int32
	height   int32
	cellSize int32
	columns  int32
	rows     int32
	popul    int32
	timeEx   int64
	tab      [][]byte
}

func (s stage) tabIndex(x, y int32) (int32, int32) {
	yInd := y / s.cellSize
	xInd := x / s.cellSize
	return xInd, yInd
}

func (s stage) newTab() [][]byte {
	tab := make([][]byte, s.rows)
	for i := range tab {
		tab[i] = make([]byte, s.columns)
	}
	s.tab = tab
	return tab
}

func (s stage) population() int32 {
	s.popul = s.rows * s.columns
	return s.popul
}
