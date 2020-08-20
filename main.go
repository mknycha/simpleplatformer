// TODO:
// - There could be a struct for a set of animations
// - Check the naming for animation rects
// - Instead of character state - there could be a function, telling if character is e.g. falling
// - Start game state, and game over state
// - Fix collisions (note that character does not take the whole tile!)
// - What should be the character width?

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
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

type platformRects struct {
	topLeftRect   *sdl.Rect
	topMiddleRect *sdl.Rect
	topRightRect  *sdl.Rect
	midLeftRect   *sdl.Rect
	midMiddleRect *sdl.Rect
	midRightRect  *sdl.Rect
}

type platformDecoration struct {
	texture *sdl.Texture
	srcRect *sdl.Rect
	dstRect *sdl.Rect
}

func (pd *platformDecoration) draw(renderer *sdl.Renderer) {
	err := renderer.Copy(pd.texture, pd.srcRect, pd.dstRect)
	if err != nil {
		log.Fatalf("could not copy platform decoration texture: %v", err)
	}
}

type platform struct {
	x           int32
	y           int32
	w           int32
	h           int32
	texture     *sdl.Texture
	sourceRects platformRects
	decorations []platformDecoration
}

func newPlatform(x, y, w, h int32, texture *sdl.Texture, sourceRects platformRects) (platform, error) {
	if w < tileDestWidth*3 {
		return platform{}, fmt.Errorf("width value: %v must be higher (at least %v)", w, tileDestWidth*3)
	}
	return platform{x, y, w, h, texture, sourceRects, []platformDecoration{}}, nil
}

// addDecoration adds a decoration tile from src of the platform texture to the position (relative to the platform)
func (p *platform) addDecoration(srcRect *sdl.Rect, x, y int32) error {
	if p.x-p.w/2+x+tileDestWidth > p.x+p.w/2 || p.x-p.w/2+x < p.x-p.w/2 {
		return fmt.Errorf("invalid decoration position x: %v. Decoration width exceeds platform width (%v)", x, p.w)
	}
	if p.y-p.h/2+y+tileDestHeight > p.y+p.h/2 || p.y-p.h/2+y < p.y-p.h/2 {
		return fmt.Errorf("invalid decoration position y: %v. Decoration height exceeds platform height (%v)", y, p.h)
	}
	dstRect := &sdl.Rect{p.x - p.w/2 + x, p.y - p.h/2 + y, tileDestWidth, tileDestHeight}
	pd := platformDecoration{p.texture, srcRect, dstRect}
	p.decorations = append(p.decorations, pd)
	return nil
}

func (p *platform) draw(renderer *sdl.Renderer) {
	// Top row
	p.drawRow(renderer, p.sourceRects.topLeftRect, p.sourceRects.topMiddleRect, p.sourceRects.topRightRect, 0)
	// Other rows
	for y := tileDestHeight; y < p.h; y += tileDestHeight - 1 {
		p.drawRow(renderer, p.sourceRects.midLeftRect, p.sourceRects.midMiddleRect, p.sourceRects.midRightRect, y)
	}
	for _, pd := range p.decorations {
		pd.draw(renderer)
	}
}

func (p *platform) drawRow(renderer *sdl.Renderer, tileLeftRect, tileMiddleRect, tileRightRect *sdl.Rect, y int32) {
	err := renderer.Copy(p.texture, tileLeftRect, &sdl.Rect{p.x - p.w/2, p.y - p.h/2 + y, tileDestWidth, tileDestHeight})
	if err != nil {
		log.Fatalf("could not copy platform left texture: %v", err)
	}
	for x := tileDestWidth; x < p.w-tileDestWidth; x += tileDestWidth {
		err = renderer.Copy(p.texture, tileMiddleRect, &sdl.Rect{p.x - p.w/2 + x, p.y - p.h/2 + y, tileDestWidth, tileDestHeight})
		if err != nil {
			log.Fatalf("could not copy platform middle texture: %v", err)
		}
	}
	err = renderer.Copy(p.texture, tileRightRect, &sdl.Rect{p.x + p.w/2 - tileDestWidth, p.y - p.h/2 + y, tileDestWidth, tileDestHeight})
	if err != nil {
		log.Fatalf("could not copy platform right texture: %v", err)
	}
}

type characterState int

type character struct {
	x                           int32
	y                           int32
	w                           int32
	h                           int32
	vy                          float32
	vx                          float32
	texture                     *sdl.Texture
	time                        int
	facedRight                  bool
	currentAnimation            []*sdl.Rect
	walkingAnimationRects       []*sdl.Rect
	standingAnimationRects      []*sdl.Rect
	jumpingUpwardAnimationRects []*sdl.Rect
	fallingAnimationRects       []*sdl.Rect
}

