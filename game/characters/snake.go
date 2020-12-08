package characters

import (
	"simpleplatformer/common"
	"simpleplatformer/constants"

	"github.com/veandco/go-sdl2/sdl"
)

func NewSnake(x, y int32, characterTexture *sdl.Texture) *Character {
	standingSnakeRects := newCharacterAnimationRects([]common.RelativeRectPosition{{0, 3}})
	walkingSnakeRects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{1, 3},
		{2, 3},
		{3, 3},
	})

	c := Character{
		X:             x,
		Y:             y,
		W:             constants.TileDestWidth,
		H:             constants.TileDestHeight,
		vx:            0,
		vy:            0,
		texture:       characterTexture,
		swooshTexture: nil,
		stamina:       constants.CharacterStaminaMax,
		time:          0,
		facedRight:    true,
		swooshes:      []*swoosh{},
		characterType: enemySnake,
	}
	standingSnakeState := standingState{
		character:      &c,
		animationRects: standingSnakeRects,
	}
	walkingSnakeState := walkingState{
		character:      &c,
		animationRects: walkingSnakeRects,
	}
	fallingSnakeState := fallingState{
		character:      &c,
		animationRects: standingSnakeRects,
	}
	hitSnakeState := hitState{
		character:      &c,
		animationRects: standingSnakeRects,
	}
	deadSnakeState := deadState{
		character:      &c,
		animationRects: standingSnakeRects,
	}
	c.updateAttack = func(enemies []*Character) {
		updateTouchAttack(&c, enemies)
	}
	c.standing = &standingSnakeState
	c.walking = &walkingSnakeState
	c.falling = &fallingSnakeState
	c.hit = &hitSnakeState
	c.dead = &deadSnakeState
	c.setState(c.falling)

	return &c
}

func updateTouchAttack(snake *Character, enemies []*Character) {
	s := *snake
	for _, e := range enemies {
		if snake == e || e.IsEnemySnake() {
			continue
		}
		if (s.Y+s.H/2) > (e.Y-e.H/2) && (s.Y-s.H/2) < (e.Y+e.H/2) { // Touches enemy vertically
			if (s.X+s.W/2) > (e.X-e.W/2) && (s.X-s.W/2) < (e.X+e.W/2) { // Touches enemy horizontally
				e.Hit(s.vx - e.vx)
				break
			}
		}
	}
}
