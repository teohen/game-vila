package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
	"math"
)

type ActionMove struct {
	name    string
	cost    float32
	require *goap.State
	outcome *goap.State
	target  Target
	world   *world.World
	from    cnts.Point
}

func NewActionMove(w *world.World, t Target, from cnts.Point) IAction {
	am := ActionMove{
		name:    "move_action",
		cost:    1,
		require: goap.StateOf("walkable"),
		outcome: goap.StateOf(fmt.Sprintf("near_%s", t.ID())),
		target:  t,
		world:   w,
		from:    from,
	}
	return &am
}

func (am *ActionMove) Cost() float32 {
	return am.cost
}

func (am *ActionMove) Simulate(current *goap.State) (*goap.State, *goap.State) {
	return am.require, am.outcome
}

func (am *ActionMove) Name() string {
	return am.name
}

func (am *ActionMove) Target() Target {
	return am.target
}

func (am *ActionMove) SetTarget(e Target) {
	am.target = e
}

func (am *ActionMove) Type() ActionType {
	return ActionMoveType
}

func (am *ActionMove) Update(t Target, from cnts.Point) {
	am.target = t
	am.from = from

	path := pathfinding.FindPath(am.world, am.from, am.target.Pos())
	if path == nil {
		am.outcome = goap.StateOf(fmt.Sprintf("!near_%s", t.ID()))
		am.cost = math.MaxFloat32
	} else {
		am.cost = float32(len(path))
		am.outcome = goap.StateOf(fmt.Sprintf("near_%s", t.ID()))
	}
}
