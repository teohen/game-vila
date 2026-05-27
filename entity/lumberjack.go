package entity

import (
	"fmt"

	"github/teohen/mgm-tto/debug"
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
	state    LumberjackState
	tree     *Tree
	hitCount int
}

func (lj *Lumberjack) Start(tree *Tree) {
	lj.state = LumberjackHitting
	lj.tree = tree
	lj.hitCount = 0
}

func (lj *Lumberjack) Update(w *world.World) (woodCollected int, done bool) {
	if lj.state != LumberjackHitting || lj.tree == nil {
		return 0, false
	}

	lj.tree.Health--
	lj.hitCount++
	lj.debugLumberjack()

	if lj.tree.Health <= 0 {
		wood := lj.tree.WoodYield
		w.Vacate(lj.tree.X, lj.tree.Y)
		lj.tree = nil
		lj.state = LumberjackIdle
		lj.hitCount = 0
		return wood, true
	}

	return 0, false
}

func (lj *Lumberjack) IsHitting() bool {
	return lj.state == LumberjackHitting && lj.tree != nil
}

func (lj *Lumberjack) HitCount() int {
	return lj.hitCount
}

func (lj *Lumberjack) debugLumberjack() {
	if debug.IsEnabled(debug.Sim) {
		if lj.tree != nil {
			fmt.Printf("[LUMBERJACK] hit=%d tree=%s health=%d\n",
				lj.hitCount, lj.tree.ID, lj.tree.Health)
		}
	}
}
