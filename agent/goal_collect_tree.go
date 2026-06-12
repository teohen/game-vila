package agent

import (
	"github/teohen/mgm-tto/goap"
	"strings"

	"github.com/google/uuid"
)

type GoalCollectTree struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction

	target Target
}

func NewGoalCollectTree(desired string, t Target) IGoal {
	name := "CollectTree"
	g := GoalCollectTree{
		id:           strings.ReplaceAll(uuid.NewString(), "-", ""),
		name:         name,
		desiredState: goap.StateOf(strings.Split(desired, ",")...),
		target:       t,
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

func (g *GoalCollectTree) ID() string {
	return g.id
}

func (g *GoalCollectTree) Type() GoalType {
	return GoalCollectTreeType
}
