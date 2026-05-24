package game_test

import (
	"testing"

	"github/teohen/mgm-tto/game"
	"github/teohen/mgm-tto/save"
)

func loadSave(t *testing.T, path string) *game.Simulation {
	t.Helper()
	s, err := save.LoadFromFile(path)
	if err != nil {
		t.Fatalf("failed to load save from %s: %v", path, err)
	}
	return game.NewSimulationFromSave(s)
}
