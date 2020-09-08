package main

import "github.com/veandco/go-sdl2/sdl"

func newGame(texCharacters *sdl.Texture, texBackground *sdl.Texture) *game {
	player := newCharacter(tileDestWidth, tileDestHeight, texCharacters)
	platforms := createPlatforms(texBackground)

	return &game{player, platforms, 0}
}

type game struct {
	player       *character
	platforms    []*platform
	shiftScreenX int32
}

func (g *game) run(r *sdl.Renderer) (generalState, bool) {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				if sdl.K_RIGHT == e.Keysym.Sym {
					if e.State == sdl.PRESSED {
						g.player.move(true)
					} else {
						g.player.vx = 0
					}
				}
				if sdl.K_LEFT == e.Keysym.Sym {
					if e.State == sdl.PRESSED {
						g.player.move(false)
					} else {
						g.player.vx = 0
					}
				}
				if sdl.K_SPACE == e.Keysym.Sym && e.State == sdl.PRESSED {
					g.player.jump()
				}
			case *sdl.QuitEvent:
				println("Quit")
				return 0, false
			}
		}
		g.player.update(g.platforms)
		if g.player.isDead() {
			return over, true
		}
		if g.player.isCloseToRightScreenEdge() {
			g.player.x -= int32(g.player.vx)
			g.shiftScreenX++
			for _, p := range g.platforms {
				p.x--
			}
		}
		if g.player.isCloseToLeftScreenEdge() && g.shiftScreenX > 0 {
			g.player.x -= int32(g.player.vx)
			g.shiftScreenX--
			for _, p := range g.platforms {
				p.x++
			}
		}
		if g.player.x < 0 {
			g.player.x = 0
		}

		r.Clear()

		for _, p := range g.platforms {
			p.draw(r)
		}
		g.player.draw(r)

		r.Present()
	}
}
