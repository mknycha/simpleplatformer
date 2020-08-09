// TODO:
// 1. Draw everything starting from center (X,Y should be center point)
// 2. Fix collisions (note that character does not take the whole tile!)

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 1000
	windowHeight = 800
	gravity      = 0.05
)

var tileDestWidth, tileDestHeight int32

type platform struct {
	X                   int32
	Y                   int32
	Texture             *sdl.Texture
	TileStartRect       *sdl.Rect
	TileMiddleRect      *sdl.Rect
	TileEndRect         *sdl.Rect
	NumberOfMiddleTiles int
}

func (p *platform) endX() int32 {
	return int32(p.NumberOfMiddleTiles+1)*tileDestWidth + p.X
}

func (p *platform) draw(renderer *sdl.Renderer) {
	renderer.Copy(p.Texture, p.TileStartRect, &sdl.Rect{p.X, p.Y, tileDestWidth, tileDestHeight})
	for i := 1; i < p.NumberOfMiddleTiles; i++ {
		renderer.Copy(p.Texture, p.TileMiddleRect, &sdl.Rect{int32(i)*tileDestWidth + p.X, p.Y, tileDestWidth, tileDestHeight})
	}
	renderer.Copy(p.Texture, p.TileEndRect, &sdl.Rect{int32(p.NumberOfMiddleTiles)*tileDestWidth + p.X, p.Y, tileDestWidth, tileDestHeight})
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
		if c.Y+c.H >= p.Y && c.X >= p.X && c.X+(c.W/4) <= p.endX() {
			c.Y = p.Y - c.H
			c.VY = 0
		}
	}
}

func (c *character) draw(renderer *sdl.Renderer) {
	tileSrcWidth := int32(32)
	tileSrcHeight := int32(32)
	if c.Walking {
		if c.DisplayedFrame > 4 {
			c.DisplayedFrame = 0
		}
		c.DisplayedFrame++
	} else {
		c.DisplayedFrame = 0
	}
	// without +1 there appears a weird line above the character head
	src := &sdl.Rect{int32(c.DisplayedFrame) * tileSrcWidth, tileSrcHeight + 1, tileSrcWidth, tileSrcHeight - 1}
	dst := &sdl.Rect{c.X, c.Y, tileDestWidth, tileDestHeight}
	var flip sdl.RendererFlip
	if c.FacedRight {
		flip = sdl.FLIP_NONE
	} else {
		flip = sdl.FLIP_HORIZONTAL
	}
	renderer.CopyEx(c.Texture, src, dst, 0, nil, flip)
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
	texCharacters, err := img.LoadTexture(renderer, "assets/characters.png")
	if err != nil {
		log.Fatalf("could not load characters texture: %v", err)
	}

	tileSourceHeight := int32(128 / 8)
	tileSourceWidth := int32(16)
	tileStartRect := &sdl.Rect{tileSourceWidth * 7, 0, tileSourceWidth, tileSourceHeight}
	tileMiddleRect := &sdl.Rect{tileSourceWidth * 8, 0, tileSourceWidth, tileSourceHeight}
	tileEndRect := &sdl.Rect{tileSourceWidth * 9, 0, tileSourceWidth, tileSourceHeight}
	tileDestHeight = int32((windowHeight / tileSourceHeight))
	tileDestWidth = int32((windowWidth / tileSourceWidth))

	platform1 := platform{0, 200, texBackground, tileStartRect, tileMiddleRect, tileEndRect, 7}
	platform2 := platform{300, 300, texBackground, tileStartRect, tileMiddleRect, tileEndRect, 4}
	player := character{0, 0, tileDestWidth, tileDestHeight, 0, 0, texCharacters, false, true, 0}
	platforms := []*platform{&platform1, &platform2}

	running := true
	for running {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
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
