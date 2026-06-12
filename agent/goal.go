package agent

import (
	"github/teohen/mgm-tto/goap"
)

type Goal struct {
	ID           string
	name         string
	desiredState *goap.State
	actions      []IAction

	target Target
}

type IGoal interface {
	EvaluatePriority() int
	DesiredState() *goap.State
	AddActions(a ...IAction)
	GetGoapActions() []goap.Action
	Target() Target
	Name() string
	GetID() string
	GetType() GoalType
}

type GoalType string

const (
	GoalCollectTreeType    GoalType = "CollectTree"
	GoalStoreInventoryType GoalType = "StoreInventory"
)
