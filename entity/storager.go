package entity

import (
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/building"
	"log"
)

type Storager struct {
	storage        *building.Storage
	inventory      int
	maxCarryWeight int
	weight         int
}

func NewStorager(max int) *Storager {
	return &Storager{
		inventory:      0,
		maxCarryWeight: max,
		weight:         0,
	}
}

func (sto *Storager) ExecuteAction(target agent.Target) bool {
	t, ok := target.(*building.Storage)
	if !ok {
		log.Fatal("storage convertion")
	}

	t.Insert(sto.inventory)
	sto.inventory = 0
	sto.storage = nil
	sto.weight = 0

	return true
}

func (sto *Storager) IncrementWood(amount int) {
	sto.inventory += amount
	sto.weight += amount * 2
}

func (sto *Storager) isOverweighted() bool {
	return sto.weight >= sto.maxCarryWeight
}
