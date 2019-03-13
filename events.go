package main

import "github.com/veandco/go-sdl2/sdl"

func handleEvents(w *sdl.Window, game *stage, edit *edit) {
	for !game.quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch e := ev.(type) {
			case *sdl.QuitEvent:
				game.quit = true
			case *sdl.MouseButtonEvent:
				mouseButtonHandling(e, game, edit)
			case *sdl.MouseMotionEvent:
				mouseMotionHandling(e, game, edit)
			case *sdl.MouseWheelEvent:
				mouseWheelHandling(e, game)
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
			game.tab = game.newTab()
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
			edit.shiftX, edit.shiftY = game.tabIndex(x, y)
		case sdl.RELEASED:
			edit.shift = false
		}
	case sdl.K_LCTRL:
		switch k.State {
		case sdl.PRESSED:
			edit.ctrl = true
			x, y, _ := sdl.GetMouseState()
			edit.ctrlX, edit.ctrlY = game.tabIndex(x, y)
		case sdl.RELEASED:
			edit.ctrl = false
		}
	}
}

func mouseButtonHandling(m *sdl.MouseButtonEvent, game *stage,
	edit *edit) {
	edit.lastX, edit.lastY = game.tabIndex(m.X, m.Y)
	if m.State == sdl.PRESSED {
		if !edit.shift && !edit.ctrl {
			if edit.toggle {
				toggleCell(&game.tab, edit.lastX, edit.lastY)
			} else if m.Button == sdl.BUTTON_LEFT {
				reviveCell(&game.tab, edit.lastX, edit.lastY)
			} else if m.Button == sdl.BUTTON_RIGHT {
				killCell(&game.tab, edit.lastX, edit.lastY)
			}
		}
	} else if m.State == sdl.RELEASED {
		if edit.shift || edit.ctrl {
			if edit.toggle {
				editRect(&game.tab, edit, toggleCell)
			} else if m.Button == sdl.BUTTON_LEFT {
				editRect(&game.tab, edit, reviveCell)
			} else if m.Button == sdl.BUTTON_RIGHT {
				editRect(&game.tab, edit, killCell)
			}
			edit.shiftX, edit.shiftY = edit.lastX, edit.lastY
			edit.ctrlX, edit.ctrlY = edit.lastX, edit.lastY
		}
	}
}

func mouseMotionHandling(m *sdl.MouseMotionEvent, game *stage,
	edit *edit) {
	x, y := game.tabIndex(m.X, m.Y)
	if m.State == sdl.BUTTON_LEFT || m.State == 4 {
		if x != edit.lastX || y != edit.lastY {
			edit.lastX, edit.lastY = x, y
			if !edit.shift && !edit.ctrl {
				if edit.toggle {
					toggleCell(&game.tab, edit.lastX, edit.lastY)
				} else if m.State == sdl.BUTTON_LEFT {
					reviveCell(&game.tab, edit.lastX, edit.lastY)
				} else if m.State == 4 {
					killCell(&game.tab, edit.lastX, edit.lastY)
				}
			}
		}
	} else {
		edit.lastX, edit.lastY = x, y
	}
}

func mouseWheelHandling(w *sdl.MouseWheelEvent, game *stage) {
	if w.Y < 0 && game.timeEx > 0 || w.Y > 0 && game.timeEx < 600 {
		game.timeEx += int64(w.Y) * 10
	}
}
