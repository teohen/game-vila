package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"strings"
)

type ActionChopTree struct {
	name    string
	cost    int32
	require *goap.State
	outcome *goap.State
	target  Target
	world   *world.World
	from    cnts.Point
}

func NewActionChopTree(r string, t Target) IAction {
	cp := ActionChopTree{
		name:    "chop_tree",
		cost:    1,
		require: goap.StateOf(strings.Split(r, ",")...),
		outcome: goap.StateOf(fmt.Sprintf("%s_health=0", t.GetID())),
		target:  t,
	}
	return &cp
}

func (cp *ActionChopTree) Cost() float32 {
	return 1
}

func (cp *ActionChopTree) Simulate(current *goap.State) (*goap.State, *goap.State) {
	return cp.require, cp.outcome
}

func (cp *ActionChopTree) Name() string {
	return cp.name
}

func (cp *ActionChopTree) Target() Target {
	return cp.target
}

func (cp *ActionChopTree) SetTarget(e Target) {
	cp.target = e
}

func (cp *ActionChopTree) Type() ActionType {
	return ActionChopTreeType
}
