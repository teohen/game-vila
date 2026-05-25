package game_test

import (
	"testing"

	"github/teohen/mgm-tto/save"
	"github/teohen/mgm-tto/simulation"
)

func loadSave(t *testing.T, path string) *simulation.Simulation {
	t.Helper()
	s, err := save.LoadFromFile(path)
	if err != nil {
		t.Fatalf("failed to load save from %s: %v", path, err)
	}
	return simulation.NewFromSave(s)
}
