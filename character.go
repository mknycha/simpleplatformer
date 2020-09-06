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

type characterState interface {
	move(bool)
	jump()
	update([]*platform)
	getAnimationRects() []*sdl.Rect
}

type standingState struct {
	character      *character
	animationRects []*sdl.Rect
}

func (s *standingState) move(right bool) {
	c := s.character
	if right {
		c.vx = 1
	} else {
		c.vx = -1
	}
	c.facedRight = right
	c.setState(c.walking)
}

func (s *standingState) jump() {
	s.character.vy = -jumpSpeed
	s.character.setState(s.character.jumping)
}

func (s *standingState) update([]*platform) {
	s.character.time = 0
}

func (s *standingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type walkingState struct {
	character      *character
	animationRects []*sdl.Rect
}

func (s *walkingState) move(right bool) {
	if right {
		s.character.vx = 1
	} else {
		s.character.vx = -1
	}
	s.character.facedRight = right
}

func (s *walkingState) jump() {
	c := s.character
	c.vy = -jumpSpeed
	c.setState(c.jumping)
}

func (s *walkingState) update(platforms []*platform) {
	c := s.character
	c.time++
	for _, p := range platforms {
		// If character collides with ANY platform from above
		if c.isTouchingFromAbove(p) {
			if c.vx == 0 {
				c.setState(c.standing)
			}
			return
		}
	}
	c.setState(c.falling)
}

func (s *walkingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type jumpingState struct {
	character      *character
	animationRects []*sdl.Rect
}

func (s *jumpingState) move(right bool) {
	if right {
		s.character.vx = 1
	} else {
		s.character.vx = -1
	}
	s.character.facedRight = right
}

func (s *jumpingState) jump() {}

func (s *jumpingState) update([]*platform) {
	s.character.time = 0
	s.character.vy += gravity
	if s.character.isFalling() {
		s.character.setState(s.character.falling)
	}
}

func (s *jumpingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type fallingState struct {
	character      *character
	animationRects []*sdl.Rect
}

func (s *fallingState) move(right bool) {
	if right {
		s.character.vx = 1
	} else {
		s.character.vx = -1
	}
	s.character.facedRight = right
}

func (s *fallingState) jump() {}

func (s *fallingState) update(platforms []*platform) {
	c := s.character
	c.time = 0
	c.vy += gravity
	for _, p := range platforms {
		if c.isTouchingFromAbove(p) {
			c.y = p.y - p.h/2 - c.h
			c.vy = 0
			if c.vx == 0 {
				c.setState(c.standing)
			} else {
				c.setState(c.walking)
			}
		}
	}
}

func (s *fallingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type character struct {
	x            int32
	y            int32
	w            int32
	h            int32
	vy           float32
	vx           float32
	texture      *sdl.Texture
	time         int
	facedRight   bool
	currentState characterState

	standing characterState
	walking  characterState
	jumping  characterState
	falling  characterState
}

func (c *character) setState(s characterState) {
	c.currentState = s
}

func newCharacter(w, h int32, texture *sdl.Texture) *character {
	standingPlayerRects := newCharacterAnimationRects([]relativeRectPosition{{0, 1}})
	walkingPlayerRects := newCharacterAnimationRects([]relativeRectPosition{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
	})
	jumpingUpwardPlayerRects := newCharacterAnimationRects([]relativeRectPosition{{6, 1}})
	fallingPlayerRects := newCharacterAnimationRects([]relativeRectPosition{{7, 1}})

	c := character{
		x:          0,
		y:          0,
		w:          w,
		h:          h,
		vx:         0,
		vy:         0,
		texture:    texture,
		time:       0,
		facedRight: true,
	}
	standingPlayerState := standingState{
		character:      &c,
		animationRects: standingPlayerRects,
	}
	walkingPlayerState := walkingState{
		character:      &c,
		animationRects: walkingPlayerRects,
	}
	jumpingPlayerState := jumpingState{
		character:      &c,
		animationRects: jumpingUpwardPlayerRects,
	}
	fallingPlayerState := fallingState{
		character:      &c,
		animationRects: fallingPlayerRects,
	}
	c.standing = &standingPlayerState
	c.walking = &walkingPlayerState
	c.jumping = &jumpingPlayerState
	c.falling = &fallingPlayerState
	c.setState(c.falling)
	return &c
}

func (c *character) update(platforms []*platform) {
	c.x += int32(c.vx)
	c.y += int32(c.vy)
	c.currentState.update(platforms)
}

func (c *character) reset() {
	c.x, c.y = 0, 0
	c.vx, c.vy = 0, 0
}

func (c *character) isTouchingFromAbove(p *platform) bool {
	return c.y+c.h >= p.y-p.h/2 && c.y+c.h <= p.y-p.h/2+5 && c.x >= p.x-p.w/2 && c.x <= p.x+p.w/2
}

func (c *character) isFalling() bool {
	return c.vy > 0
}

func (c *character) isJumpingUpward() bool {
	return c.vy < 0
}

func (c *character) isDead() bool {
	return c.y-c.h > windowHeight
}

func (c *character) move(right bool) {
	c.currentState.move(right)
}

func (c *character) jump() {
	c.currentState.jump()
}

func (c *character) draw(renderer *sdl.Renderer) {
	currentAnimationRects := c.currentState.getAnimationRects()
	displayedFrame := c.time / 10 % len(currentAnimationRects)
	src := currentAnimationRects[displayedFrame]
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
