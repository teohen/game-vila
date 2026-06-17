package tests

import (
	"github/teohen/mgm-tto/agent"
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

	advanceTicks(sim, 3)

	if deer.Pos().X == 0 && deer.Pos().Y == 0 {
		t.Errorf("expected deer to have moved from starting position, got %v", deer.Pos())
	}
}

func TestDeerRunsAwayWhenCloseToVillager(t *testing.T) {
	sim := simulation.NewEmpty(5)

	deer := npc.NewDeer(1, 1, sim.World())
	vil := npc.NewVillager(1, 2, sim.World(), sim.CollectResourceAt)
	sim.AddDeer(deer)
	sim.AddVillager(vil)

	advanceTicks(sim, 1)

	goals := deer.GetGoals()
	foundRunAway := false
	for _, g := range goals {
		if g.Type() == agent.GoalRunAwayType {
			foundRunAway = true
			break
		}
	}
	if !foundRunAway {
		t.Errorf("expected deer to have a RunAway goal, got %d goals", len(goals))
	}

	advanceTicks(sim, 1)
	if deer.State != npc.StateRoaming {
		t.Errorf("expected deer to be roaming after choosing RunAway goal, got %s", deer.State)
	}

	advanceTicks(sim, 3)
	if deer.Pos().Y >= 1 {
		t.Errorf("expected deer to move away from villager (Y should decrease), got %v", deer.Pos())
	}
}
