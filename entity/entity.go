package entity

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/world"
)

type EntityType string

const (
	EntityTypeTree     EntityType = "tree"
	EntityTypeVillager EntityType = "villager"
)

type Entity interface {
	Tick(entities *[]Entity, w *world.World)
	Draw()
	Pos() cnts.Point
	GetID() string
	GetType() EntityType
}
