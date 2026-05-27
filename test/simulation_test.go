package game_test

import (
	"testing"

	"github/teohen/mgm-tto/entity"
)

func TestVillagerWalksToTarget(t *testing.T) {
	s := loadSave(t, "testdata/sim_walk_target.json")

	s.SetTarget("v1", 2, 2)
	s.AdvanceTicks(4)

	x, y := s.Pos("v1")
	if x != 2 || y != 2 {
		t.Errorf("expected villager at (2,2) after 4 ticks, got (%d,%d)", x, y)
	}
	if s.TickCount() != 4 {
		t.Errorf("expected tick count 4, got %d", s.TickCount())
	}
}

func TestVillagerArrivesInExactTicks(t *testing.T) {
	s := loadSave(t, "testdata/sim_arrive_exact.json")

	s.SetTarget("v1", 3, 3)
	s.AdvanceTicks(4)

	x, y := s.Pos("v1")
	if x != 3 || y != 3 {
		t.Errorf("expected (3,3), got (%d,%d) after %d ticks", x, y, s.TickCount())
	}
}

func TestMultipleVillagersMoveIndependently(t *testing.T) {
	s := loadSave(t, "testdata/sim_two_move.json")

	s.SetTarget("v1", 2, 0)
	s.SetTarget("v2", 4, 2)

	s.AdvanceTicks(2)

	x1, y1 := s.Pos("v1")
	if x1 != 2 || y1 != 0 {
		t.Errorf("expected v1 at (2,0) after 2 ticks, got (%d,%d)", x1, y1)
	}

	x2, y2 := s.Pos("v2")
	if x2 != 4 || y2 != 2 {
		t.Errorf("expected v2 at (4,2) after 2 ticks, got (%d,%d)", x2, y2)
	}
}

func TestVillagerDoesNotWalkOntoOccupiedCell(t *testing.T) {
	s := loadSave(t, "testdata/sim_occupied_cell.json")

	s.AdvanceTicks(6)

	vx, vy := s.Pos("v1")
	if vx == 0 && vy == 0 {
		t.Error("villager did not move")
	}
	if vx == 2 && vy == 2 {
		t.Errorf("villager walked onto tree's cell at (2,2); should have stopped adjacent")
	}
	if !s.World().IsOccupied(2, 2) {
		t.Error("tree cell (2,2) should still be occupied after villager arrives")
	}
}

func TestMultipleChopJobs_WalkedInSequence(t *testing.T) {
	s := loadSave(t, "testdata/sim_chop_sequence.json")

	s.AdvanceTicks(4)

	x, y := s.Pos("v1")
	if x != 3 || y != 0 {
		t.Fatalf("expected v1 at (3,0) after first job, got (%d,%d)", x, y)
	}

	s.AdvanceTicks(3)

	x, y = s.Pos("v1")
	if x != 6 || y != 0 {
		t.Errorf("expected v1 at (6,0) after both jobs, got (%d,%d)", x, y)
	}
}

func TestTwoVillagers_EachGetsAChopJob(t *testing.T) {
	s := loadSave(t, "testdata/sim_two_chop.json")

	s.AdvanceTicks(4)

	x1, y1 := s.Pos("v1")
	x2, y2 := s.Pos("v2")

	if x1 != 3 || y1 != 0 {
		t.Errorf("expected v1 at (3,0), got (%d,%d)", x1, y1)
	}
	if x2 != 6 || y2 != 0 {
		t.Errorf("expected v2 at (6,0), got (%d,%d)", x2, y2)
	}
}

func TestVillagerChopsTree_BasicFlow(t *testing.T) {
	s := loadSave(t, "testdata/lumberjack_basic.json")

	s.AdvanceTicks(4)

	x, y := s.Pos("v1")
	if x != 2 || y != 0 {
		t.Fatalf("expected v1 adjacent to tree at (2,0), got (%d,%d)", x, y)
	}

	tree := s.TreeAt(3, 0)
	if tree == nil {
		t.Fatal("tree should exist before chopping")
	}
	if tree.Health != 3 {
		t.Errorf("expected tree health 3 before chop, got %d", tree.Health)
	}

	if !s.World().IsOccupied(3, 0) {
		t.Error("tree cell should be occupied before chop")
	}

	s.AdvanceTicks(3)

	if s.TreeAt(3, 0) != nil {
		t.Error("tree should be removed from world after chop")
	}

	if s.World().IsOccupied(3, 0) {
		t.Error("tree cell should be vacated after destruction")
	}

	wood := s.VillagerWood("v1")
	if wood != 5 {
		t.Errorf("expected Wood=5 after chop, got %d", wood)
	}
}

