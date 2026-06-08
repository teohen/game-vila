package agent

import (
	"fmt"
	"github/teohen/mgm-tto/goap"
	"strings"
)

type GoalStoreInventory struct {
	ID           string
	name         string
	desiredState *goap.State
	actions      []IAction
	Type         GoalType

	target Target
}

func NewGoalStoreInventory(desired string, t Target) IGoal {
	name := "StoreInventory"
	g := GoalStoreInventory{
		ID:           fmt.Sprintf("%s_%s", name, t.GetID()),
		name:         name,
		desiredState: goap.StateOf(strings.Split(desired, ",")...),
		target:       t,
		Type:         GoalStoreInventoryType,
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

func (gsi *GoalStoreInventory) GetID() string {
	return gsi.ID
}

func (gsi *GoalStoreInventory) GetType() GoalType {
	return GoalStoreInventoryType
}
