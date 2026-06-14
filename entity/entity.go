package entity

import (
	"github/teohen/mgm-tto/cnts"
)

type EntityType string

const (
	EntityTypeTree     EntityType = "tree"
	EntityTypeVillager EntityType = "villager"
)

type Entity interface {
	Tick()
	Draw()
	Pos() cnts.Point
	ID() string
	Type() EntityType
}
