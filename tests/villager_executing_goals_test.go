package tests

import (
	"github/teohen/mgm-tto/job"
	"github/teohen/mgm-tto/npc"
	"github/teohen/mgm-tto/resource"
	"github/teohen/mgm-tto/simulation"
	"testing"
)

func advanceTicks(sim *simulation.Simulation, n int) {
	for i := 0; i < n; i++ {
		sim.Tick()
	}
}

func TestFullCollectWoodGoal(t *testing.T) {
	sim := simulation.NewEmpty(5)
	vil := npc.NewVillager(0, 0, sim.World(), sim.CollectResourceAt)
	tree := resource.NewTree(3, 3, 100, 10)
	sim.AddVillager(vil)
	sim.AddTree(tree)
	job.GetJobQueue().Push(*job.NewJob(job.JobChopTreeType, tree))

	advanceTicks(sim, 15)
	if sim.Resources()[0].Type() != resource.ResourceWoodType {
		t.Fatalf("expected resources to be %T got = %T", resource.ResourceWoodType, sim.Resources()[0].Type())
	}

	advanceTicks(sim, 1)

	if len(sim.Resources()) != 0 {
		t.Fatalf("expected len(sim.Resources()) = %d, got=%d", 0, len(sim.Resources()))
	}

	if vil.Pos() != tree.Pos() {
		t.Error("expeted vil.Pos() to equal tree.Pos()")
		t.Errorf("vil: %v\n", vil.Pos())
		t.Errorf("tree: %v\n", tree.Pos())
		t.Fail()
	}

	if vil.Storager().Inventory != 10 {
		t.Fatalf("expected vil.Storager().Inventory = %d, got=%d", 10, vil.Storager().Inventory)
	}

}

func TestFullCollectTwoTreesSequentially(t *testing.T) {
	job.GetJobQueue().Jobs = job.GetJobQueue().Jobs[:0]

	sim := simulation.NewEmpty(10)
	vil := npc.NewVillager(0, 0, sim.World(), sim.CollectResourceAt)
	tree1 := resource.NewTree(2, 2, 100, 10)
	tree2 := resource.NewTree(4, 4, 100, 15)
	sim.AddVillager(vil)
	sim.AddTree(tree1)
	sim.AddTree(tree2)
	job.GetJobQueue().Push(*job.NewJob(job.JobChopTreeType, tree1))
	job.GetJobQueue().Push(*job.NewJob(job.JobChopTreeType, tree2))

	advanceTicks(sim, 13)

	foundWood := false
	for _, r := range sim.Resources() {
		if r.Type() == resource.ResourceWoodType {
			foundWood = true
			break
		}
	}
	if !foundWood {
		t.Fatalf("expected wood on ground after first tree chopped, resources: %v", sim.Resources())
	}

	advanceTicks(sim, 17)

	if len(sim.Resources()) != 0 {
		t.Fatalf("expected 0 resources on ground, got %d", len(sim.Resources()))
	}

	if vil.Pos() != tree2.Pos() {
		t.Errorf("expected vil.Pos()=%v, got %v", tree2.Pos(), vil.Pos())
	}

	expected := 25
	if vil.Storager().Inventory != expected {
		t.Fatalf("expected inventory=%d, got %d", expected, vil.Storager().Inventory)
	}
}
