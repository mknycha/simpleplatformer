package characters

func conditionalSwitchToAttackingState(c *Character) {
	if !c.CanAttack() {
		return
	}
	c.swooshes = append(c.swooshes, newSwooshForCharacter(c))
	c.setState(c.attacking)
}

func setVelocityAndSwitchToDeadState(c *Character, newVX float32) {
	c.vx = newVX
	c.vy = -2
	c.setState(c.dead)
}
