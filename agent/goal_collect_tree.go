package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"log"
	"strings"

	"github.com/google/uuid"
)

type GoalCollectTree struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction
	target       Target
}

func NewGoalCollectTree(t Target) IGoal {
	name := "CollectTree"
	g := GoalCollectTree{
		id:           strings.ReplaceAll(uuid.NewString(), "-", ""),
		name:         name,
		desiredState: goap.StateOf(fmt.Sprintf("%s_health=0", t.ID())),
		target:       t,
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

func (g *GoalCollectTree) UpdateActions(t Target, desired *goap.State, w *world.World, pos cnts.Point) {
	newActions := make([]IAction, 0)
	for _, action := range g.actions {
		switch action.Type() {
		case ActionMoveType:
			newActions = append(newActions, NewActionMove("walkable", "near", t, w, pos))
		case ActionChopTreeType:
			newActions = append(newActions, NewActionChopTree("near", t))
		case ActionPutIntoType:
			newActions = append(newActions, NewActionPutInto("near", desired, t))
		}
	}
}

func (g *GoalCollectTree) IsRelevant(w *world.World, from cnts.Point, state *goap.State) bool {
	ok, err := state.Match(goap.StateOf("!overweighted"))
	if err != nil {
		log.Fatal("invalid state passed to match", err.Error())
	}
	return ok
}
