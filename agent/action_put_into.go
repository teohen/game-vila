package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
)

type ActionPutInto struct {
	name    string
	cost    int32
	require *goap.State
	outcome *goap.State
	target  Target
	world   *world.World
	from    cnts.Point
}

func NewActionPutInto(t Target) IAction {
	tpos := t.Pos()
	pi := ActionPutInto{
		name:    "put_into",
		cost:    1,
		require: goap.StateOf(fmt.Sprintf("at_%s", tpos.String())),
		outcome: goap.StateOf("inventory_stored"),
		target:  t,
	}
	return &pi
}

func (pi *ActionPutInto) Cost() float32 {
	return 1
}

func (pi *ActionPutInto) Simulate(current *goap.State) (*goap.State, *goap.State) {
	return pi.require, pi.outcome
}

func (pi *ActionPutInto) Name() string {
	return pi.name
}

func (pi *ActionPutInto) Target() Target {
	return pi.target
}

func (pi *ActionPutInto) SetTarget(e Target) {
	pi.target = e
}

func (pi *ActionPutInto) Type() ActionType {
	return ActionPutIntoType
}

func (pi *ActionPutInto) Update(t Target, from cnts.Point) {
	pi.from = from
	pi.target = t
}
