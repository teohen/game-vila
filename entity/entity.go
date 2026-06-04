package entity

type Entity interface {
	Tick(entities []Entity)
	Draw()
	Pos() (int, int)
}

type MovementEvent int

const (
	EventNone    MovementEvent = 0
	EventIdle    MovementEvent = 1
	EventArrived MovementEvent = 2
	EventStuck   MovementEvent = 3
)
