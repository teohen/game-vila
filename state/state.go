package state

type State struct {
	World     WorldDTO      `json:"world"`
	NPCs      []NPCDTO      `json:"npcs"`
	Resources []ResourceDTO `json:"resources"`
	Buildings []BuildingDTO `json:"buildings"`
	Inventory InventoryDTO  `json:"inventory"`
}

type WorldDTO struct {
	Occupied [][2]int `json:"occupied"`
}

type NPCDTO struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	State     string `json:"state"`
	Inventory *int   `json:"inventory,omitempty"`
}

type ResourceDTO struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Amount        int    `json:"amount"`
	IsCollectable bool   `json:"isCollectable"`
}

type BuildingDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type InventoryDTO struct {
	Wood  int `json:"wood"`
	Meat  int `json:"meat"`
	Hide  int `json:"hide"`
	Stone int `json:"stone"`
	Iron  int `json:"iron"`
}
