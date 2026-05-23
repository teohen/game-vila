package game

import (
	"testing"

	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/save"
	"github/teohen/mgm-tto/world"
)

func TestAxeSelection_SingleTree_CreatesOneJob(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 5, Cols: 5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Trees: []save.TreeSave{
			{ID: "tree-1", X: 2, Y: 2, Health: 3, WoodYield: 5},
		},
	})

	s.ProcessAxeSelection([][2]int{{2, 2}})

	jobs := s.QueueJobs()
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Type != entity.JobTypeChopTrees {
		t.Errorf("expected JobTypeChopTrees (%d), got %d", entity.JobTypeChopTrees, jobs[0].Type)
	}
	if jobs[0].TargetX != 2 || jobs[0].TargetY != 2 {
		t.Errorf("expected target (2,2), got (%d,%d)", jobs[0].TargetX, jobs[0].TargetY)
	}
}

func TestAxeSelection_MultipleTrees_CreatesJobsForEach(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 5, Cols: 5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Trees: []save.TreeSave{
			{ID: "tree-1", X: 1, Y: 1, Health: 3, WoodYield: 5},
			{ID: "tree-2", X: 3, Y: 3, Health: 3, WoodYield: 5},
			{ID: "tree-3", X: 4, Y: 1, Health: 3, WoodYield: 5},
		},
	})

	s.ProcessAxeSelection([][2]int{{1, 1}, {3, 3}, {4, 1}})

	jobs := s.QueueJobs()
	if len(jobs) != 3 {
		t.Fatalf("expected 3 jobs, got %d", len(jobs))
	}
	for _, job := range jobs {
		if job.Type != entity.JobTypeChopTrees {
			t.Errorf("expected JobTypeChopTrees, got %d", job.Type)
		}
	}
}

func TestAxeSelection_MixedSelection_SkipsEmptyCells(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 5, Cols: 5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Trees: []save.TreeSave{
			{ID: "tree-1", X: 2, Y: 2, Health: 3, WoodYield: 5},
		},
	})

	s.ProcessAxeSelection([][2]int{{2, 2}, {0, 0}, {4, 4}})

	jobs := s.QueueJobs()
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job (only tree cell), got %d", len(jobs))
	}
	if jobs[0].TargetX != 2 || jobs[0].TargetY != 2 {
		t.Errorf("expected target (2,2), got (%d,%d)", jobs[0].TargetX, jobs[0].TargetY)
	}
}

func TestAxeSelection_NoTrees_NoJobs(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 5, Cols: 5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
	})

	s.ProcessAxeSelection([][2]int{{0, 0}, {1, 1}})

	jobs := s.QueueJobs()
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestAxeSelection_EmptySelection_NoJobs(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 5, Cols: 5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Trees: []save.TreeSave{
			{ID: "tree-1", X: 2, Y: 2, Health: 3, WoodYield: 5},
		},
	})

	s.ProcessAxeSelection([][2]int{})

	jobs := s.QueueJobs()
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs for empty selection, got %d", len(jobs))
	}
}

func TestAxeSelection_AppendsToExistingQueue(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 5, Cols: 5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Trees: []save.TreeSave{
			{ID: "tree-1", X: 2, Y: 2, Health: 3, WoodYield: 5},
			{ID: "tree-2", X: 4, Y: 4, Health: 3, WoodYield: 5},
		},
		Jobs: []save.JobSave{
			{Type: 0, TargetX: 0, TargetY: 0},
			{Type: 0, TargetX: 1, TargetY: 1},
		},
	})

	s.ProcessAxeSelection([][2]int{{2, 2}, {4, 4}})

	jobs := s.QueueJobs()
	if len(jobs) != 4 {
		t.Fatalf("expected 4 jobs (2 existing + 2 new), got %d", len(jobs))
	}
	if jobs[2].TargetX != 2 || jobs[2].TargetY != 2 {
		t.Errorf("expected 3rd job target (2,2), got (%d,%d)", jobs[2].TargetX, jobs[2].TargetY)
	}
	if jobs[3].TargetX != 4 || jobs[3].TargetY != 4 {
		t.Errorf("expected 4th job target (4,4), got (%d,%d)", jobs[3].TargetX, jobs[3].TargetY)
	}
}

func TestAxeSelection_OutOfBoundsCell_NoJob(t *testing.T) {
	s := NewSimulationFromSave(save.Save{
		World: save.WorldSave{
			Rows: 5, Cols: 5,
			Cells: repeatGrid(5, 5, int(world.Empty)),
		},
		Trees: []save.TreeSave{
			{ID: "tree-1", X: 2, Y: 2, Health: 3, WoodYield: 5},
		},
	})

	s.ProcessAxeSelection([][2]int{{-1, -1}, {10, 10}})

	jobs := s.QueueJobs()
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs for out-of-bounds cells, got %d", len(jobs))
	}
}
