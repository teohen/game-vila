package agent

import (
	"github/teohen/mgm-tto/goap"
)

type IGoal interface {
	EvaluatePriority() int
	DesiredState() *goap.State
	AddActions(a ...IAction)
	GetGoapActions() []goap.Action
	Target() Target
	Name() string
	ID() string
	Type() GoalType
}

type GoalType string

const (
	GoalCollectTreeType    GoalType = "CollectTree"
	GoalStoreInventoryType GoalType = "StoreInventory"
)
