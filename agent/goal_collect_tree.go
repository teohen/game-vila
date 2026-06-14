package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"log"
)

type GoalCollectTree struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction
	target       Target
}

func NewGoalCollectTree(w *world.World, t Target, from cnts.Point) IGoal {
	actions := []IAction{NewActionMove(w, t, from), NewActionChopTree(t)}
	g := GoalCollectTree{
		id:   cnts.NewID(),
		name: "CollectTree",
		// TODO: revise the action to not decrement but zero the health or just state the tree downed
		desiredState: goap.StateOf(fmt.Sprintf("%s_health=0", t.ID())),
		target:       t,
		actions:      actions,
	}

	return &g
}

func (g *GoalCollectTree) DesiredState() *goap.State {
	return g.desiredState
}

func (g *GoalCollectTree) Target() Target {
	return g.target
}

func (g *GoalCollectTree) SetActions(a ...IAction) {
	g.actions = a
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

func (g *GoalCollectTree) Actions() []IAction {
	return g.actions
}

func (g *GoalCollectTree) IsRelevant(from cnts.Point, state *goap.State) bool {
	ok, err := state.Match(goap.StateOf("!overweighted"))
	if err != nil {
		log.Fatal("invalid state passed to match", err.Error())
	}
	return ok
}
