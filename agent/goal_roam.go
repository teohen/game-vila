package agent

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
)

type GoalRoam struct {
	id             string
	name           string
	desiredState   *goap.State
	actions        []IAction
	target         Target
	allowedActions []string
}

func NewGoalRoam(t Target, actions []IAction) IGoal {
	name := "Roam"
	g := GoalRoam{
		id:           cnts.NewID(),
		name:         name,
		desiredState: goap.StateOf("near"),
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

func (gr *GoalRoam) UpdateActions(t Target, desired *goap.State, w *world.World, pos cnts.Point) {
	newActions := make([]IAction, 0)
	for _, action := range gr.actions {
		ac := action.Type()
		switch ac {
		case ActionMoveType:
			newActions = append(newActions, NewActionMove("walkable", "near", t, w, pos))
		case ActionChopTreeType:
			newActions = append(newActions, NewActionChopTree("near", t))
		case ActionPutIntoType:
			newActions = append(newActions, NewActionPutInto("near", desired, t))
		}
	}

	gr.actions = newActions

}

func (gr *GoalRoam) IsRelevant(w *world.World, from cnts.Point, state *goap.State) bool {
	return true
}
