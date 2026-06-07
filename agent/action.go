package agent

import (
	"github/teohen/mgm-tto/goap"
)

type IAction interface {
	Simulate(current *goap.State) (*goap.State, *goap.State)
	Cost() float32
	Name() string
	Target() Target
	SetTarget(Target)
	Type() ActionType
}

type ActionType string

const (
	ActionMoveType     ActionType = "ActionMove"
	ActionChopTreeType ActionType = "ChopTree"
	ActionPutIntoType  ActionType = "PutInto"
)
