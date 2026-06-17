package entity

import (
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/building"
	"log"
)

type Storager struct {
	storage        *building.Storage
	Inventory      int
	MaxCarryWeight int
}

func NewStorager(max int) *Storager {
	return &Storager{
		Inventory:      0,
		MaxCarryWeight: max,
	}
}

func (sto *Storager) ExecuteAction(target agent.Target) bool {
	t, ok := target.(*building.Storage)
	if !ok {
		log.Fatal("storage convertion")
	}

	t.Insert(sto.Inventory)
	sto.Inventory = 0
	sto.storage = nil

	return true
}

func (sto *Storager) IncrementInventory(amount int) {
	sto.Inventory += amount
}

func (sto *Storager) CalculateWeight() int {
	return sto.Inventory * 2
}

func (sto *Storager) IsOverweighted() bool {
	return sto.CalculateWeight() >= sto.MaxCarryWeight
}
