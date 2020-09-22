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
	attack()
	update([]*platforms.Platform)
	kill(float32)
	getAnimationRects() []*sdl.Rect
}

type standingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *standingState) move(right bool) {
	c := s.character
	if right {
		c.vx = constants.CharacterXSpeed
	} else {
		c.vx = -constants.CharacterXSpeed
	}
	c.facedRight = right
	c.setState(c.walking)
}

func (s *standingState) jump() {
	s.character.vy = -constants.JumpSpeed
	s.character.setState(s.character.jumping)
}

func (s *standingState) attack() {
	conditionalSwitchToAttackingState(s.character)
}

func (s *standingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *standingState) update([]*platforms.Platform) {}

func (s *standingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type walkingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *walkingState) move(right bool) {
	if right {
		s.character.vx = constants.CharacterXSpeed
	} else {
		s.character.vx = -constants.CharacterXSpeed
	}
	s.character.facedRight = right
}

func (s *walkingState) jump() {
	c := s.character
	c.vy = -constants.JumpSpeed
	c.setState(c.jumping)
}

func (s *walkingState) attack() {
	conditionalSwitchToAttackingState(s.character)
}

func (s *walkingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *walkingState) update(platforms []*platforms.Platform) {
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
	character      *Character
	animationRects []*sdl.Rect
}

func (s *jumpingState) move(right bool) {
	if right {
		s.character.vx = constants.CharacterXSpeed
	} else {
		s.character.vx = -constants.CharacterXSpeed
	}
	s.character.facedRight = right
}

func (s *jumpingState) jump() {}

func (s *jumpingState) attack() {}

func (s *jumpingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *jumpingState) update([]*platforms.Platform) {
	s.character.time = 0
	s.character.vy += constants.Gravity
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
		s.character.vx = constants.CharacterXSpeed
	} else {
		s.character.vx = -constants.CharacterXSpeed
	}
	s.character.facedRight = right
}

func (s *fallingState) jump() {}

func (s *fallingState) attack() {}

func (s *fallingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *fallingState) update(platforms []*platforms.Platform) {
	c := s.character
	c.time = 0
	c.vy += constants.Gravity
	for _, p := range platforms {
		if c.isTouchingFromAbove(p) {
			c.Y = p.Y - p.H/2 - c.H
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

type attackingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *attackingState) move(right bool) {}

func (s *attackingState) jump() {}

func (s *attackingState) attack() {}

func (s *attackingState) kill(float32) {
	s.character.setState(s.character.dead)
}

func (s *attackingState) update(platforms []*platforms.Platform) {
	c := s.character
	c.vx = 0
	c.stamina = 0
	c.time++
	if c.time > len(s.getAnimationRects())*10 {
		c.setState(c.standing)
	}
}

func (s *attackingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type deadState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *deadState) move(right bool) {}

func (s *deadState) jump() {}

func (s *deadState) attack() {}

func (s *deadState) kill(float32) {}

func (s *deadState) update(platforms []*platforms.Platform) {
	c := s.character
	c.time++
	c.vy += constants.Gravity
}

func (s *deadState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type Character struct {
	X             int32
	Y             int32
	W             int32
	H             int32
	vy            float32
	vx            float32
	texture       *sdl.Texture
	swooshTexture *sdl.Texture
	time          int
	facedRight    bool
	currentState  characterState
	swooshes      []*swoosh
	stamina       int

	standing  characterState
	walking   characterState
	jumping   characterState
	attacking characterState
	dead      characterState
	falling   characterState
}

func (c *Character) setState(s characterState) {
	c.time = 0
	c.currentState = s
}

func NewPlayerCharacter(x, y int32, characterTexture *sdl.Texture, swooshTexture *sdl.Texture) *Character {
	standingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{{0, 1}})
	walkingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
	})
	jumpingUpwardPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{{6, 1}})
	fallingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{{7, 1}})
	attackingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{12, 1},
		{11, 1},
		{12, 1},
		{13, 1},
	})
	deadPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{9, 1},
		{10, 1},
	})

	c := Character{
		X:             x,
		Y:             y,
		W:             constants.TileDestWidth,
		H:             constants.TileDestHeight,
		vx:            0,
		vy:            0,
		texture:       characterTexture,
		swooshTexture: swooshTexture,
		stamina:       constants.CharacterStaminaMax,
		time:          0,
		facedRight:    true,
		swooshes:      []*swoosh{},
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
	attackingPlayerState := attackingState{
		character:      &c,
		animationRects: attackingPlayerRects,
	}
	deadPlayerState := deadState{
		character:      &c,
		animationRects: deadPlayerRects,
	}
	c.standing = &standingPlayerState
	c.walking = &walkingPlayerState
	c.jumping = &jumpingPlayerState
	c.falling = &fallingPlayerState
	c.attacking = &attackingPlayerState
	c.dead = &deadPlayerState
	c.setState(c.falling)
	return &c
}

