package main

type edit struct {
	lastX  int32
	lastY  int32
	shiftX int32
	shiftY int32
	ctrlX  int32
	ctrlY  int32
	shift  bool
	ctrl   bool
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
	if edit.shift && !edit.ctrl {
		minX, minY, maxX, maxY := sqrPoints(
			edit.lastX, edit.lastY, edit.shiftX, edit.shiftY)
		for i := minY; i <= maxY; i++ {
			for j := minX; j <= maxX; j++ {
				f(tab, int32(j), int32(i))
			}
		}
	} else if !edit.shift && edit.ctrl {
		minX, minY, maxX, maxY := sqrPoints(
			edit.lastX, edit.lastY, edit.ctrlX, edit.ctrlY)
		if maxX-minX <= maxY-minY {
			for i := minY; i <= maxY; i++ {
				f(tab, edit.ctrlX, int32(i))
			}
		} else {
			for i := minX; i <= maxX; i++ {
				f(tab, int32(i), edit.ctrlY)
			}
		}
	}
}

func reviveCell(tab *[][]byte, x, y int32) {
	(*tab)[y][x] = 1
}

func killCell(tab *[][]byte, x, y int32) {
	(*tab)[y][x] = 0
}
