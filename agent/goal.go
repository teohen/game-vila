package agent

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
)

type IGoal interface {
	DesiredState() *goap.State
	SetActions(a ...IAction)
	GetGoapActions() []goap.Action
	Target() Target
	Name() string
	ID() string
	Type() GoalType
	IsRelevant(w *world.World, from cnts.Point, state *goap.State) bool
}

type GoalType string

const (
	GoalCollectTreeType    GoalType = "CollectTree"
	GoalStoreInventoryType GoalType = "StoreInventory"
)
