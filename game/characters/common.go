package characters

import (
	"simpleplatformer/constants"
	"simpleplatformer/game/ladders"
)

func conditionalSwitchToAttackingState(c *Character) {
	if !c.CanAttack() {
		return
	}
	c.swooshes = append(c.swooshes, newSwooshForCharacter(c))
	c.setState(c.attacking)
}

func prepareAndSetHitState(c *Character, newVX float32) {
	c.facedRight = true
	if newVX > 0 {
		c.facedRight = false
	}
	c.vx = newVX
	c.vy = constants.CharacterVYWhenHit
	c.health--
	c.setState(c.hit)
}

func prepareAndSetShowingAlarmState(c *Character) {
	c.vy = -2
	c.setState(c.showingAlarm)
}

func setVelocityAndSwitchToDeadState(c *Character, newVX float32) {
	c.vx = newVX
	c.vy = constants.CharacterVYWhenHit
	c.setState(c.dead)
}

func setVelocityAndSwitchFacedRight(c *Character, newVX float32) {
	if newVX > 0 {
		c.facedRight = true
	} else if newVX < 0 {
		c.facedRight = false
	}
	c.vx = newVX
}

func conditionalClimbLadder(c *Character, newVY float32, lads []*ladders.Ladder) {
	if newVY == 0 {
		return
	}
	for _, l := range lads {
		if c.isTouchingLadder(l) {
			c.X = l.X
			c.vy = newVY
			c.setState(c.climbing)
		}
	}
}
