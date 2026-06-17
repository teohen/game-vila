package tests

import (
	"github/teohen/mgm-tto/npc"
	"github/teohen/mgm-tto/simulation"
	"testing"
)

func TestDeerAddsRoamAndExecutesItIfNotInThreat(t *testing.T) {
	sim := simulation.NewEmpty(5)
	deer := npc.NewDeer(0, 0, sim.World())
	sim.AddDeer(deer)

	advanceTicks(sim, 2)

	if deer.State != npc.StateRoaming {
		t.Errorf("expected deer to be roaming after choosing Roam goal, got %s", deer.State)
	}

	advanceTicks(sim, 8)

	if deer.Pos().X == 0 && deer.Pos().Y == 0 {
		t.Errorf("expected deer to have moved from starting position, got %v", deer.Pos())
	}
}
