package game

import (
	"errors"
	"log"
	"simpleplatformer/constants"
	"simpleplatformer/game/characters"
	"simpleplatformer/game/platforms"
)

type patrollingStateInterface interface {
	update([]*platforms.Platform, *characters.Character, []*characters.Character)
	String() string // useful for debugging state
}

type aiEnemyController interface {
	setState(state patrollingStateInterface)
	update([]*platforms.Platform, *characters.Character, []*characters.Character)
	shiftPatrollingReferencePointRight()
	shiftPatrollingReferencePointLeft()
}

type slasherPatrollingStateMoveRight struct {
	ctrl *aiEnemySlasherController
}

func (s *slasherPatrollingStateMoveRight) update(platforms []*platforms.Platform, playerCharacter *characters.Character, _ []*characters.Character) {
	showAlarmIfNoticedPlayer(s.ctrl, playerCharacter)
	ch := s.ctrl.character
	ch.Move(constants.CharacterVX)
	if ch.X > s.ctrl.startX+(3*constants.TileDestWidth) || ch.IsCloseToPlatformRightEdge(platforms) {
		s.ctrl.setState(s.ctrl.patrollingStand)
	}
}

func (s *slasherPatrollingStateMoveRight) String() string {
	return "patrollingStateMoveRight"
}

type slasherPatrollingStateMoveLeft struct {
	ctrl *aiEnemySlasherController
}

func (s *slasherPatrollingStateMoveLeft) update(platforms []*platforms.Platform, playerCharacter *characters.Character, _ []*characters.Character) {
	showAlarmIfNoticedPlayer(s.ctrl, playerCharacter)
	ch := s.ctrl.character
	ch.Move(-constants.CharacterVX)
	if ch.X < s.ctrl.startX-(3*constants.TileDestWidth) || ch.IsCloseToPlatformLeftEdge(platforms) {
		s.ctrl.setState(s.ctrl.patrollingStand)
	}
}

func (s *slasherPatrollingStateMoveLeft) String() string {
	return "patrollingStateMoveLeft"
}

type slasherPatrollingStateStand struct {
	ctrl *aiEnemySlasherController
}

