// TODO:
// - create a level struct to encapsulate platforms,
// - Fix collisions (note that character does not take the whole tile!)
// - What should be the character width?

package main

import (
	"fmt"
	"log"
	"os"
	"simpleplatformer/common"
	"simpleplatformer/constants"
	"simpleplatformer/game"
	"simpleplatformer/game/platforms"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var state = common.Start

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %s\n", err)
		return
	}

	window, err := sdl.CreateWindow("Platformer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(constants.WindowWidth), int32(constants.WindowHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	// explosionBytes, audioSpec := sdl.LoadWAV("balloons/explode.wav")
	// if audioSpec == nil {
	// 	log.Println(sdl.GetError())
	// 	return
	// }
	// audioID, err := sdl.OpenAudioDevice("", false, audioSpec, nil, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// defer sdl.CloseAudioDevice(audioID)
	// defer sdl.FreeWAV(explosionBytes)

	// audioState := audioState{explosionBytes, audioID, audioSpec}

	// go func() {
	// 	sdl.Delay(5000)
	// 	e := sdl.QuitEvent{Type: sdl.QUIT}
	// 	sdl.PushEvent(&e)
	// }()

	var elapsedTime float32
	var g *game.Game

	renderer.SetDrawColor(uint8(66), uint8(135), uint8(245), uint8(0))

	texBackground, err := img.LoadTexture(renderer, "assets/sheet.png")
	if err != nil {
		log.Fatalf("could not load background texture: %v", err)
	}
	defer texBackground.Destroy()
	texCharacters, err := img.LoadTexture(renderer, "assets/characters.png")
	if err != nil {
		log.Fatalf("could not load characters texture: %v", err)
	}
	defer texCharacters.Destroy()

	running := true
	for running {
		frameStart := time.Now()
		if state == common.Start {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch e := event.(type) {
				case *sdl.KeyboardEvent:
					if sdl.K_SPACE == e.Keysym.Sym && e.State == sdl.PRESSED {
						g = game.NewGame(texCharacters, texBackground)
						state = common.Play
					}
				case *sdl.QuitEvent:
					println("Quit")
					running = false
					break
				}
			}
			renderer.Clear()

			displayTitle(renderer, texBackground)

			renderer.Present()
		} else if state == common.Over {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch e := event.(type) {
				case *sdl.KeyboardEvent:
					if sdl.K_SPACE == e.Keysym.Sym && e.State == sdl.PRESSED {
						g = game.NewGame(texCharacters, texBackground)
						state = common.Play
					}
				case *sdl.QuitEvent:
					println("Quit")
					running = false
					break
				}
			}
			renderer.Clear()

			err = drawText(renderer, "Game over")
			if err != nil {
				log.Fatal(err)
			}

			renderer.Present()
		} else if state == common.Play {
			var newState common.GeneralState
			newState, running = g.Run(renderer)
			if !running {
				break
			}
			state = newState
		}
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}

func displayTitle(r *sdl.Renderer, texBackground *sdl.Texture) {
	platform, err := platforms.NewWalkablePlatform(constants.WindowWidth/2, constants.WindowHeight*0.9, constants.WindowWidth, constants.WindowHeight*0.2, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	err = platform.AddUpperLeftDecoration(constants.TileDestWidth*2, 0)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	decorationWidthInTiles := 23
	for i := 3; i < decorationWidthInTiles; i++ {
		err = platform.AddUpperMiddleDecoration(constants.TileDestWidth*int32(i), 0)
		if err != nil {
			log.Fatalf("could not add a decoration: %v", err)
		}
	}
	err = platform.AddUpperRightDecoration(constants.TileDestWidth*int32(decorationWidthInTiles), 0)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	err = platform.AddLowerMiddleDecoration(constants.TileDestWidth*3, constants.TileDestHeight)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	err = platform.AddLowerMiddleDecoration(constants.TileDestWidth*7, constants.TileDestHeight*2)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	err = platform.AddLowerMiddleDecoration(constants.TileDestWidth*13, constants.TileDestHeight)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}

	platform.Draw(r)

	err = drawText(r, "King's Quest")
	if err != nil {
		log.Fatal(err)
	}
}

func drawText(r *sdl.Renderer, text string) error {
	f, err := ttf.OpenFont("assets/test.ttf", 60)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	_, _, w, h, err := t.Query()
	if err != nil {
		return fmt.Errorf("could not query texture: %v", err)
	}
	dstRect := &sdl.Rect{constants.WindowWidth/2 - w/2, constants.WindowHeight/2 - h/2, w, h}
	if err := r.Copy(t, nil, dstRect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}
