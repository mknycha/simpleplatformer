package characters

func conditionalSwitchToAttackingState(c *Character) {
	if !c.CanAttack() {
		return
	}
	c.swooshes = append(c.swooshes, newSwooshForCharacter(c))
	c.setState(c.attacking)
}
