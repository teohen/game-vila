package entity

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
	"strings"
)

type IAction interface {
	Simulate(current *goap.State) (*goap.State, *goap.State)
	Cost() float32
	Name() string
	Target() Entity
}

type Action struct {
	name    string
	cost    int
	require *goap.State
	outcome *goap.State
	target  Entity
	world   *world.World
	from    cnts.Point
}

func NewAction(name, require, outcome string, t *Entity) *Action {
	return &Action{
		name:    name,
		require: goap.StateOf(strings.Split(require, ",")...),
		outcome: goap.StateOf(strings.Split(outcome, ",")...),
		target:  *t,
	}
}

func (a *Action) Simulate(current *goap.State) (*goap.State, *goap.State) {
	requre := a.require
	outcome := a.outcome
	if a.name == "move_to" {
		path := pathfinding.FindPath(a.world, a.from, a.Target().Pos())

		if path == nil {
			outcome = goap.StateOf("!near_tree")
			fmt.Println(a.require.String())
			fmt.Println(outcome.String())
		}
	}

	return requre, outcome
}

func (a *Action) Cost() float32 {
	return 1
}

func (a *Action) Name() string {
	return a.name
}

func (a *Action) Target() Entity {
	return a.target
}
