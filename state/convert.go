package state

import (
	"github/teohen/mgm-tto/npc"
	"github/teohen/mgm-tto/simulation"
)

func FromSimulation(sim *simulation.Simulation) State {
	w := sim.World()

	var occupied [][2]int
	for r := 0; r < w.Rows(); r++ {
		for c := 0; c < w.Cols(); c++ {
			if w.Occupied[r][c] {
				occupied = append(occupied, [2]int{c, r})
			}
		}
	}

	var npcsDTO []NPCDTO
	for _, n := range sim.NPCs() {
		dto := NPCDTO{
			ID:   n.ID(),
			Type: string(n.Type()),
			X:    n.Pos().X,
			Y:    n.Pos().Y,
		}
		switch v := n.(type) {
		case *npc.Villager:
			dto.State = string(v.State)
			inv := v.Storager().Inventory
			dto.Inventory = &inv
		case *npc.Deer:
			dto.State = string(v.State)
		}
		npcsDTO = append(npcsDTO, dto)
	}

	var resourcesDTO []ResourceDTO
	for _, r := range sim.Resources() {
		resourcesDTO = append(resourcesDTO, ResourceDTO{
			ID:            r.ID(),
			Type:          string(r.Type()),
			X:             r.Pos().X,
			Y:             r.Pos().Y,
			Amount:        r.Amount(),
			IsCollectable: r.Collectable(),
		})
	}

	var buildingsDTO []BuildingDTO
	for _, b := range sim.Buildings() {
		buildingsDTO = append(buildingsDTO, BuildingDTO{
			ID:   b.ID(),
			Type: string(b.Type()),
			X:    b.Pos().X,
			Y:    b.Pos().Y,
		})
	}

	inv := sim.Inventory
	inventoryDTO := InventoryDTO{}
	if inv != nil {
		inventoryDTO = InventoryDTO{
			Wood:  inv.Wood,
			Meat:  inv.Meat,
			Hide:  inv.Hide,
			Stone: inv.Stone,
			Iron:  inv.Iron,
		}
	}

	return State{
		World:     WorldDTO{Occupied: occupied},
		NPCs:      npcsDTO,
		Resources: resourcesDTO,
		Buildings: buildingsDTO,
		Inventory: inventoryDTO,
	}
}
