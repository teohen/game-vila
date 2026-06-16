package entity

import (
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/cnts"
)

type CollectResource func(pos cnts.Point) int
type AddToInventory func(amount int)

type Collecter struct {
	collectResource func(pos cnts.Point) int
	addToInventory  func(amount int)
}

func NewCollecter(colResource CollectResource, addToInventory AddToInventory) *Collecter {
	return &Collecter{
		collectResource: colResource,
		addToInventory:  addToInventory,
	}
}

// TODO: maybe develop a channel for the simulation to pass stuff down to the traits
// like in this case the i need to know the amount of wood that was just spawned
// in the simulation because of the lumberjack
func (c *Collecter) ExecuteAction(target agent.Target) bool {
	amout := c.collectResource(target.Pos())
	if amout > 0 {
		c.addToInventory(amout)
	}

	return true
}
