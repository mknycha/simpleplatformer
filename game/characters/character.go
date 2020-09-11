package characters

import (
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"
	"simpleplatformer/game/platforms"

	"github.com/veandco/go-sdl2/sdl"
)

func newCharacterAnimationRects(positions []common.RelativeRectPosition) []*sdl.Rect {
	arr := []*sdl.Rect{}
	for _, p := range positions {
		r := sdl.Rect{
			X: int32(p.XIndex) * constants.CharacterSourceWidth,
			// Without this + 1 there appears weird line above character's head
			Y: int32(p.YIndex)*constants.CharacterSourceHeight + 1,
			W: constants.CharacterSourceWidth,
			// Without this - 1 character starts to levitate
			H: constants.CharacterSourceHeight - 1,
		}
		arr = append(arr, &r)
	}
	return arr
}

type characterState interface {
	move(bool)
	jump()
	update([]*platforms.Platform)
	getAnimationRects() []*sdl.Rect
}

type standingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *standingState) move(right bool) {
	c := s.character
	if right {
		c.VX = 1
	} else {
		c.VX = -1
	}
	c.facedRight = right
	c.setState(c.walking)
}

func (s *standingState) jump() {
	s.character.VY = -constants.JumpSpeed
	s.character.setState(s.character.jumping)
}

func (s *standingState) update([]*platforms.Platform) {
	s.character.time = 0
}

func (s *standingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type walkingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *walkingState) move(right bool) {
	if right {
		s.character.VX = 1
	} else {
		s.character.VX = -1
	}
	s.character.facedRight = right
}

func (s *walkingState) jump() {
	c := s.character
	c.VY = -constants.JumpSpeed
	c.setState(c.jumping)
}

func (s *walkingState) update(platforms []*platforms.Platform) {
	c := s.character
	c.time++
	for _, p := range platforms {
		// If character collides with ANY platform from above
		if c.isTouchingFromAbove(p) {
			if c.VX == 0 {
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
	character      *Character
	animationRects []*sdl.Rect
}

func (s *jumpingState) move(right bool) {
	if right {
		s.character.VX = 1
	} else {
		s.character.VX = -1
	}
	s.character.facedRight = right
}

func (s *jumpingState) jump() {}

func (s *jumpingState) update([]*platforms.Platform) {
	s.character.time = 0
	s.character.VY += constants.Gravity
	if s.character.isFalling() {
		s.character.setState(s.character.falling)
	}
}

func (s *jumpingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type fallingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *fallingState) move(right bool) {
	if right {
		s.character.VX = 1
	} else {
		s.character.VX = -1
	}
	s.character.facedRight = right
}

func (s *fallingState) jump() {}

func (s *fallingState) update(platforms []*platforms.Platform) {
	c := s.character
	c.time = 0
	c.VY += constants.Gravity
	for _, p := range platforms {
		if c.isTouchingFromAbove(p) {
			c.Y = p.Y - p.H/2 - c.H
			c.VY = 0
			if c.VX == 0 {
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

type Character struct {
	X            int32
	Y            int32
	W            int32
	H            int32
	VY           float32
	VX           float32
	texture      *sdl.Texture
	time         int
	facedRight   bool
	currentState characterState

	standing characterState
	walking  characterState
	jumping  characterState
	falling  characterState
}

func (c *Character) setState(s characterState) {
	c.currentState = s
}

func NewCharacter(w, h int32, texture *sdl.Texture) *Character {
	standingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{{0, 1}})
	walkingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
	})
	jumpingUpwardPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{{6, 1}})
	fallingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{{7, 1}})

	c := Character{
		X:          0,
		Y:          0,
		W:          w,
		H:          h,
		VX:         0,
		VY:         0,
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

func (c *Character) Update(platforms []*platforms.Platform) {
	c.X += int32(c.VX)
	c.Y += int32(c.VY)
	c.currentState.update(platforms)
}

func (c *Character) reset() {
	c.X, c.Y = 0, 0
	c.VX, c.VY = 0, 0
}

func (c *Character) ResetVX() {
	c.VX = 0
}

func (c *Character) ResetVY() {
	c.VY = 0
}

func (c *Character) isTouchingFromAbove(p *platforms.Platform) bool {
	return c.Y+c.H >= p.Y-p.H/2 && c.Y+c.H <= p.Y-p.H/2+5 && c.X >= p.X-p.W/2 && c.X <= p.X+p.W/2
}

func (c *Character) isFalling() bool {
	return c.VY > 0
}

func (c *Character) isJumpingUpward() bool {
	return c.VY < 0
}

func (c *Character) IsDead() bool {
	return c.Y-c.H > constants.WindowHeight
}

func (c *Character) IsCloseToRightScreenEdge() bool {
	return c.X+(constants.TileDestWidth*5) > constants.WindowWidth
}

func (c *Character) IsCloseToLeftScreenEdge() bool {
	return c.X < (constants.TileDestWidth * 5)
}

func (c *Character) Move(right bool) {
	c.currentState.move(right)
}

func (c *Character) Jump() {
	c.currentState.jump()
}

func (c *Character) Draw(renderer *sdl.Renderer) {
	currentAnimationRects := c.currentState.getAnimationRects()
	displayedFrame := c.time / 10 % len(currentAnimationRects)
	src := currentAnimationRects[displayedFrame]
	characterDestWidth := constants.CharacterDestWidth
	characterDestHeight := constants.CharacterDestHeight
	dst := &sdl.Rect{c.X - characterDestWidth/2, c.Y - characterDestHeight/2, characterDestWidth, characterDestHeight}
	var flip sdl.RendererFlip
	if c.facedRight {
		flip = sdl.FLIP_NONE
	} else {
		flip = sdl.FLIP_HORIZONTAL
	}
	err := renderer.CopyEx(c.texture, src, dst, 0, nil, flip)
	if err != nil {
		log.Fatalf("could not copy Character texture: %v", err)
	}
}
