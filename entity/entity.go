package entity

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/world"
)

type Entity interface {
	Tick(entities *[]Entity, w *world.World)
	Draw()
	Pos() cnts.Point
	GetID() string
}

type MovementEvent int

const (
	EventNone    MovementEvent = 0
	EventIdle    MovementEvent = 1
	EventArrived MovementEvent = 2
	EventStuck   MovementEvent = 3
)
