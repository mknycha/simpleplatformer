package game

import (
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"
	"simpleplatformer/game/characters"
	"simpleplatformer/game/ladders"
	"simpleplatformer/game/platforms"

	"github.com/veandco/go-sdl2/sdl"
)

func NewGame(texCharacters *sdl.Texture, texBackground *sdl.Texture, texSwoosh *sdl.Texture) *Game {
	tileDestWidth := constants.TileDestWidth
	tileDestHeight := constants.TileDestHeight
	player := characters.NewPlayerCharacter(0, tileDestHeight*7, texCharacters, texSwoosh)
	platforms := createPlatforms(texBackground)
	l1, err := ladders.NewLadder(tileDestWidth*4, tileDestHeight*4+tileDestHeight/2, tileDestWidth, tileDestHeight*13, texBackground)
	if err != nil {
		log.Fatalf("could not create a ladder: %v", err)
	}
	ladders := []*ladders.Ladder{&l1}
	slasher1 := characters.NewEnemyCharacter(
		tileDestWidth*12,
		tileDestHeight*10,
		texCharacters,
		texSwoosh,
	)
	slasher2 := characters.NewEnemyCharacter(
		tileDestWidth*6,
		tileDestHeight*6,
		texCharacters,
		texSwoosh,
	)
	snake1 := characters.NewSnake(
		tileDestWidth*20,
		tileDestHeight*10,
		texCharacters,
	)
	snake2 := characters.NewSnake(
		tileDestWidth*22,
		tileDestHeight*10,
		texCharacters,
	)
	enemies := []*characters.Character{slasher1, slasher2, snake1, snake2}
	aiControllers := []aiEnemyController{}
	for _, e := range enemies {
		aiCtrl, err := newAiControllerForEnemy(e)
		if err != nil {
			log.Fatalf("could not create enemy controller: %v", err)
		}
		aiControllers = append(aiControllers, aiCtrl)
	}
	return &Game{
		player:        player,
		platforms:     platforms,
		ladders:       ladders,
		enemies:       enemies,
		aiControllers: aiControllers,
	}
}

// TODO: Move. Can we reuse here logic used for swooshes?
func updateEnemies(platforms []*platforms.Platform, ladders []*ladders.Ladder, enemies []*characters.Character, player *characters.Character) []*characters.Character {
	result := []*characters.Character{}
	for _, e := range enemies {
		if !e.IsDead() {
			e.Update(platforms, ladders, append(enemies, player))
			result = append(result, e)
		}
	}
	return result
}

type Game struct {
	player        *characters.Character
	platforms     []*platforms.Platform
	ladders       []*ladders.Ladder
	enemies       []*characters.Character
	aiControllers []aiEnemyController
	shiftScreenX  int32
	shiftScreenY  int32
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
	if keyState[sdl.SCANCODE_UP] != 0 {
		g.player.Climb(-constants.CharacterVY, g.ladders)
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		g.player.Climb(constants.CharacterVY, g.ladders)
	}
	if keyState[sdl.SCANCODE_UP] == 0 && keyState[sdl.SCANCODE_DOWN] == 0 {
		g.player.Climb(0, g.ladders)
	}

	g.player.Update(g.platforms, g.ladders, g.enemies)
	if g.player.IsDead() {
		return common.Over, true
	}
	if g.player.IsCloseToRightScreenEdge() {
		g.player.X -= constants.CharacterVX
		g.shiftScreenX++
		for _, p := range g.platforms {
			p.X--
		}
		for _, l := range g.ladders {
			l.X--
		}
		for _, e := range g.enemies {
			e.X--
		}
		for _, ctrl := range g.aiControllers {
			ctrl.shiftPatrollingReferencePointLeft()
		}
	}
	if g.player.IsCloseToLeftScreenEdge() && g.shiftScreenX > 0 {
		g.player.X += constants.CharacterVX
		g.shiftScreenX--
		for _, p := range g.platforms {
			p.X++
		}
		for _, l := range g.ladders {
			l.X++
		}
		for _, e := range g.enemies {
			e.X++
		}
		for _, ctrl := range g.aiControllers {
			ctrl.shiftPatrollingReferencePointRight()
		}
	}
	if g.player.IsCloseToLowerScreenEdge() && g.shiftScreenY > 0 {
		diff := g.player.Y + constants.ScreenMarginHeight - constants.WindowHeight
		g.player.Y -= diff
		g.shiftScreenY -= diff
		for _, p := range g.platforms {
			p.Y -= diff
		}
		for _, l := range g.ladders {
			l.Y -= diff
		}
		for _, e := range g.enemies {
			e.Y -= diff
		}
	}
	if g.player.IsCloseToUpperScreenEdge() {
		diff := constants.ScreenMarginHeight - g.player.Y
		g.player.Y += diff
		g.shiftScreenY += diff
		for _, p := range g.platforms {
			p.Y += diff
		}
		for _, l := range g.ladders {
			l.Y += diff
		}
		for _, e := range g.enemies {
			e.Y += diff
		}
	}
	if g.player.X < 0 {
		g.player.X = 0
	}

	for _, ctrl := range g.aiControllers {
		ctrl.update(g.platforms, g.player)
	}

	g.enemies = updateEnemies(g.platforms, g.ladders, g.enemies, g.player)

	r.Clear()

	for _, p := range g.platforms {
		p.Draw(r)
	}
	for _, l := range g.ladders {
		l.Draw(r)
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
	platform1, err := platforms.NewWalkablePlatform(tileDestWidth*6, tileDestHeight*8, tileDestWidth*5, tileDestHeight*20, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	platform2, err := platforms.NewWalkablePlatform(tileDestWidth*2, tileDestHeight*14, tileDestWidth*5, tileDestHeight*6, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	platform3, err := platforms.NewWalkablePlatform(tileDestWidth*19, tileDestHeight*14, tileDestWidth*22, tileDestHeight*6, texBackground)
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
