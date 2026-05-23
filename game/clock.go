package game

import (
	"fmt"

	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/debug"
)

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

	c.debugClock(fired, dtMs)
	return fired
}

func (c *Clock) TickCount() int {
	return c.tickCount
}

func (c *Clock) debugClock(fired int, dtMs float64) {
	if debug.IsEnabled(debug.Clock) && fired > 0 && c.tickCount%120 == 0 {
		fmt.Printf("[DEBUG] Clock dt=%.1fms fired=%d accum=%.1f tick=%d\n",
			dtMs, fired, c.accumulator, c.tickCount)
	}
}
