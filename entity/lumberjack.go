package entity

import "github/teohen/mgm-tto/events"

const (
	LUMBERJACK_HIT = 20
)

type LumberjackState string

const (
	StateLumberjackIdle    LumberjackState = "idle"
	StateLumberjackHitting LumberjackState = "hitting"
)

type Lumberjack struct {
	state LumberjackState
	tree  *Tree
	hit   int
}

func NewLumberjack() Lumberjack {
	return Lumberjack{
		state: StateLumberjackIdle,
		tree:  nil,
		hit:   LUMBERJACK_HIT,
	}
}

func (lj *Lumberjack) Start(tree *Tree) {
	lj.tree = tree
	lj.state = StateLumberjackHitting
}

func (lj *Lumberjack) Update() (woodCollected int, done bool) {

	lj.tree.Health -= lj.hit

	if lj.tree.Health <= 0 {
		wood := lj.tree.WoodYield

		lj.state = StateLumberjackIdle
		events.Emit(events.GameEvent{
			Type: events.EventTreeCut,
			Payload: map[string]interface{}{
				"pos": lj.tree.Pos(),
			},
		})
		lj.tree = nil
		return wood, true
	}

	return 0, false
}

func (lj *Lumberjack) IsHitting() bool {
	return lj.state == StateLumberjackHitting && lj.tree != nil
}
