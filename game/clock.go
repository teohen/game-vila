package game

import "github/teohen/mgm-tto/constants"

type Clock struct {
	accumulator float64
	tickCount   int
	interval    float64
}

func newClock() Clock {
	return Clock{
		interval: float64(constants.TickInterval),
	}
}

func (c *Clock) Advance(dtMs float64) int {
	c.accumulator += dtMs
	fired := 0
	for c.accumulator >= c.interval {
		c.accumulator -= c.interval
		c.tickCount++
		fired++
	}
	return fired
}

func (c *Clock) TickCount() int {
	return c.tickCount
}
