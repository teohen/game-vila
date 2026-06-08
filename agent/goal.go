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

// func NewGoal(name, desired string, t Target) IGoal {
// 	g := Goal{
// 		ID:           strings.ReplaceAll(uuid.NewString(), "-", ""),
// 		name:         name,
// 		desiredState: goap.StateOf(strings.Split(desired, ",")...),
// 		target:       t,
// 	}

// 	return &g
// }

// func (g *Goal) DesiredState() *goap.State {
// 	return g.desiredState
// }

// func (g *Goal) EvaluatePriority() int {
// 	return 1
// }

// func (g *Goal) Target() Target {
// 	return g.target
// }

// func (g *Goal) AddActions(a ...IAction) {
// 	g.actions = append(g.actions, a...)
// }
// func (g *Goal) GetGoapActions() []goap.Action {
// 	a := make([]goap.Action, 0)
// 	for _, act := range g.actions {
// 		a = append(a, act)
// 	}
// 	return a
// }
// func (g *Goal) Name() string {
// 	return g.name
// }

// func (g *Goal) GetID() string {
// 	return g.ID
// }
