package game

import (
	"simpleplatformer/constants"
	"simpleplatformer/game/characters"
)

func showAlarmIfNoticedPlayer(ctrl *aiController, playerCharacter *characters.Character) {
	c := ctrl.character
	// If player is lower or higher than the character
	if (c.Y > playerCharacter.Y+constants.CharacterDestHeight) || (c.Y < playerCharacter.Y-constants.CharacterDestHeight) {
		return
	}
	if c.IsFacedRight() {
		if playerCharacter.X > c.X && playerCharacter.X < c.X+(8*constants.CharacterSourceWidth) {
			ctrl.setState(ctrl.alarmed)
		}
	} else if playerCharacter.X < c.X && playerCharacter.X > c.X-(8*constants.CharacterDestWidth) {
		ctrl.setState(ctrl.alarmed)
	}
}
