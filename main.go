// TODO:
// 1. Fix collisions (note that character does not take the whole tile!)
// 2. Render the lower part of the platform as well
// 3. What should be the character width?
// 4. Reduce animation speed

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 900
	windowHeight = 600
	gravity      = 0.05
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

type platform struct {
	X       int32
	Y       int32
	W       int32
	H       int32
	Texture *sdl.Texture
}

func (p *platform) draw(renderer *sdl.Renderer) {
	tileStartRect := &sdl.Rect{tileSourceWidth * 10, 0, tileSourceWidth, tileSourceHeight}
	tileMiddleRect := &sdl.Rect{tileSourceWidth * 11, 0, tileSourceWidth, tileSourceHeight}
	tileEndRect := &sdl.Rect{tileSourceWidth * 12, 0, tileSourceWidth, tileSourceHeight}
	// First row
	// First tile
	err := renderer.Copy(p.Texture, tileStartRect, &sdl.Rect{p.X - p.W/2, p.Y - p.H/2, tileDestWidth, tileDestHeight})
	if err != nil {
		log.Fatalf("could not copy texture: %v", err)
	}
	// Middle
	for tempW := tileDestWidth; tempW < p.W-tileDestWidth; tempW += tileDestWidth {
		err = renderer.Copy(p.Texture, tileMiddleRect, &sdl.Rect{p.X - p.W/2 + tempW, p.Y - p.H/2, tileDestWidth, tileDestHeight})
		if err != nil {
			log.Fatalf("could not copy texture: %v", err)
		}
	}
	// Last tile
	err = renderer.Copy(p.Texture, tileEndRect, &sdl.Rect{p.X + p.W/2 - tileDestWidth, p.Y - p.H/2, tileDestWidth, tileDestHeight})
	if err != nil {
		log.Fatalf("could not copy texture: %v", err)
	}
}

type character struct {
	X              int32
	Y              int32
	W              int32
	H              int32
	VY             float32
	VX             float32
	Texture        *sdl.Texture
	Walking        bool // Should I change it to enumerable state?
	FacedRight     bool
	DisplayedFrame int
}

func (c *character) update(tileDestWidth int32, platforms []*platform) {
	c.X += int32(c.VX)
	c.Y += int32(c.VY)
	c.VY += gravity
	if c.VX > 0 {
		c.FacedRight = true
	}
	if c.VX < 0 {
		c.FacedRight = false
	}
	// Walking animation
	if c.VX != 0 && c.VY == gravity {
		c.Walking = true
	} else {
		c.Walking = false
	}
	for _, p := range platforms {
		// If character collides with a platform from above
		// Right now it transports the character whenever he is under the platform
		if c.Y+c.H/2 >= p.Y-p.H && c.X >= p.X-p.W/2 && c.X <= p.X+p.W/2 {
			c.Y = p.Y - p.H/2 - c.H
			c.VY = 0
		}
	}
}

func (c *character) draw(renderer *sdl.Renderer) {
	if c.Walking {
		if c.DisplayedFrame > 4 {
			c.DisplayedFrame = 0
		}
		c.DisplayedFrame++
	} else {
		c.DisplayedFrame = 0
	}
	// without +1 there appears a weird line above the character head
	src := &sdl.Rect{int32(c.DisplayedFrame) * characterSourceWidth, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1}
	dst := &sdl.Rect{c.X - characterDestWidth/2, c.Y - characterDestHeight/2, characterDestWidth, characterDestHeight}
	var flip sdl.RendererFlip
	if c.FacedRight {
		flip = sdl.FLIP_NONE
	} else {
		flip = sdl.FLIP_HORIZONTAL
	}
	err := renderer.CopyEx(c.Texture, src, dst, 0, nil, flip)
	if err != nil {
		log.Fatalf("could not copy character texture: %v", err)
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

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

	platform1 := platform{windowWidth / 2, windowHeight / 2, windowWidth, 50, texBackground}
	platform2 := platform{windowWidth / 3, windowHeight / 3, windowWidth / 4, 50, texBackground}
	player := character{0, 0, tileDestWidth, tileDestHeight, 0, 0, texCharacters, false, true, 0}
	platforms := []*platform{&platform1, &platform2}

	running := true
	for running {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				// TODO: Refactor
				if sdl.K_RIGHT == e.Keysym.Sym {
					if e.State == sdl.PRESSED {
						player.VX = 1
					} else {
						player.VX = 0
					}
				}
				if sdl.K_LEFT == e.Keysym.Sym {
					if e.State == sdl.PRESSED {
						player.VX = -1
					} else {
						player.VX = 0
					}
				}
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		player.update(tileDestWidth, platforms)

		renderer.Clear()

		platform1.draw(renderer)
		platform2.draw(renderer)
		player.draw(renderer)

		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}
