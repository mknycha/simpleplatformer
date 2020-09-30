package game

import (
	"simpleplatformer/constants"
	"simpleplatformer/game/characters"
)

type patrollingStateInterface interface {
	update()
	String() string // useful for debugging state
}

type patrollingStateMoveRight struct {
	ctrl *aiController
}

func (s *patrollingStateMoveRight) update() {
	ch := s.ctrl.character
	ch.Move(constants.CharacterVX)
	if ch.X > s.ctrl.startX+(3*constants.TileDestWidth) {
		s.ctrl.setState(s.ctrl.newPatrollingStateStand)
	}
}

func (s *patrollingStateMoveRight) String() string {
	return "patrollingStateMoveRight"
}

type patrollingStateMoveLeft struct {
	ctrl *aiController
}

func (s *patrollingStateMoveLeft) update() {
	ch := s.ctrl.character
	ch.Move(-constants.CharacterVX)
	if ch.X < s.ctrl.startX-(3*constants.TileDestWidth) {
		s.ctrl.setState(s.ctrl.newPatrollingStateStand)
	}
}

func (s *patrollingStateMoveLeft) String() string {
	return "patrollingStateMoveLeft"
}

type patrollingStateStand struct {
	ctrl *aiController
}

func (s *patrollingStateStand) update() {
	s.ctrl.time++
	s.ctrl.character.Move(0)
	if s.ctrl.time > 100 {
		s.ctrl.time = 0
		if s.ctrl.character.IsFacedRight() {
			s.ctrl.setState(s.ctrl.newPatrollingStateMoveLeft)
		} else {
			s.ctrl.setState(s.ctrl.newPatrollingStateMoveRight)
		}
	}
}

func (s *patrollingStateStand) String() string {
	return "patrollingStateStand"
}

func newAiController(ch *characters.Character) *aiController {
	ctrl := &aiController{
		character: ch,
		startX:    ch.X,
		time:      0,
	}
	ctrl.newPatrollingStateStand = &patrollingStateStand{ctrl}
	ctrl.newPatrollingStateMoveRight = &patrollingStateMoveRight{ctrl}
	ctrl.newPatrollingStateMoveLeft = &patrollingStateMoveLeft{ctrl}
	ctrl.setState(ctrl.newPatrollingStateMoveRight)
	return ctrl
}

type aiController struct {
	character *characters.Character
	startX    int32
	time      int

	newCurrentPatrollingState   patrollingStateInterface
	newPatrollingStateStand     patrollingStateInterface
	newPatrollingStateMoveRight patrollingStateInterface
	newPatrollingStateMoveLeft  patrollingStateInterface
}

func (ai *aiController) setState(state patrollingStateInterface) {
	ai.newCurrentPatrollingState = state
}

func (ai *aiController) update() {
	ai.newCurrentPatrollingState.update()
}
