package entity

import (
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/events"
	"log"
)

const (
	LUMBERJACK_HIT = 20
)

type LumberjackState string

const (
	StateLumberjackIdle    LumberjackState = "idle"
	StateLumberjackHitting LumberjackState = "hitting"
)

type IncrementWood func(amount int)

type Lumberjack struct {
	State         LumberjackState
	tree          *Tree
	hit           int
	incrementWood IncrementWood
}

func NewLumberjack(iw IncrementWood) *Lumberjack {
	return &Lumberjack{
		incrementWood: iw,
		State:         StateLumberjackIdle,
		tree:          nil,
		hit:           LUMBERJACK_HIT,
	}
}

func (lj *Lumberjack) Start(tree *Tree) {
	lj.tree = tree
	lj.State = StateLumberjackHitting
}

func (lj *Lumberjack) Update() bool {

	lj.tree.Health -= lj.hit

	if lj.tree.Health <= 0 {
		lj.State = StateLumberjackIdle
		events.Emit(events.GameEvent{
			Type: events.EventTreeCut,
			Payload: map[string]interface{}{
				"pos": lj.tree.Pos(),
			},
		})

		lj.incrementWood(lj.tree.WoodYield)
		lj.tree = nil
		return true
	}

	return false
}

func (lj *Lumberjack) ExecuteAction(target agent.Target) bool {
	t, ok := target.(*Tree)
	if !ok {
		log.Fatal("TREE CONVERTION NOT WOTK")
	}
	if lj.State == StateLumberjackIdle {
		lj.Start(t)
	} else {
		if done := lj.Update(); done {
			return true
		}
	}
	return false
}
