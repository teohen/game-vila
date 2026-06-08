package agent

import (
	"fmt"
	"github/teohen/mgm-tto/goap"
	"strings"
)

type GoalCollectTree struct {
	ID           string
	name         string
	desiredState *goap.State
	actions      []IAction
	Type         GoalType

	target Target
}

func NewGoalCollectTree(desired string, t Target) IGoal {
	name := "CollectTree"
	g := GoalCollectTree{
		ID:           fmt.Sprintf("%s_%s", name, t.GetID()),
		name:         name,
		desiredState: goap.StateOf(strings.Split(desired, ",")...),
		target:       t,
		Type:         GoalCollectTreeType,
	}

	return &g
}

func (g *GoalCollectTree) DesiredState() *goap.State {
	return g.desiredState
}

func (g *GoalCollectTree) EvaluatePriority() int {
	return 1
}

func (g *GoalCollectTree) Target() Target {
	return g.target
}

func (g *GoalCollectTree) AddActions(a ...IAction) {
	g.actions = append(g.actions, a...)
}
func (g *GoalCollectTree) GetGoapActions() []goap.Action {
	a := make([]goap.Action, 0)
	for _, act := range g.actions {
		a = append(a, act)
	}
	return a
}
func (g *GoalCollectTree) Name() string {
	return g.name
}

func (g *GoalCollectTree) GetID() string {
	return g.ID
}

func (g *GoalCollectTree) GetType() GoalType {
	return GoalCollectTreeType
}