func TestVillagerCarrying_BlocksNewJobs(t *testing.T) {
	s := loadSave(t, "testdata/lumberjack_basic.json")
	s.SetMaxCarryWeight("v1", 10)

	s.AdvanceTicks(7)

	wood := s.VillagerWood("v1")
	if wood != 5 {
		t.Fatalf("expected villager to have wood after chop, got %d", wood)
	}
	x0, y0 := s.Pos("v1")

	s.PushJob(entity.Job{Type: entity.JobTypeMove, TargetX: 9, TargetY: 9})
	s.AdvanceTicks(10)

	x, y := s.Pos("v1")
	if x != x0 || y != y0 {
		t.Errorf("expected v1 to stay at (%d,%d) while carrying, moved to (%d,%d)", x0, y0, x, y)
	}

	wood = s.VillagerWood("v1")
	if wood != 5 {
		t.Errorf("expected wood to remain 5 while carrying, got %d", wood)
	}
}

func TestTwoVillagers_ChopDifferentTrees(t *testing.T) {
	s := loadSave(t, "testdata/lumberjack_two_villagers.json")

	s.AdvanceTicks(4)

	x1, y1 := s.Pos("v1")
	if x1 != 2 || y1 != 0 {
		t.Errorf("expected v1 adjacent to tree-1 at (2,0), got (%d,%d)", x1, y1)
	}

	x2, y2 := s.Pos("v2")
	if x2 != 7 || y2 != 0 {
		t.Errorf("expected v2 adjacent to tree-2 at (7,0), got (%d,%d)", x2, y2)
	}

	if s.TreeAt(3, 0) == nil || s.TreeAt(6, 0) == nil {
		t.Fatal("both trees should exist before chopping")
	}

	s.AdvanceTicks(3)

	if s.TreeAt(3, 0) != nil {
		t.Error("tree-1 should be removed after chop")
	}
	if s.TreeAt(6, 0) != nil {
		t.Error("tree-2 should be removed after chop")
	}

	wood1 := s.VillagerWood("v1")
	if wood1 != 5 {
		t.Errorf("expected v1 Wood=5, got %d", wood1)
	}

	wood2 := s.VillagerWood("v2")
	if wood2 != 5 {
		t.Errorf("expected v2 Wood=5, got %d", wood2)
	}
}

func TestVillagerMultipleChopJobs_CarriesAfterFirst(t *testing.T) {
	s := loadSave(t, "testdata/lumberjack_two_villagers.json")
	s.SetMaxCarryWeight("v1", 10)

	s.AdvanceTicks(7)

	wood := s.VillagerWood("v1")
	if wood != 5 {
		t.Fatalf("expected v1 Wood=5 after chop, got %d", wood)
	}

	x0, y0 := s.Pos("v1")

	s.PushJob(entity.Job{Type: entity.JobTypeChopTrees, TargetX: 9, TargetY: 0})
	s.AdvanceTicks(10)

	x, y := s.Pos("v1")
	if x != x0 || y != y0 {
		t.Errorf("expected v1 to stay at (%d,%d) carrying wood, moved to (%d,%d)", x0, y0, x, y)
	}

	wood = s.VillagerWood("v1")
	if wood != 5 {
		t.Errorf("expected v1 wood to remain 5, got %d", wood)
	}
}

func TestVillagerSkipsNilTree(t *testing.T) {
	s := loadSave(t, "testdata/lumberjack_nil_tree.json")

	s.AdvanceTicks(5)

	x, y := s.Pos("v1")
	if x != 3 || y != 0 {
		t.Fatalf("expected v1 at (3,0) after walk + nil-tree skip, got (%d,%d)", x, y)
	}

	wood := s.VillagerWood("v1")
	if wood != 0 {
		t.Errorf("expected Wood=0 (no tree to chop), got %d", wood)
	}
}

func TestVillagerAcceptsJobs_WhenPartiallyCarrying(t *testing.T) {
	s := loadSave(t, "testdata/lumberjack_high_capacity.json")

	s.AdvanceTicks(7)

	wood := s.VillagerWood("v1")
	if wood != 5 {
		t.Fatalf("expected Wood=5 after chop, got %d", wood)
	}

	currentWeight := wood * entity.WoodWeightPerUnit
	if currentWeight >= s.VillagerMaxCarryWeight("v1") {
		t.Fatal("villager should have room to carry more")
	}

	s.PushJob(entity.Job{Type: entity.JobTypeMove, TargetX: 9, TargetY: 0})
	s.AdvanceTicks(8)

	x, y := s.Pos("v1")
	if x != 9 || y != 0 {
		t.Errorf("expected v1 to walk to (9,0) while partially carrying, got (%d,%d)", x, y)
	}
}

func TestVillagerSkipsNilTree_ReturnsToIdle(t *testing.T) {
	s := loadSave(t, "testdata/lumberjack_nil_tree.json")

	s.AdvanceTicks(4)

	s.PushJob(entity.Job{Type: entity.JobTypeMove, TargetX: 4, TargetY: 0})
	s.AdvanceTicks(3)

	x, y := s.Pos("v1")
	if x != 4 || y != 0 {
		t.Errorf("expected v1 at (4,0) after nil-tree skip + move, got (%d,%d)", x, y)
	}
}
