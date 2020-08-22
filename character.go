package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func newCharacterAnimationRects(positions []relativeRectPosition) []*sdl.Rect {
	arr := []*sdl.Rect{}
	for _, p := range positions {
		r := sdl.Rect{
			X: int32(p.xIndex) * characterSourceWidth,
			// Without this + 1 there appears weird line above character's head
			Y: int32(p.yIndex)*characterSourceHeight + 1,
			W: characterSourceWidth,
			// Without this - 1 character starts to levitate
			H: characterSourceHeight - 1,
		}
		arr = append(arr, &r)
	}
	return arr
}

type animationsRectsSet struct {
	standingAnimationRects      []*sdl.Rect
	walkingAnimationRects       []*sdl.Rect
	jumpingUpwardAnimationRects []*sdl.Rect
	fallingAnimationRects       []*sdl.Rect
}

type character struct {
	x       int32
	y       int32
	w       int32
	h       int32
	vy      float32
	vx      float32
	texture *sdl.Texture
	animationsRectsSet
	time                  int
	facedRight            bool
	currentAnimationRects []*sdl.Rect
}

func newCharacter(x, y, w, h int32, texture *sdl.Texture) character {
	standingPlayerRects := newCharacterAnimationRects([]relativeRectPosition{{0, 1}})
	walkingPlayerRects := newCharacterAnimationRects([]relativeRectPosition{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
	})
	jumpingUpwardPlayerRects := newCharacterAnimationRects([]relativeRectPosition{{6, 1}})
	fallingPlayerRects := newCharacterAnimationRects([]relativeRectPosition{{7, 1}})
	animations := animationsRectsSet{
		standingAnimationRects:      standingPlayerRects,
		walkingAnimationRects:       walkingPlayerRects,
		jumpingUpwardAnimationRects: jumpingUpwardPlayerRects,
		fallingAnimationRects:       fallingPlayerRects,
	}
	return character{x, y, w, h, 0, 0, texture, animations, 0, true, animations.standingAnimationRects}
}

func (c *character) update(tileDestWidth int32, platforms []*platform) {
	c.x += int32(c.vx)
	c.y += int32(c.vy)
	if c.isWalking() {
		c.time++
		c.currentAnimationRects = c.walkingAnimationRects
	} else if c.isJumpingUpward() {
		c.time = 0
		c.currentAnimationRects = c.jumpingUpwardAnimationRects
	} else if c.isFalling() {
		c.time = 0
		c.currentAnimationRects = c.fallingAnimationRects
	} else {
		// character is standing
		c.time = 0
		c.currentAnimationRects = c.standingAnimationRects
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

func (c *character) isWalking() bool {
	return c.vx != 0 && c.vy == 0
}

func (c *character) isFalling() bool {
	return c.vy > 0
}

func (c *character) isJumpingUpward() bool {
	return c.vy < 0
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
	displayedFrame := c.time / 10 % len(c.currentAnimationRects)
	src := c.currentAnimationRects[displayedFrame]
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
