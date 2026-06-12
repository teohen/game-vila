package entity

import (
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/world"
)

type EntityType string

const (
	EntityTypeTree     EntityType = "tree"
	EntityTypeVillager EntityType = "villager"
)

type Entity interface {
	Tick(w *world.World, entities *[]Entity, buildings *building.BuildingsList)
	Draw()
	Pos() cnts.Point
	ID() string
	Type() EntityType
}
