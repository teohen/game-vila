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

func TestVillagerDoesNotMoveWithoutTarget(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows:  3,
			Cols:  3,
			Cells: repeatGrid(3, 3, int(world.Empty)),
		},
		Villagers: []save.VillagerSave{
			{ID: "v1", Name: "test", Type: 1, X: 1, Y: 1},
		},
	})

	s.AdvanceTicks(10)

	x, y := s.Pos("v1")
	if x != 1 || y != 1 {
		t.Errorf("expected villager to stay at (1,1), moved to (%d,%d)", x, y)
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
