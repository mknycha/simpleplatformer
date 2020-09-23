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
	deadSnakeState := deadState{
		character:      &c,
		animationRects: standingSnakeRects,
	}
	c.standing = &standingSnakeState
	c.walking = &walkingSnakeState
	c.falling = &fallingSnakeState
	c.dead = &deadSnakeState
	c.setState(c.falling)

	return &c
}
