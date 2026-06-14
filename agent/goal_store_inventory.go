package agent

import (
	"fmt"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
	"strings"

	"github.com/google/uuid"
)

type GoalStoreInventory struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction
	target       Target
	inventory    int
}

func NewGoalStoreInventory(inventory int) IGoal {
	name := "StoreInventory"
	g := GoalStoreInventory{
		id:        strings.ReplaceAll(uuid.NewString(), "-", ""),
		name:      name,
		inventory: inventory,
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

func (gsi *GoalStoreInventory) SetActions(a ...IAction) {
	gsi.actions = a
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

func (g *GoalStoreInventory) UpdateActions(t Target, desired *goap.State, w *world.World, pos cnts.Point) {
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

func (gsi *GoalStoreInventory) IsRelevant(w *world.World, from cnts.Point, state *goap.State) bool {
	match, err := state.Match(goap.StateOf("overweighted"))
	if err != nil || !match {
		return false
	}

	near := pathfinding.FindClosest(w, from, building.Get().GetBuildingsOf(building.StorageType))
	if near.X == -1 {
		return false
	}

	b := building.Get().GetBuildingAt(near)
	storage, ok := b.(*building.Storage)
	if !ok {
		return false
	}

	gsi.target = storage
	gsi.desiredState = goap.StateOf(fmt.Sprintf("%s_wood=%d", storage.ID(), storage.Wood+gsi.inventory))
	return true
}
