package agent

import (
	"github/teohen/mgm-tto/goap"
)

type IGoal interface {
	DesiredState() *goap.State
	SetActions(a ...IAction)
	GetGoapActions() []goap.Action
	Target() Target
	Name() string
	ID() string
	Type() GoalType
	IsRelevant(state *goap.State) bool
}

type GoalType string

const (
	GoalCollectTreeType    GoalType = "CollectTree"
	GoalStoreInventoryType GoalType = "StoreInventory"
)
