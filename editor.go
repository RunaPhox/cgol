package main

type edit struct {
	lastX  int32
	lastY  int32
	shiftX int32
	shiftY int32
	shift  bool
	toggle bool
}

type cell func(*[][]byte, int32, int32)

func toggleCell(tab *[][]byte, x, y int32) {
	if (*tab)[y][x] == 0 {
		(*tab)[y][x] = 1
	} else {
		(*tab)[y][x] = 0
	}
}

func editRect(tab *[][]byte, edit *edit, f cell) {
	minX, minY, maxX, maxY := sqrPoints(
		edit.lastX, edit.lastY, edit.shiftX, edit.shiftY)
	for i := minY; i <= maxY; i++ {
		for j := minX; j <= maxX; j++ {
			f(tab, int32(j), int32(i))
		}
	}
}

func reviveCell(tab *[][]byte, x, y int32) {
	(*tab)[y][x] = 1
}

func killCell(tab *[][]byte, x, y int32) {
	(*tab)[y][x] = 0
}
