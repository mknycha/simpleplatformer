// TODO:
// - create a level struct to encapsulate platforms,
// - Fix collisions (note that character does not take the whole tile!)
// - What should be the character width?

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 860
	windowHeight = 510
	gravity      = 0.05
	jumpSpeed    = 4
)

const (
	scaleX                = windowWidth / 288
	scaleY                = windowHeight / 172
	tileSourceWidth       = int32(16)
	tileSourceHeight      = int32(128 / 8)
	tileDestWidth         = int32(tileSourceWidth * scaleX)
	tileDestHeight        = int32(tileSourceHeight * scaleY)
	characterSourceWidth  = int32(32)
	characterSourceHeight = int32(32)
	characterDestWidth    = int32(characterSourceWidth * scaleX)
	characterDestHeight   = int32(characterSourceHeight * scaleY)
)

type gameState int

const (
	start gameState = iota
	play
	over
)

var state = start

var drawingFromX int32

type relativeRectPosition struct{ xIndex, yIndex int }

func main() {
	drawingFromX = 0
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %s\n", err)
		return
	}

	window, err := sdl.CreateWindow("Platformer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth), int32(windowHeight), sdl.WINDOW_SHOWN)
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

	player := newCharacter(tileDestWidth, tileDestHeight, texCharacters)
	platforms := createPlatforms(texBackground)

	running := true
	for running {
		frameStart := time.Now()
		if state == start {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch e := event.(type) {
				case *sdl.KeyboardEvent:
					if sdl.K_SPACE == e.Keysym.Sym && e.State == sdl.PRESSED {
						state = play
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
		} else if state == over {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch e := event.(type) {
				case *sdl.KeyboardEvent:
					if sdl.K_SPACE == e.Keysym.Sym && e.State == sdl.PRESSED {
						state = play
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
		} else if state == play {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch e := event.(type) {
				case *sdl.KeyboardEvent:
					if sdl.K_RIGHT == e.Keysym.Sym {
						if e.State == sdl.PRESSED {
							player.move(true)
						} else {
							player.vx = 0
						}
					}
					if sdl.K_LEFT == e.Keysym.Sym {
						if e.State == sdl.PRESSED {
							player.move(false)
						} else {
							player.vx = 0
						}
					}
					if sdl.K_SPACE == e.Keysym.Sym && e.State == sdl.PRESSED {
						player.jump()
					}
				case *sdl.QuitEvent:
					println("Quit")
					running = false
					break
				}
			}
			player.update(platforms)
			if player.isDead() {
				state = over
				player.reset()
			}
			if player.x+(tileDestWidth*3) > windowWidth {
				for _, p := range platforms {
					p.x -= 3
					player.x -= int32(player.vx)
					drawingFromX++
				}
			}
			if player.x < (tileDestWidth*3) && drawingFromX > 0 {
				for _, p := range platforms {
					p.x += 3
					player.x -= int32(player.vx)
					drawingFromX--
				}
			}
			if player.x < drawingFromX {
				player.x = 0
			}

			renderer.Clear()

			for _, p := range platforms {
				p.draw(renderer)
			}
			player.draw(renderer)

			renderer.Present()
		}
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}

func createPlatforms(texBackground *sdl.Texture) []*platform {
	platform1, err := newWalkablePlatform(windowWidth/3, windowHeight*0.7, windowWidth/4, windowHeight*0.5, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	platform2, err := newWalkablePlatform(windowWidth*0.1, windowHeight*0.9, windowWidth*0.25, windowHeight*0.5, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	platform3, err := newWalkablePlatform(tileDestWidth*16, tileDestHeight*14, tileDestWidth*24, tileDestHeight*5, texBackground)
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
	return []*platform{&platform1, &platform2, &platform3}
}

func displayTitle(r *sdl.Renderer, texBackground *sdl.Texture) {
	platform, err := newWalkablePlatform(windowWidth/2, windowHeight*0.9, windowWidth, windowHeight*0.2, texBackground)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	err = platform.addUpperLeftDecoration(tileDestWidth*2, 0)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	decorationWidthInTiles := 23
	for i := 3; i < decorationWidthInTiles; i++ {
		err = platform.addUpperMiddleDecoration(tileDestWidth*int32(i), 0)
		if err != nil {
			log.Fatalf("could not add a decoration: %v", err)
		}
	}
	err = platform.addUpperRightDecoration(tileDestWidth*int32(decorationWidthInTiles), 0)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	err = platform.addLowerMiddleDecoration(tileDestWidth*3, tileDestHeight)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	err = platform.addLowerMiddleDecoration(tileDestWidth*7, tileDestHeight*2)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}
	err = platform.addLowerMiddleDecoration(tileDestWidth*13, tileDestHeight)
	if err != nil {
		log.Fatalf("could not add a decoration: %v", err)
	}

	platform.draw(r)

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
	dstRect := &sdl.Rect{windowWidth/2 - w/2, windowHeight/2 - h/2, w, h}
	if err := r.Copy(t, nil, dstRect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}
