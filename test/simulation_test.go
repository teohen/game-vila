package game_test

import (
	"testing"
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
