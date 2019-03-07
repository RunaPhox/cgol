package main

import "github.com/veandco/go-sdl2/sdl"

func handleEvents(w *sdl.Window, game *stage, edit *edit) {
	for !game.quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				game.quit = true
			case *sdl.MouseButtonEvent:
				mouseButtonHandling(e, &game.tab, edit)
			case *sdl.MouseMotionEvent:
				mouseMotionHandling(e, &game.tab, edit)
			case *sdl.KeyboardEvent:
				keyboardHandling(e, w, game, edit)
			}
		}
	}
}

func keyboardHandling(k *sdl.KeyboardEvent, w *sdl.Window,
	game *stage, edit *edit) {
	if k.Type == sdl.KEYUP {
		switch k.Keysym.Sym {
		case sdl.K_SPACE:
			game.pause = !game.pause
		case sdl.K_c:
			game.tab = newTab(rows, columns)
		case sdl.K_f:
			toggleFullscreen(w)
		case sdl.K_q:
			game.quit = true
		case sdl.K_t:
			edit.toggle = !edit.toggle
		case sdl.K_w:
			if game.pause {
				game.wrap = !game.wrap
			}
		}
	}

	switch k.Keysym.Sym {
	case sdl.K_LSHIFT:
		switch k.State {
		case sdl.PRESSED:
			edit.shift = true
			x, y, _ := sdl.GetMouseState()
			edit.shiftX, edit.shiftY = tabIndex(x, y)
		case sdl.RELEASED:
			edit.shift = false
		}
	}
}

func mouseButtonHandling(m *sdl.MouseButtonEvent, tab *[][]byte,
	edit *edit) {
	edit.lastX, edit.lastY = tabIndex(m.X, m.Y)
	if m.State == sdl.PRESSED {
		if !edit.shift {
			if edit.toggle {
				toggleCell(tab, edit.lastX, edit.lastY)
			} else if m.Button == sdl.BUTTON_LEFT {
				reviveCell(tab, edit.lastX, edit.lastY)
			} else if m.Button == sdl.BUTTON_RIGHT {
				killCell(tab, edit.lastX, edit.lastY)
			}
		}
	} else if m.State == sdl.RELEASED {
		if edit.shift {
			if edit.toggle {
				editRect(tab, edit, toggleCell)
			} else if m.Button == sdl.BUTTON_LEFT {
				editRect(tab, edit, reviveCell)
			} else if m.Button == sdl.BUTTON_RIGHT {
				editRect(tab, edit, killCell)
			}
			edit.shiftX, edit.shiftY = edit.lastX, edit.lastY
		}
	}
}

func mouseMotionHandling(m *sdl.MouseMotionEvent, tab *[][]byte,
	edit *edit) {
	x, y := tabIndex(m.X, m.Y)
	if m.State&(sdl.BUTTON_LEFT|sdl.BUTTON_RIGHT) > 0 {
		if x != edit.lastX || y != edit.lastY {
			edit.lastX, edit.lastY = x, y
			if !edit.shift {
				if edit.toggle {
					toggleCell(tab, edit.lastX, edit.lastY)
				} else if m.State == sdl.BUTTON_LEFT {
					reviveCell(tab, edit.lastX, edit.lastY)
				} else if m.State&sdl.BUTTON_RIGHT > 0 {
					killCell(tab, edit.lastX, edit.lastY)
				}
			}
		}
	} else {
		edit.lastX, edit.lastY = x, y
	}
}
