package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Save struct {
	Version   int            `json:"version"`
	World     WorldSave      `json:"world"`
	Villagers []VillagerSave `json:"villagers,omitempty"`
	Trees     []TreeSave     `json:"trees,omitempty"`
	Jobs      []JobSave      `json:"jobs,omitempty"`
	Camera    CameraSave     `json:"camera"`
}

type CameraSave struct {
	TargetX float64 `json:"target_x"`
	TargetY float64 `json:"target_y"`
	Zoom    float64 `json:"zoom"`
}

type WorldSave struct {
	Rows  int     `json:"rows"`
	Cols  int     `json:"cols"`
	Cells [][]int `json:"cells"`
}

type VillagerSave struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	TargetX *int   `json:"target_x,omitempty"`
	TargetY *int   `json:"target_y,omitempty"`
	State   string `json:"state,omitempty"`
}

type TreeSave struct {
	ID        string `json:"id"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Health    int    `json:"health"`
	WoodYield int    `json:"wood_yield"`
}

type JobSave struct {
	Type    int `json:"type"`
	TargetX int `json:"target_x"`
	TargetY int `json:"target_y"`
}

var savePath = "test/testdata/sim_chop_sequence.json"

func SaveToFile(path string, s Save) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("save: create directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("save: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("save: write file: %w", err)
	}

	return nil
}

func LoadFromFile(path string) (Save, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Save{}, fmt.Errorf("save: read file: %w", err)
	}

	var s Save
	if err := json.Unmarshal(data, &s); err != nil {
		return Save{}, fmt.Errorf("save: unmarshal: %w", err)
	}

	savePath = path

	return s, nil
}

func GetSavePath() string {
	return savePath
}
