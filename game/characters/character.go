package characters

import (
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"
	"simpleplatformer/game/ladders"
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
	move(float32)
	jump()
	attack()
	update([]*platforms.Platform, []*ladders.Ladder)
	hit(float32)
	kill(float32)
	showAlarm()
	climb(float32, []*ladders.Ladder)
	getAnimationRects() []*sdl.Rect
}

type standingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *standingState) move(newVX float32) {
	if newVX == 0 {
		return
	}
	c := s.character
	setVelocityAndSwitchFacedRight(c, newVX)
	c.setState(c.walking)
}

func (s *standingState) jump() {
	s.character.vy = -constants.JumpSpeed
	s.character.setState(s.character.jumping)
}

func (s *standingState) attack() {
	conditionalSwitchToAttackingState(s.character)
}

func (s *standingState) hit(newVX float32) {
	prepareAndSetHitState(s.character, newVX)
}

func (s *standingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *standingState) showAlarm() {
	prepareAndSetShowingAlarmState(s.character)
}

func (s *standingState) climb(newVY float32, lads []*ladders.Ladder) {
	conditionalClimbLadder(s.character, newVY, lads)
}

func (s *standingState) update([]*platforms.Platform, []*ladders.Ladder) {}

func (s *standingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type walkingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *walkingState) move(newVX float32) {
	setVelocityAndSwitchFacedRight(s.character, newVX)
}

func (s *walkingState) jump() {
	c := s.character
	c.vy = -constants.JumpSpeed
	c.setState(c.jumping)
}

func (s *walkingState) attack() {
	conditionalSwitchToAttackingState(s.character)
}

func (s *walkingState) hit(newVX float32) {
	prepareAndSetHitState(s.character, newVX)
}

func (s *walkingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *walkingState) showAlarm() {
	prepareAndSetShowingAlarmState(s.character)
}

func (s *walkingState) climb(newVY float32, lads []*ladders.Ladder) {
	conditionalClimbLadder(s.character, newVY, lads)
}