func (c *character) update(tileDestWidth int32, platforms []*platform) {
	c.x += int32(c.vx)
	c.y += int32(c.vy)
	if c.vx != 0 && c.vy == 0 {
		c.time++
		c.currentAnimation = c.walkingAnimationRects
	} else if c.vy < 0 { // jumping, going upward
		c.time = 0
		c.currentAnimation = c.jumpingUpwardAnimationRects
	} else if c.vy > 0 { // falling down
		c.time = 0
		c.currentAnimation = c.fallingAnimationRects
	} else {
		// character is standing
		c.time = 0
		c.currentAnimation = c.standingAnimationRects
	}
	c.vy += gravity
	for _, p := range platforms {
		// If character collides with a platform from above
		if c.y+c.h >= p.y-p.h/2 && c.y+c.h <= p.y-p.h/2+5 && c.x >= p.x-p.w/2 && c.x <= p.x+p.w/2 {
			// If character is standing or falling down
			if c.vy >= 0 {
				c.y = p.y - p.h/2 - c.h
				c.vy = 0
			}
		}
	}
}

func (c *character) move(right bool) {
	if right {
		c.vx = 1
	} else {
		c.vx = -1
	}
	c.facedRight = right
}

func (c *character) jump() {
	if c.vy == 0 {
		c.vy = -jumpSpeed
	}
}

func (c *character) draw(renderer *sdl.Renderer) {
	displayedFrame := c.time / 10 % len(c.currentAnimation)
	src := c.currentAnimation[displayedFrame]
	dst := &sdl.Rect{c.x - characterDestWidth/2, c.y - characterDestHeight/2, characterDestWidth, characterDestHeight}
	var flip sdl.RendererFlip
	if c.facedRight {
		flip = sdl.FLIP_NONE
	} else {
		flip = sdl.FLIP_HORIZONTAL
	}
	err := renderer.CopyEx(c.texture, src, dst, 0, nil, flip)
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

	walkablePlatformRects := platformRects{
		topLeftRect:   &sdl.Rect{tileSourceWidth * 10, 0, tileSourceWidth, tileSourceHeight},
		topMiddleRect: &sdl.Rect{tileSourceWidth * 11, 0, tileSourceWidth, tileSourceHeight},
		topRightRect:  &sdl.Rect{tileSourceWidth * 12, 0, tileSourceWidth, tileSourceHeight},
		midLeftRect:   &sdl.Rect{tileSourceWidth * 10, tileSourceHeight, tileSourceWidth, tileSourceHeight},
		midMiddleRect: &sdl.Rect{tileSourceWidth * 11, tileSourceHeight, tileSourceWidth, tileSourceHeight},
		midRightRect:  &sdl.Rect{tileSourceWidth * 12, tileSourceHeight, tileSourceWidth, tileSourceHeight},
	}
	// topLeftDecorationRect := &sdl.Rect{tileSourceWidth*7 + 1, 0, tileSourceWidth, tileSourceHeight - 1}
	// topMiddleDecorationRect := &sdl.Rect{tileSourceWidth * 8, 0, tileSourceWidth, tileSourceHeight - 1}
	// topRightDecorationRect := &sdl.Rect{tileSourceWidth * 9, 0, tileSourceWidth - 1, tileSourceHeight - 1}
	// midMiddleDecorationRect := &sdl.Rect{tileSourceWidth*7 + 1, tileDestHeight, tileSourceWidth - 2, tileSourceHeight - 1}

	platform1, err := newPlatform(windowWidth/3, windowHeight*0.7, windowWidth/4, windowHeight*0.5, texBackground, walkablePlatformRects)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
	platform2, err := newPlatform(windowWidth*0.1, windowHeight*0.9, windowWidth*0.25, windowHeight*0.5, texBackground, walkablePlatformRects)
	if err != nil {
		log.Fatalf("could not create a platform: %v", err)
	}
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

	// without +1 there appears a weird line above the character head
	standingPlayerRects := []*sdl.Rect{
		{0, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1},
	}
	walkingPlayerRects := []*sdl.Rect{
		{characterSourceWidth, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1},
		{characterSourceWidth * 2, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1},
		{characterSourceWidth * 3, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1},
		{characterSourceWidth * 4, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1},
	}
	jumpingUpwardPlayerRects := []*sdl.Rect{
		{characterSourceWidth * 6, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1},
	}
	fallingPlayerRects := []*sdl.Rect{
		{characterSourceWidth * 7, characterSourceHeight + 1, characterSourceWidth, characterSourceHeight - 1},
	}
	player := character{0, 0, tileDestWidth, tileDestHeight, 0, 0, texCharacters, 0, true, standingPlayerRects, walkingPlayerRects, standingPlayerRects, jumpingUpwardPlayerRects, fallingPlayerRects}
	platforms := []*platform{&platform1, &platform2}

	running := true
	for running {
		frameStart := time.Now()
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
		player.update(tileDestWidth, platforms)

		renderer.Clear()

		for _, p := range platforms {
			p.draw(renderer)
		}
		player.draw(renderer)

		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}
