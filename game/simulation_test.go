package game

import (
	"testing"

	"github/teohen/mgm-tto/save"
	"github/teohen/mgm-tto/world"
)

func TestVillagerWalksToTarget(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows:  3,
			Cols:  3,
			Cells: [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "test", Type: 1, X: 0, Y: 0},
		},
	})

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
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows:  5,
			Cols:  5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "test", Type: 1, X: 1, Y: 1},
		},
	})

	s.SetTarget("v1", 3, 3)
	s.AdvanceTicks(4)

	x, y := s.Pos("v1")
	if x != 3 || y != 3 {
		t.Errorf("expected (3,3), got (%d,%d) after %d ticks", x, y, s.TickCount())
	}
}

func TestVillagerMovesRandomlyWhenNoJobs(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows:  5,
			Cols:  5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "test", Type: 1, X: 2, Y: 2},
		},
	})

	s.AdvanceTicks(10)

	x, y := s.Pos("v1")
	if x == 2 && y == 2 {
		t.Errorf("expected villager to move to random target, but stayed at (2,2)")
	}
	if !s.World().IsWalkable(x, y) {
		t.Errorf("villager ended at non-walkable cell (%d,%d)", x, y)
	}
}

func TestMultipleVillagersMoveIndependently(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows:  5,
			Cols:  5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "alice", Type: 1, X: 0, Y: 0},
			{ID: "v2", Name: "bob", Type: 1, X: 4, Y: 4},
		},
	})

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
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows:  5,
			Cols:  5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "test", Type: 1, X: 0, Y: 0},
		},
		Trees: []save.TreeSave{
			{ID: "tree-1", X: 2, Y: 2, Health: 10, WoodYield: 5},
		},
		Jobs: []save.JobSave{
			{TargetX: 2, TargetY: 2},
		},
	})

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
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 10, Cols: 10,
			Cells: repeatGrid(10, 10, int(world.Empty)),
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "test", Type: 1, X: 0, Y: 0},
		},
		Jobs: []save.JobSave{
			{TargetX: 3, TargetY: 0},
			{TargetX: 6, TargetY: 0},
		},
	})

	// Tick 1: idle → pop job → set target (no movement)
	// Ticks 2-4: walk to (3,0) and arrive → pop second job → set target
	s.AdvanceTicks(4)

	x, y := s.Pos("v1")
	if x != 3 || y != 0 {
		t.Fatalf("expected v1 at (3,0) after first job, got (%d,%d)", x, y)
	}

	// Ticks 5-7: walk to (6,0) and arrive
	s.AdvanceTicks(3)

	x, y = s.Pos("v1")
	if x != 6 || y != 0 {
		t.Errorf("expected v1 at (6,0) after both jobs, got (%d,%d)", x, y)
	}
}

func TestTwoVillagers_EachGetsAChopJob(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 10, Cols: 10,
			Cells: repeatGrid(10, 10, int(world.Empty)),
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "alice", Type: 1, X: 0, Y: 0},
			{ID: "v2", Name: "bob", Type: 1, X: 9, Y: 0},
		},
		Jobs: []save.JobSave{
			{TargetX: 3, TargetY: 0},
			{TargetX: 6, TargetY: 0},
		},
	})

	// Tick 1: both idle → pop jobs → set targets (no movement)
	// Ticks 2-4: v1 walks 3 cells to (3,0), v2 walks 3 cells to (6,0)
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

func repeatGrid(rows, cols int, val int) [][]int {
	grid := make([][]int, rows)
	for r := range grid {
		grid[r] = make([]int, cols)
		for c := range grid[r] {
			grid[r][c] = val
		}
	}
	return grid
}
