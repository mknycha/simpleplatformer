package game

import (
	"simpleplatformer/game/characters"
)

func showAlarmIfNoticedPlayer(ctrl *aiEnemySlasherController, playerCharacter *characters.Character) {
	c := ctrl.character
	if c.CharacterWithinSight(playerCharacter) {
		ctrl.setState(ctrl.alarmed)
	}
}
