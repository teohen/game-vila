package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
)

type GoalRoam struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction
	target       Target
}

func NewGoalRoam(w *world.World, t Target, from cnts.Point) IGoal {

	actions := []IAction{NewActionMove(w, t, from)}

	g := GoalRoam{
		id:           cnts.NewID(),
		name:         "Roam",
		desiredState: goap.StateOf(fmt.Sprintf("near_%s", t.ID())),
		target:       t,

		actions: actions,
	}

	return &g
}

func (gr *GoalRoam) DesiredState() *goap.State {
	return gr.desiredState
}

func (gr *GoalRoam) Target() Target {
	return gr.target
}

func (gr *GoalRoam) SetActions(a ...IAction) {
	gr.actions = a
}

func (gr *GoalRoam) GetGoapActions() []goap.Action {
	a := make([]goap.Action, 0)
	for _, act := range gr.actions {
		a = append(a, act)
	}
	return a
}
func (gr *GoalRoam) Name() string {
	return gr.name
}

func (gr *GoalRoam) ID() string {
	return gr.id
}

func (gr *GoalRoam) Type() GoalType {
	return GoalRoamType
}

func (gr *GoalRoam) Actions() []IAction {
	return gr.actions
}

func (gr *GoalRoam) IsRelevant(from cnts.Point, state *goap.State) bool {
	match, err := state.Match(goap.StateOf("!threatned"))
	if err != nil {
		return false
	}

	if !match {
		return false
	}
	return true
}