func NewEnemyCharacter(x, y int32, characterTexture *sdl.Texture, swooshTexture *sdl.Texture) *Character {
	standingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{{0, 0}})
	walkingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
	})
	jumpingUpwardEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{{6, 0}})
	fallingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{{7, 0}})
	attackingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{12, 0},
		{11, 0},
		{12, 0},
		{13, 0},
	})
	deadEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{9, 0},
		{10, 0},
	})

	c := Character{
		X:             x,
		Y:             y,
		W:             constants.TileDestWidth,
		H:             constants.TileDestHeight,
		vx:            0,
		vy:            0,
		texture:       characterTexture,
		swooshTexture: swooshTexture,
		stamina:       constants.CharacterStaminaMax,
		time:          0,
		facedRight:    true,
		swooshes:      []*swoosh{},
	}
	standingEnemyState := standingState{
		character:      &c,
		animationRects: standingEnemyRects,
	}
	walkingEnemyState := walkingState{
		character:      &c,
		animationRects: walkingEnemyRects,
	}
	jumpingEnemyState := jumpingState{
		character:      &c,
		animationRects: jumpingUpwardEnemyRects,
	}
	fallingEnemyState := fallingState{
		character:      &c,
		animationRects: fallingEnemyRects,
	}
	attackingEnemyState := attackingState{
		character:      &c,
		animationRects: attackingEnemyRects,
	}
	deadEnemyState := deadState{
		character:      &c,
		animationRects: deadEnemyRects,
	}
	c.standing = &standingEnemyState
	c.walking = &walkingEnemyState
	c.jumping = &jumpingEnemyState
	c.falling = &fallingEnemyState
	c.attacking = &attackingEnemyState
	c.dead = &deadEnemyState
	c.setState(c.falling)
	return &c
}

func (c *Character) Update(platforms []*platforms.Platform, enemies []*Character) {
	c.X += int32(c.vx)
	c.Y += int32(c.vy)
	if !c.CanAttack() {
		c.stamina++
	}
	c.currentState.update(platforms)
	c.updateAttack(enemies)
	c.swooshes = updateSwooshes(c.swooshes)
}

func (c *Character) updateAttack(enemies []*Character) {
	for _, s := range c.swooshes {
		for _, e := range enemies {
			if (s.x+s.w/2) > (e.X-e.W/2) && (s.x-s.w/2) < (e.X+e.W/2) {
				e.Kill(s.vx)
				s.destroyed = true
				break
			}
		}
	}
}

func (c *Character) CanAttack() bool {
	return c.stamina >= constants.CharacterStaminaMax
}

func (c *Character) reset() {
	c.X, c.Y = 0, 0
	c.vx, c.vy = 0, 0
}

func (c *Character) ResetVX() {
	c.vx = 0
}

func (c *Character) isTouchingFromAbove(p *platforms.Platform) bool {
	return c.Y+c.H >= p.Y-p.H/2 && c.Y+c.H <= p.Y-p.H/2+5 && c.X >= p.X-p.W/2 && c.X <= p.X+p.W/2
}

func (c *Character) isFalling() bool {
	return c.vy > 0
}

func (c *Character) isJumpingUpward() bool {
	return c.vy < 0
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

func (c *Character) Attack() {
	c.currentState.attack()
}

func (c *Character) Kill(newVX float32) {
	c.currentState.kill(newVX)
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
	// Draw swooshes made by character
	for _, s := range c.swooshes {
		s.draw(renderer)
	}
}
