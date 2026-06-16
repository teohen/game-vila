package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
)

type ActionPickUp struct {
	name    string
	cost    int32
	require *goap.State
	outcome *goap.State
	target  Target
}

func NewActionPickUp(t Target) IAction {

	targetPos := t.Pos()

	cp := ActionPickUp{
		name:    "pick_up",
		cost:    1,
		require: goap.StateOf(fmt.Sprintf("at_%s", targetPos.String())),
		outcome: goap.StateOf("picked_up_%s", targetPos.String()),
		target:  t,
	}
	return &cp
}

func (cp *ActionPickUp) Cost() float32 {
	return 1
}

func (cp *ActionPickUp) Simulate(current *goap.State) (*goap.State, *goap.State) {
	return cp.require, cp.outcome
}

func (cp *ActionPickUp) Name() string {
	return cp.name
}

func (cp *ActionPickUp) Target() Target {
	return cp.target
}

func (cp *ActionPickUp) SetTarget(e Target) {
	cp.target = e
}

func (cp *ActionPickUp) Type() ActionType {
	return ActionPickUpType
}

func (cp *ActionPickUp) Update(t Target, from cnts.Point) {
	cp.target = t
}
