package game

import (
	"log"
	"simpleplatformer/constants"
	"simpleplatformer/game/characters"
	"simpleplatformer/game/platforms"
)

type patrollingStateInterface interface {
	update([]*platforms.Platform, *characters.Character)
	String() string // useful for debugging state
}

type patrollingStateMoveRight struct {
	ctrl *aiController
}

func (s *patrollingStateMoveRight) update(platforms []*platforms.Platform, playerCharacter *characters.Character) {
	showAlarmIfNoticedPlayer(s.ctrl, playerCharacter)
	ch := s.ctrl.character
	ch.Move(constants.CharacterVX)
	if ch.X > s.ctrl.startX+(3*constants.TileDestWidth) || ch.IsCloseToPlatformRightEdge(platforms) {
		s.ctrl.setState(s.ctrl.patrollingStand)
	}
}

func (s *patrollingStateMoveRight) String() string {
	return "patrollingStateMoveRight"
}

type patrollingStateMoveLeft struct {
	ctrl *aiController
}

func (s *patrollingStateMoveLeft) update(platforms []*platforms.Platform, playerCharacter *characters.Character) {
	showAlarmIfNoticedPlayer(s.ctrl, playerCharacter)
	ch := s.ctrl.character
	ch.Move(-constants.CharacterVX)
	if ch.X < s.ctrl.startX-(3*constants.TileDestWidth) || ch.IsCloseToPlatformLeftEdge(platforms) {
		s.ctrl.setState(s.ctrl.patrollingStand)
	}
}

func (s *patrollingStateMoveLeft) String() string {
	return "patrollingStateMoveLeft"
}

type patrollingStateStand struct {
	ctrl *aiController
}

func (s *patrollingStateStand) update(_ []*platforms.Platform, playerCharacter *characters.Character) {
	showAlarmIfNoticedPlayer(s.ctrl, playerCharacter)
	s.ctrl.time++
	s.ctrl.character.Move(0)
	if s.ctrl.time > 100 {
		s.ctrl.time = 0
		if s.ctrl.character.IsFacedRight() {
			s.ctrl.setState(s.ctrl.patrollingMoveLeft)
		} else {
			s.ctrl.setState(s.ctrl.patrollingMoveRight)
		}
	}
}

func (s *patrollingStateStand) String() string {
	return "patrollingStateStand"
}

type alarmedState struct {
	ctrl *aiController
}

func (s *alarmedState) update(_ []*platforms.Platform, _ *characters.Character) {
	s.ctrl.character.Move(0)
	s.ctrl.character.ShowAlarm()
	// If finished showing alarm
	if s.ctrl.character.CanAttack() {
		s.ctrl.setState(s.ctrl.patrollingStand)
	}
}

func (s *alarmedState) String() string {
	return "alarmedState"
}

func newAiController(ch *characters.Character) *aiController {
	ctrl := &aiController{
		character: ch,
		startX:    ch.X,
		time:      0,
	}
	ctrl.patrollingStand = &patrollingStateStand{ctrl}
	ctrl.patrollingMoveRight = &patrollingStateMoveRight{ctrl}
	ctrl.patrollingMoveLeft = &patrollingStateMoveLeft{ctrl}
	ctrl.alarmed = &alarmedState{ctrl}
	ctrl.setState(ctrl.patrollingMoveRight)
	return ctrl
}

type aiController struct {
	character *characters.Character
	startX    int32
	time      int

	currentPatrollingState patrollingStateInterface
	patrollingStand        patrollingStateInterface
	patrollingMoveRight    patrollingStateInterface
	patrollingMoveLeft     patrollingStateInterface
	alarmed                patrollingStateInterface
}

func (ai *aiController) setState(state patrollingStateInterface) {
	log.Println("switching to state:", state.String())
	ai.currentPatrollingState = state
}

func (ai *aiController) update(platforms []*platforms.Platform, playerCharacter *characters.Character) {
	ai.currentPatrollingState.update(platforms, playerCharacter)
}
