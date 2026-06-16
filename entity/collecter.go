package entity

import (
	"github/teohen/mgm-tto/agent"
)

type Collecter struct {
}

func NewCollecter() *Collecter {
	return &Collecter{}
}

func (sto *Collecter) ExecuteAction(target agent.Target) bool {
	// t, ok := target.(*building.Storage)
	// if !ok {
	// 	log.Fatal("storage convertion")
	// }

	// t.Insert(sto.inventory)
	// sto.inventory = 0
	// sto.storage = nil
	// sto.weight = 0

	return true
}
