package agent

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
	"math"
	"strings"
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

func NewActionMove(r, o string, t Target, w *world.World, from cnts.Point) IAction {
	am := ActionMove{
		name:    "move_action",
		cost:    1,
		require: goap.StateOf(strings.Split(r, ",")...),
		outcome: goap.StateOf(strings.Split(o, ",")...),
		target:  t,
		world:   w,
		from:    from,
	}

	path := pathfinding.FindPath(am.world, am.from, am.target.Pos())
	if path == nil {
		am.outcome = goap.StateOf("!near")
		am.cost = math.MaxFloat32
	} else {
		am.cost = float32(len(path))
		am.outcome = goap.StateOf("near")
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
