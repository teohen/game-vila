package agent

import (
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
)

type GoalStoreInventory struct {
	id           string
	name         string
	desiredState *goap.State
	actions      []IAction
	target       Target
	inventory    int
	w            *world.World
}

func NewGoalStoreInventory(w *world.World, t Target, from cnts.Point) IGoal {

	actions := []IAction{NewActionMove(w, t, from), NewActionPutInto(t)}

	g := GoalStoreInventory{
		name:         "StoreInventory",
		id:           cnts.NewID(),
		desiredState: goap.StateOf("inventory_incremented"),
		target:       t,
		w:            w,
		actions:      actions,
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

func (gsi *GoalStoreInventory) Actions() []IAction {
	return gsi.actions
}

func (gsi *GoalStoreInventory) IsRelevant(from cnts.Point, state *goap.State) bool {
	near := pathfinding.FindClosest(gsi.w, from, building.Get().GetBuildingsOf(building.StorageType))
	if near.X == -1 {
		return false
	}

	b := building.Get().GetBuildingAt(near)
	storage, ok := b.(*building.Storage)
	if !ok {
		return false
	}

	gsi.target = storage
	return true
}