func (s *slasherPatrollingStateStand) update(_ []*platforms.Platform, playerCharacter *characters.Character, _ []*characters.Character) {
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

func (s *slasherPatrollingStateStand) String() string {
	return "patrollingStateStand"
}

type slasherAlarmedState struct {
	ctrl *aiEnemySlasherController
}

func (s *slasherAlarmedState) update(_ []*platforms.Platform, _ *characters.Character, _ []*characters.Character) {
	s.ctrl.character.Move(0)
	s.ctrl.character.ShowAlarm()
	// If finished showing alarm
	if !s.ctrl.character.FinishedShowingAlarm() {
		s.ctrl.setState(s.ctrl.chasing)
	}
}

func (s *slasherAlarmedState) String() string {
	return "alarmedState"
}

type slasherChasingState struct {
	ctrl *aiEnemySlasherController
}

func (s *slasherChasingState) update(platforms []*platforms.Platform, playerCharacter *characters.Character, enemies []*characters.Character) {
	c := s.ctrl.character
	otherEnemyWithinAttackRange := false
	for _, e := range enemies {
		if e == c {
			continue
		}
		if c.CharacterWithinSight(e) && c.CharacterWithinAttackRange(e) {
			otherEnemyWithinAttackRange = true
			break
		}
	}
	if c.CharacterClose(playerCharacter) {
		s.ctrl.cooldownTime = constants.AiCooldownTime
		if c.CharacterWithinAttackRange(playerCharacter) {
			c.Attack()
		}
		if playerCharacter.X-constants.CharacterDestWidth/2 > c.X && !otherEnemyWithinAttackRange && !c.IsCloseToPlatformRightEdge(platforms) {
			c.Move(constants.CharacterVX)
		} else if playerCharacter.X+constants.CharacterDestWidth/2 < c.X && !otherEnemyWithinAttackRange && !c.IsCloseToPlatformLeftEdge(platforms) {
			c.Move(-constants.CharacterVX)
		} else {
			c.Move(0)
		}
	} else {
		if c.IsCloseToPlatformRightEdge(platforms) || c.IsCloseToPlatformLeftEdge(platforms) || otherEnemyWithinAttackRange {
			c.Move(0)
		}
		s.ctrl.cooldownTime--
		if s.ctrl.cooldownTime == 0 {
			s.ctrl.setState(s.ctrl.patrollingMoveLeft)
		}
	}
}

func (s *slasherChasingState) String() string {
	return "chasingState"
}

func newAiEnemySlasherController(ch *characters.Character) aiEnemyController {
	ctrl := &aiEnemySlasherController{
		character: ch,
		startX:    ch.X,
		time:      0,
	}
	ctrl.patrollingStand = &slasherPatrollingStateStand{ctrl}
	ctrl.patrollingMoveRight = &slasherPatrollingStateMoveRight{ctrl}
	ctrl.patrollingMoveLeft = &slasherPatrollingStateMoveLeft{ctrl}
	ctrl.alarmed = &slasherAlarmedState{ctrl}
	ctrl.chasing = &slasherChasingState{ctrl}
	ctrl.setState(ctrl.patrollingMoveRight)
	return ctrl
}

type aiEnemySlasherController struct {
	character    *characters.Character
	startX       int32
	time         int
	cooldownTime int

	currentPatrollingState patrollingStateInterface
	patrollingStand        patrollingStateInterface
	patrollingMoveRight    patrollingStateInterface
	patrollingMoveLeft     patrollingStateInterface
	alarmed                patrollingStateInterface
	chasing                patrollingStateInterface
}

func (ai *aiEnemySlasherController) setState(state patrollingStateInterface) {
	log.Println("switching to state:", state.String())
	ai.currentPatrollingState = state
}

func (ai *aiEnemySlasherController) update(platforms []*platforms.Platform, playerCharacter *characters.Character, enemies []*characters.Character) {
	ai.currentPatrollingState.update(platforms, playerCharacter, enemies)
}

func (ai *aiEnemySlasherController) shiftPatrollingReferencePointRight() {
	ai.startX++
}

func (ai *aiEnemySlasherController) shiftPatrollingReferencePointLeft() {
	ai.startX--
}

type snakePatrollingStateMoveRight struct {
	ctrl *aiEnemySnakeController
}

func (s *snakePatrollingStateMoveRight) update(platforms []*platforms.Platform, playerCharacter *characters.Character, _ []*characters.Character) {
	ch := s.ctrl.character
	ch.Move(constants.CharacterVX)
	if ch.X > s.ctrl.startX+(3*constants.TileDestWidth) || ch.IsCloseToPlatformRightEdge(platforms) {
		s.ctrl.setState(s.ctrl.patrollingStand)
	}
}

func (s *snakePatrollingStateMoveRight) String() string {
	return "patrollingStateMoveRight"
}

type snakePatrollingStateMoveLeft struct {
	ctrl *aiEnemySnakeController
}

func (s *snakePatrollingStateMoveLeft) update(platforms []*platforms.Platform, playerCharacter *characters.Character, _ []*characters.Character) {
	ch := s.ctrl.character
	ch.Move(-constants.CharacterVX)
	if ch.X < s.ctrl.startX-(3*constants.TileDestWidth) || ch.IsCloseToPlatformLeftEdge(platforms) {
		s.ctrl.setState(s.ctrl.patrollingStand)
	}
}

func (s *snakePatrollingStateMoveLeft) String() string {
	return "patrollingStateMoveLeft"
}

type snakePatrollingStateStand struct {
	ctrl *aiEnemySnakeController
}

func (s *snakePatrollingStateStand) update(_ []*platforms.Platform, playerCharacter *characters.Character, _ []*characters.Character) {
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

func (s *snakePatrollingStateStand) String() string {
	return "patrollingStateStand"
}

func newAiEnemySnakeController(ch *characters.Character) aiEnemyController {
	ctrl := &aiEnemySnakeController{
		character: ch,
		startX:    ch.X,
		time:      0,
	}
	ctrl.patrollingStand = &snakePatrollingStateStand{ctrl}
	ctrl.patrollingMoveRight = &snakePatrollingStateMoveRight{ctrl}
	ctrl.patrollingMoveLeft = &snakePatrollingStateMoveLeft{ctrl}
	ctrl.setState(ctrl.patrollingMoveRight)
	return ctrl
}

type aiEnemySnakeController struct {
	character    *characters.Character
	startX       int32
	time         int
	cooldownTime int

	currentPatrollingState patrollingStateInterface
	patrollingStand        patrollingStateInterface
	patrollingMoveRight    patrollingStateInterface
	patrollingMoveLeft     patrollingStateInterface
}

func (ai *aiEnemySnakeController) setState(state patrollingStateInterface) {
	log.Println("switching to state:", state.String())
	ai.currentPatrollingState = state
}

func (ai *aiEnemySnakeController) update(platforms []*platforms.Platform, playerCharacter *characters.Character, enemies []*characters.Character) {
	ai.currentPatrollingState.update(platforms, playerCharacter, enemies)
}

func (ai *aiEnemySnakeController) shiftPatrollingReferencePointRight() {
	ai.startX++
}

func (ai *aiEnemySnakeController) shiftPatrollingReferencePointLeft() {
	ai.startX--
}

func newAiControllerForEnemy(ch *characters.Character) (aiEnemyController, error) {
	switch {
	case ch.IsEnemySlasher():
		return newAiEnemySlasherController(ch), nil
	case ch.IsEnemySnake():
		return newAiEnemySnakeController(ch), nil
	}
	return nil, errors.New("unknown enemy type")
}
