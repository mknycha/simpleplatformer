package game

import (
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"
	"simpleplatformer/game/characters"
	"simpleplatformer/game/platforms"

	"github.com/veandco/go-sdl2/sdl"
)

func NewGame(texCharacters *sdl.Texture, texBackground *sdl.Texture, texSwoosh *sdl.Texture) *Game {
	tileDestWidth := constants.TileDestWidth
	tileDestHeight := constants.TileDestHeight
	player := characters.NewPlayerCharacter(0, 0, texCharacters, texSwoosh)
	platforms := createPlatforms(texBackground)
	enemy := characters.NewEnemyCharacter(
		tileDestWidth*19,
		tileDestHeight*10,
		texCharacters,
		texSwoosh,
	)
	enemies := []*characters.Character{enemy}
	aiControllers := []*aiController{
		newAiController(enemy),
	}

	return &Game{player, platforms, enemies, aiControllers, 0}
}

// TODO: Move. Can we reuse here logic used for swooshes?
func updateEnemies(platforms []*platforms.Platform, enemies []*characters.Character, player *characters.Character) []*characters.Character {
	result := []*characters.Character{}
	for _, e := range enemies {
		e.Update(platforms, []*characters.Character{player})
		result = append(result, e)
		// TODO: Destroy when fell off the screen
		// if e.X == false {
		// 	e.update()
		// 	result = append(result, e)
		// }
	}
	return result
}

type Game struct {
	player        *characters.Character
	platforms     []*platforms.Platform
	enemies       []*characters.Character
	aiControllers []*aiController
	shiftScreenX  int32
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
		g.player.Move(-constants.CharacterVX)
	}
	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		g.player.Move(constants.CharacterVX)
	}
	if keyState[sdl.SCANCODE_LEFT] == 0 && keyState[sdl.SCANCODE_RIGHT] == 0 {
		g.player.Move(0)
	}
	if keyState[sdl.SCANCODE_SPACE] != 0 {
		g.player.Jump()
	}
	if keyState[sdl.SCANCODE_LCTRL] != 0 {
		g.player.Attack()
	}

	g.player.Update(g.platforms, g.enemies)
	if g.player.IsDead() {
		return common.Over, true
	}
	if g.player.IsCloseToRightScreenEdge() {
		g.player.X -= constants.CharacterVX
		g.shiftScreenX++
		for _, p := range g.platforms {
			p.X--
		}
		for _, e := range g.enemies {
			e.X--
		}
	}
	if g.player.IsCloseToLeftScreenEdge() && g.shiftScreenX > 0 {
		g.player.X += constants.CharacterVX
		g.shiftScreenX--
		for _, p := range g.platforms {
			p.X++
		}
		for _, e := range g.enemies {
			e.X++
		}
	}
	if g.player.X < 0 {
		g.player.X = 0
	}

	// TODO: Make a list
	g.aiControllers[0].update()

	g.enemies = updateEnemies(g.platforms, g.enemies, g.player)

	r.Clear()

	for _, p := range g.platforms {
		p.Draw(r)
	}
	for _, e := range g.enemies {
		e.Draw(r)
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
