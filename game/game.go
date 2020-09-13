package game

import (
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"
	"simpleplatformer/game/characters"
	"simpleplatformer/game/platforms"

	"github.com/veandco/go-sdl2/sdl"
)

func NewGame(texCharacters *sdl.Texture, texBackground *sdl.Texture) *Game {
	player := characters.NewCharacter(constants.TileDestWidth, constants.TileDestHeight, texCharacters)
	platforms := createPlatforms(texBackground)

	return &Game{player, platforms, 0}
}

type Game struct {
	player       *characters.Character
	platforms    []*platforms.Platform
	shiftScreenX int32
}

func (g *Game) Run(r *sdl.Renderer, keyState []uint8) (common.GeneralState, bool) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return 0, false
		}
	}

	if keyState[sdl.SCANCODE_LEFT] != 0 {
		g.player.Move(false)
	}
	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		g.player.Move(true)
	}
	if keyState[sdl.SCANCODE_LEFT] == 0 && keyState[sdl.SCANCODE_RIGHT] == 0 {
		g.player.ResetVX()
	}
	if keyState[sdl.SCANCODE_SPACE] != 0 {
		g.player.Jump()
	}
	if keyState[sdl.SCANCODE_LCTRL] != 0 {
		g.player.Attack()
	}

	g.player.Update(g.platforms)
	if g.player.IsDead() {
		return common.Over, true
	}
	if g.player.IsCloseToRightScreenEdge() {
		g.player.X -= constants.CharacterXSpeed
		g.shiftScreenX++
		for _, p := range g.platforms {
			p.X--
		}
	}
	if g.player.IsCloseToLeftScreenEdge() && g.shiftScreenX > 0 {
		g.player.X += constants.CharacterXSpeed
		g.shiftScreenX--
		for _, p := range g.platforms {
			p.X++
		}
	}
	if g.player.X < 0 {
		g.player.X = 0
	}

	r.Clear()

	for _, p := range g.platforms {
		p.Draw(r)
	}
	g.player.Draw(r)

	r.Present()

	return common.Play, true
}

func createPlatforms(texBackground *sdl.Texture) []*platforms.Platform {
	tileDestWidth := constants.TileDestWidth
	tileDestHeight := constants.TileDestHeight
	platform1, err := platforms.NewWalkablePlatform(tileDestWidth*6, tileDestHeight*12, tileDestWidth*5, tileDestHeight*8, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	platform2, err := platforms.NewWalkablePlatform(tileDestWidth*2, tileDestHeight*14, tileDestWidth*5, tileDestHeight*5, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	platform3, err := platforms.NewWalkablePlatform(tileDestWidth*19, tileDestHeight*14, tileDestWidth*22, tileDestHeight*5, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	// topLeftDecorationRect := &sdl.Rect{tileSourceWidth*7 + 1, 0, tileSourceWidth, tileSourceHeight - 1}
	// topMiddleDecorationRect := &sdl.Rect{tileSourceWidth * 8, 0, tileSourceWidth, tileSourceHeight - 1}
	// topRightDecorationRect := &sdl.Rect{tileSourceWidth * 9, 0, tileSourceWidth - 1, tileSourceHeight - 1}
	// midMiddleDecorationRect := &sdl.Rect{tileSourceWidth*7 + 1, tileDestHeight, tileSourceWidth - 2, tileSourceHeight - 1}
	// msg := "could not add decoration to platform2: %v"
	// err = platform2.addDecoration(topLeftDecorationRect, tileDestWidth*2, 0)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	// err = platform2.addDecoration(topMiddleDecorationRect, tileDestWidth*3, 0)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	// err = platform2.addDecoration(topRightDecorationRect, tileDestWidth*4, 0)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	// err = platform2.addDecoration(topLeftDecorationRect, tileDestWidth*10, 0)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	// err = platform2.addDecoration(topMiddleDecorationRect, tileDestWidth*11, 0)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	// err = platform2.addDecoration(topRightDecorationRect, tileDestWidth*12, 0)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	// err = platform2.addDecoration(midMiddleDecorationRect, tileDestWidth*3, tileDestHeight)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	// err = platform2.addDecoration(midMiddleDecorationRect, tileDestWidth*7, tileDestHeight*2)
	// if err != nil {
	// 	log.Fatalf(msg, err)
	// }
	return []*platforms.Platform{&platform1, &platform2, &platform3}
}
