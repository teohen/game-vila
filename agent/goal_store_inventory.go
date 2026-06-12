package agent

import (
	"github/teohen/mgm-tto/goap"
	"strings"

	"github.com/google/uuid"
)

type GoalStoreInventory struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction

	target Target
}

func NewGoalStoreInventory(desired string, t Target) IGoal {
	name := "StoreInventory"
	g := GoalStoreInventory{
		id:           strings.ReplaceAll(uuid.NewString(), "-", ""),
		name:         name,
		desiredState: goap.StateOf(strings.Split(desired, ",")...),
		target:       t,
	}

	return &g
}

func (gsi *GoalStoreInventory) DesiredState() *goap.State {
	return gsi.desiredState
}

func (gsi *GoalStoreInventory) EvaluatePriority() int {
	return 1
}

func (gsi *GoalStoreInventory) Target() Target {
	return gsi.target
}

func (gsi *GoalStoreInventory) AddActions(a ...IAction) {
	gsi.actions = append(gsi.actions, a...)
}
func (gsi *GoalStoreInventory) GetGoapActions() []goap.Action {
	a := make([]goap.Action, 0)
	for _, act := range gsi.actions {
		a = append(a, act)
	}
	return a
}
func (gsi *GoalStoreInventory) Name() string {
	return gsi.name
}

func (gsi *GoalStoreInventory) ID() string {
	return gsi.id
}

func (gsi *GoalStoreInventory) Type() GoalType {
	return GoalStoreInventoryType
}
