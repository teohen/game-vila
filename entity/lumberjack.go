package entity

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
}

func (lj *Lumberjack) Update() (woodCollected int, done bool) {
	if lj.state != StateLumberjackIdle || lj.tree == nil {
		return 0, false
	}

	lj.tree.Health -= lj.hit

	if lj.tree.Health <= 0 {
		// emit the event
		wood := lj.tree.WoodYield
		// w.Vacate(lj.tree.X, lj.tree.Y)
		// lj.tree.ID = ""
		// lj.state = StateLumberjackIdle
		return wood, true
	}

	return 0, false
}

func (lj *Lumberjack) IsHitting() bool {
	return lj.state == StateLumberjackHitting && lj.tree != nil
}
