package entity

import (
	"github/teohen/mgm-tto/world"
)

type LumberjackState int

const (
	LumberjackIdle    LumberjackState = 0
	LumberjackHitting LumberjackState = 1
)

func (s LumberjackState) String() string {
	switch s {
	case LumberjackIdle:
		return "idle"
	case LumberjackHitting:
		return "hitting"
	default:
		return "unknown"
	}
}

type Lumberjack struct {
	state LumberjackState
	tree  *Tree
	hit   int
}

func (lj *Lumberjack) Start(tree *Tree) {
	lj.state = LumberjackHitting
	lj.tree = tree
}

func (lj *Lumberjack) Update(w *world.World) (woodCollected int, done bool) {
	if lj.state != LumberjackHitting || lj.tree == nil {
		return 0, false
	}

	lj.tree.Health -= lj.hit

	if lj.tree.Health <= 0 {
		wood := lj.tree.WoodYield
		w.Vacate(lj.tree.X, lj.tree.Y)
		lj.tree.ID = ""
		lj.state = LumberjackIdle
		return wood, true
	}

	return 0, false
}

func (lj *Lumberjack) IsHitting() bool {
	return lj.state == LumberjackHitting && lj.tree != nil
}
