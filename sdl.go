package main

import "github.com/veandco/go-sdl2/sdl"

type sdlContext struct {
	w *sdl.Window
	r *sdl.Renderer
}

func initSdl(game *stage) (*sdl.Window, *sdl.Renderer, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, nil, err
	}
	w, err := sdl.CreateWindow("title", sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED, game.width, game.height,
		sdl.WINDOW_SHOWN|sdl.WINDOW_FULLSCREEN)

	if err != nil {
		return w, nil, err
	}

	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return w, r, err
	}

	return w, r, nil
}

func closeSdl(wrap sdlContext) {
	if wrap.r != nil {
		wrap.r.Destroy()
	}

	if wrap.w != nil {
		wrap.w.Destroy()
	}

	sdl.Quit()
}

func toggleFullscreen(w *sdl.Window) {
	fs := w.GetFlags()&sdl.WINDOW_FULLSCREEN != 0
	if fs {
		w.SetFullscreen(0)
	} else {
		w.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	}
}
