package save

type Save struct {
	Version   int           `json:"version"`
	World     WorldSave     `json:"world"`
	Villagers []VillagerSave `json:"villagers,omitempty"`
	Trees     []TreeSave    `json:"trees,omitempty"`
	Jobs      []JobSave     `json:"jobs,omitempty"`
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