func (s *walkingState) update(platforms []*platforms.Platform, _ []*ladders.Ladder) {
	c := s.character
	c.time++
	for _, p := range platforms {
		// If character collides with ANY platform from above
		if c.isTouchingPlatformFromAbove(p) {
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

func (s *jumpingState) move(newVX float32) {
	setVelocityAndSwitchFacedRight(s.character, newVX)
}

func (s *jumpingState) jump() {}

func (s *jumpingState) attack() {}

func (s *jumpingState) hit(newVX float32) {
	prepareAndSetHitState(s.character, newVX)
}

func (s *jumpingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *jumpingState) showAlarm() {
	prepareAndSetShowingAlarmState(s.character)
}

func (s *jumpingState) climb(newVY float32, lads []*ladders.Ladder) {
	conditionalClimbLadder(s.character, newVY, lads)
}

func (s *jumpingState) update([]*platforms.Platform, []*ladders.Ladder) {
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

func (s *fallingState) move(newVX float32) {
	setVelocityAndSwitchFacedRight(s.character, newVX)
}

func (s *fallingState) jump() {}

func (s *fallingState) attack() {}

func (s *fallingState) hit(newVX float32) {
	prepareAndSetHitState(s.character, newVX)
}

func (s *fallingState) kill(newVX float32) {
	setVelocityAndSwitchToDeadState(s.character, newVX)
}

func (s *fallingState) showAlarm() {}

func (s *fallingState) climb(newVY float32, lads []*ladders.Ladder) {
	conditionalClimbLadder(s.character, newVY, lads)
}

func (s *fallingState) update(platforms []*platforms.Platform, _ []*ladders.Ladder) {
	c := s.character
	c.time = 0
	c.vy += constants.Gravity
	for _, p := range platforms {
		if c.isTouchingPlatformFromAbove(p) {
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

func (s *attackingState) move(float32) {}

func (s *attackingState) jump() {}

func (s *attackingState) attack() {}

func (s *attackingState) hit(newVX float32) {
	prepareAndSetHitState(s.character, newVX)
}

func (s *attackingState) kill(float32) {
	s.character.setState(s.character.dead)
}

func (s *attackingState) showAlarm() {}

func (s *attackingState) climb(newVX float32, lads []*ladders.Ladder) {}

func (s *attackingState) update(platforms []*platforms.Platform, _ []*ladders.Ladder) {
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

type hitState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *hitState) move(float32) {}

func (s *hitState) jump() {}

func (s *hitState) attack() {}

func (s *hitState) hit(float32) {}

func (s *hitState) kill(float32) {
	s.character.setState(s.character.dead)
}

func (s *hitState) showAlarm() {}

func (s *hitState) climb(newVX float32, lads []*ladders.Ladder) {}

func (s *hitState) update(platforms []*platforms.Platform, _ []*ladders.Ladder) {
	c := s.character
	if c.health <= 0 {
		c.setState(c.dead)
		return
	}
	c.time++
	c.vy += constants.Gravity
	for _, p := range platforms {
		if c.isTouchingPlatformFromAbove(p) {
			c.Y = p.Y - p.H/2 - c.H
			c.vy = 0
		}
	}
	if c.time > constants.HitStateLength {
		c.resetVX()
		c.setState(c.falling)
	}
}

func (s *hitState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type showingAlarmState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *showingAlarmState) move(float32) {}

func (s *showingAlarmState) jump() {}

func (s *showingAlarmState) attack() {}

func (s *showingAlarmState) hit(newVX float32) {
	prepareAndSetHitState(s.character, newVX)
}

func (s *showingAlarmState) kill(float32) {
	s.character.setState(s.character.dead)
}

func (s *showingAlarmState) showAlarm() {}

func (s *showingAlarmState) climb(newVX float32, lads []*ladders.Ladder) {}

func (s *showingAlarmState) update(platforms []*platforms.Platform, _ []*ladders.Ladder) {
	c := s.character
	c.vy += constants.Gravity
	for _, p := range platforms {
		if c.isTouchingPlatformFromAbove(p) {
			c.Y = p.Y - p.H/2 - c.H
			c.vy = 0
			c.setState(c.standing)
		}
	}
}

func (s *showingAlarmState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type climbingState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *climbingState) move(float32) {}

func (s *climbingState) jump() {
	c := s.character
	c.vy = 0
	c.setState(c.jumping)
}

func (s *climbingState) attack() {}

func (s *climbingState) hit(newVX float32) {
	prepareAndSetHitState(s.character, newVX)
}

func (s *climbingState) kill(float32) {
	s.character.setState(s.character.dead)
}

func (s *climbingState) showAlarm() {}

func (s *climbingState) climb(newVY float32, lads []*ladders.Ladder) {
	s.character.vy = newVY
}

func (s *climbingState) update(platforms []*platforms.Platform, ladders []*ladders.Ladder) {
	c := s.character
	c.vx = 0
	c.Y += int32(c.vy)
	if c.vy == 0 {
		c.time = 0
	} else {
		c.time++
	}
	for _, l := range ladders {
		if c.isTouchingLadder(l) {
			return
		}
	}
	for _, p := range platforms {
		if c.isTouchingPlatformFromAbove(p) {
			c.Y = p.Y - p.H/2 - c.H
			c.vy = 0
			c.setState(c.standing)
			return
		}
	}
	c.vy = 0
	c.setState(c.falling)
}

func (s *climbingState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type deadState struct {
	character      *Character
	animationRects []*sdl.Rect
}

func (s *deadState) move(float32) {}

func (s *deadState) jump() {}

func (s *deadState) attack() {}

func (s *deadState) kill(float32) {}

func (s *deadState) hit(float32) {}

func (s *deadState) showAlarm() {}

func (s *deadState) climb(newVX float32, lads []*ladders.Ladder) {}

func (s *deadState) update([]*platforms.Platform, []*ladders.Ladder) {
	c := s.character
	c.time++
	c.vy += constants.Gravity
}

func (s *deadState) getAnimationRects() []*sdl.Rect {
	return s.animationRects
}

type characterType int

const (
	player characterType = iota
	enemySlasher
	enemySnake
)

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
	health        int
	updateAttack  func([]*Character)
	characterType characterType

	standing     characterState
	walking      characterState
	jumping      characterState
	attacking    characterState
	hit          characterState
	dead         characterState
	falling      characterState
	climbing     characterState
	showingAlarm characterState
}

// IsPlayer returns true if the character is of player type
func (c *Character) IsPlayer() bool {
	return c.characterType == player
}

// IsEnemySlasher returns true if the character is of enemy slasher type
func (c *Character) IsEnemySlasher() bool {
	return c.characterType == enemySlasher
}

// IsEnemySnake returns true if the character is of enemy snake type
func (c *Character) IsEnemySnake() bool {
	return c.characterType == enemySnake
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
	hitPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{9, 1},
		{10, 1},
	})
	climbingPlayerRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{19, 1},
		{20, 1},
		{21, 1},
		{22, 1},
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
		health:        constants.DefaultPlayerHealth,
		time:          0,
		facedRight:    true,
		swooshes:      []*swoosh{},
		characterType: player,
	}
	c.updateAttack = func(enemies []*Character) {
		updateSwooshAttack(&c, enemies)
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
	hitPlayerState := hitState{
		character:      &c,
		animationRects: hitPlayerRects,
	}
	climbingPlayerState := climbingState{
		character:      &c,
		animationRects: climbingPlayerRects,
	}
	deadPlayerState := deadState{
		character:      &c,
		animationRects: hitPlayerRects,
	}
	c.standing = &standingPlayerState
	c.walking = &walkingPlayerState
	c.jumping = &jumpingPlayerState
	c.falling = &fallingPlayerState
	c.attacking = &attackingPlayerState
	c.hit = &hitPlayerState
	c.climbing = &climbingPlayerState
	c.dead = &deadPlayerState
	c.setState(c.falling)
	return &c
}

func NewEnemyCharacter(x, y int32, characterTexture *sdl.Texture, swooshTexture *sdl.Texture) *Character {
	standingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{{0, 0}})
	walkingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{1, 0},
		{2, 0},
		{3, 0},
		{4, 0},
	})
	jumpingUpwardEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{{6, 0}})
	fallingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{{7, 0}})
	attackingEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{12, 0},
		{11, 0},
		{12, 0},
		{13, 0},
	})
	hitEnemyRects := newCharacterAnimationRects([]common.RelativeRectPosition{
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
		health:        constants.DefaultEnemyHealth,
		time:          0,
		facedRight:    true,
		swooshes:      []*swoosh{},
		characterType: enemySlasher,
	}
	c.updateAttack = func(enemies []*Character) {
		updateSwooshAttack(&c, enemies)
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
	hitEnemyState := hitState{
		character:      &c,
		animationRects: hitEnemyRects,
	}
	deadEnemyState := deadState{
		character:      &c,
		animationRects: hitEnemyRects,
	}
	showingAlarmEnemyState := showingAlarmState{
		character:      &c,
		animationRects: standingEnemyRects,
	}
	c.standing = &standingEnemyState
	c.walking = &walkingEnemyState
	c.jumping = &jumpingEnemyState
	c.falling = &fallingEnemyState
	c.attacking = &attackingEnemyState
	c.hit = &hitEnemyState
	c.dead = &deadEnemyState
	c.showingAlarm = &showingAlarmEnemyState
	c.setState(c.falling)
	return &c
}

func (c *Character) Update(platforms []*platforms.Platform, ladders []*ladders.Ladder, enemies []*Character) {
	c.X += int32(c.vx)
	c.Y += int32(c.vy)
	if !c.CanAttack() {
		c.stamina++
	}
	c.currentState.update(platforms, ladders)
	c.updateAttack(enemies)
}

func updateSwooshAttack(c *Character, enemies []*Character) {
	for _, s := range c.swooshes {
		for _, e := range enemies {
			if c == e {
				continue
			}
			if (s.y+s.h/2) > (e.Y-e.H/2) && (s.y-s.h/2) < (e.Y+e.H/2) { // Touches enemy vertically
				if (s.x+s.w/2) > (e.X-e.W/2) && (s.x-s.w/2) < (e.X+e.W/2) { // Touches enemy horizontally
					e.Hit(s.vx)
					s.destroyed = true
					break
				}
			}
		}
	}
	c.swooshes = updateSwooshes(c.swooshes)
}

func (c *Character) CanAttack() bool {
	return c.stamina >= constants.CharacterStaminaMax
}

func (c *Character) reset() {
	c.X, c.Y = 0, 0
	c.vx, c.vy = 0, 0
}

func (c *Character) resetVX() {
	c.vx = 0
}

func (c *Character) isTouchingPlatformFromAbove(p *platforms.Platform) bool {
	return c.Y+c.H >= p.Y-p.H/2 && c.Y+c.H <= p.Y-p.H/2+5 && c.X >= p.X-p.W/2 && c.X <= p.X+p.W/2
}

func (c *Character) isTouchingLadder(l *ladders.Ladder) bool {
	// Additional c.H allows character to get on the platform that's on the same level as top of the ladder
	// This >= does not allow character to fall down the platform when climbing down the ladder
	return c.X > l.X-l.W/2 && c.X < l.X+l.W/2 && c.Y >= l.Y-l.H/2-c.H && c.Y+c.H <= l.Y+l.H/2
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

func (c *Character) IsCloseToUpperScreenEdge() bool {
	return c.Y < (constants.ScreenMarginHeight)
}

func (c *Character) IsCloseToLowerScreenEdge() bool {
	return c.Y+(constants.ScreenMarginHeight) > constants.WindowHeight
}

func (c *Character) IsFacedRight() bool {
	return c.facedRight
}

func (c *Character) IsCloseToPlatformLeftEdge(platforms []*platforms.Platform) bool {
	for _, p := range platforms {
		if c.isTouchingPlatformFromAbove(p) {
			return c.X < (p.X - p.W/2 + c.W/2)
		}
	}
	return false
}

func (c *Character) IsCloseToPlatformRightEdge(platforms []*platforms.Platform) bool {
	for _, p := range platforms {
		if c.isTouchingPlatformFromAbove(p) {
			return c.X > (p.X + p.W/2 - c.W/2)
		}
	}
	return false
}

func (c *Character) FinishedShowingAlarm() bool {
	return c.currentState == c.standing
}

// CharacterClose returns true if the other character is relatively close horizontally and on the same height
func (c *Character) CharacterClose(otherCharacter *Character) bool {
	if c.OnSameHeight(otherCharacter) {
		if otherCharacter.X > c.X-(constants.CharacterSightLimit) || otherCharacter.X < c.X+constants.CharacterSightLimit {
			return true
		}
	}
	return false
}

// OnSameHeight returns false if other character is tile lower or tile higher than the character
func (c *Character) OnSameHeight(otherCharacter *Character) bool {
	if (c.Y > otherCharacter.Y+constants.CharacterDestHeight) || (c.Y < otherCharacter.Y-constants.CharacterDestHeight) {
		return false
	}
	return true
}

// CharacterWithinAttackRange returns true if the other character is in range of the potential attack
func (c *Character) CharacterWithinAttackRange(otherCharacter *Character) bool {
	distance := otherCharacter.X - c.X
	if distance > 0 && distance < constants.CharacterDestWidth {
		return true
	}
	if distance < 0 && distance > -constants.CharacterDestWidth {
		return true
	}
	return false
}

// CharacterWithinSight returns true if the other character can be seen
func (c *Character) CharacterWithinSight(otherCharacter *Character) bool {
	if c.OnSameHeight(otherCharacter) {
		if c.IsFacedRight() {
			if otherCharacter.X > c.X && otherCharacter.X < c.X+constants.CharacterSightLimit {
				return true
			}
		} else if otherCharacter.X < c.X && otherCharacter.X > c.X-constants.CharacterSightLimit {
			return true
		}
	}
	return false
}

func (c *Character) Move(newVX float32) {
	c.currentState.move(newVX)
}

func (c *Character) Jump() {
	c.currentState.jump()
}

func (c *Character) Attack() {
	c.currentState.attack()
}

func (c *Character) Hit(newVX float32) {
	c.currentState.hit(newVX)
}

func (c *Character) Kill(newVX float32) {
	c.currentState.kill(newVX)
}

func (c *Character) ShowAlarm() {
	c.currentState.showAlarm()
}

func (c *Character) Climb(newVY float32, lads []*ladders.Ladder) {
	c.currentState.climb(newVY, lads)
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
