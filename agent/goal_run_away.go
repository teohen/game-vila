package agent

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
)

type GoalRunAway struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction
	target       Target
}

func NewGoalRunAway(w *world.World, t Target, from cnts.Point) IGoal {

	actions := []IAction{NewActionMove(w, t, from)}

	g := GoalRunAway{
		id:           cnts.NewID(),
		name:         "Roam",
		desiredState: goap.StateOf("!near"),
		target:       t,

		actions: actions,
	}

	return &g
}

func (g *GoalRunAway) DesiredState() *goap.State {
	return g.desiredState
}

func (g *GoalRunAway) Target() Target {
	return g.target
}

func (g *GoalRunAway) SetActions(a ...IAction) {
	g.actions = a
}

func (g *GoalRunAway) GetGoapActions() []goap.Action {
	a := make([]goap.Action, 0)
	for _, act := range g.actions {
		a = append(a, act)
	}
	return a
}
func (g *GoalRunAway) Name() string {
	return g.name
}

func (g *GoalRunAway) ID() string {
	return g.id
}

func (g *GoalRunAway) Type() GoalType {
	return GoalRunAwayType
}

func (g *GoalRunAway) Actions() []IAction {
	return g.actions
}

func (g *GoalRunAway) IsRelevant(from cnts.Point, state *goap.State) bool {
	match, err := state.Match(goap.StateOf("threatned"))
	if err != nil {
		return false
	}

	if !match {
		return false
	}

	return true
}
