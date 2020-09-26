package characters

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
	c.vy = -2
	c.health--
	c.setState(c.hit)
}

func setVelocityAndSwitchToDeadState(c *Character, newVX float32) {
	c.vx = newVX
	c.vy = -2
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
