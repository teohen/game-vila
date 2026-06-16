package entity

import (
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/events"
	"github/teohen/mgm-tto/resource"
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

type Lumberjack struct {
	State LumberjackState
	tree  *resource.Tree
	hit   int
}

func NewLumberjack() *Lumberjack {
	return &Lumberjack{
		State: StateLumberjackIdle,
		tree:  nil,
		hit:   LUMBERJACK_HIT,
	}
}

func (lj *Lumberjack) Start(tree *resource.Tree) {
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
				"pos":       lj.tree.Pos(),
				"woodYield": lj.tree.WoodYield,
			},
		})
		lj.tree = nil
		return true
	}

	return false
}

func (lj *Lumberjack) ExecuteAction(target agent.Target) bool {
	t, ok := target.(*resource.Tree)
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
