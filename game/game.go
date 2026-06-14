package game

import (
	"fmt"
	"math/rand"

	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/save"
	"github/teohen/mgm-tto/simulation"
	"github/teohen/mgm-tto/ui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	sim   *simulation.Simulation
	UI    *ui.UI
	clock Clock
}

func New() Game {
	sim := simulation.New()
	g := Game{
		sim:   sim,
		UI:    ui.New(sim),
		clock: newClock(),
	}

	hit := false
	var villager *entity.Villager
	var storage *building.Storage
	// var deer *entity.Deer
	for !hit {
		x := rand.Intn(cnts.GridCols)
		y := rand.Intn(cnts.GridRows)

		if !sim.World().IsOccupied(x, y) {
			villager = entity.NewVillager(x, y, sim.World())
			hit = true
		}
	}

	hit = false

	for !hit {
		x := rand.Intn(cnts.GridCols)
		y := rand.Intn(cnts.GridRows)

		if !sim.World().IsOccupied(x, y) {
			storage = building.NewStorage(x, y)
			hit = true
		}
	}

	hit = false

	for !hit {
		x := rand.Intn(cnts.GridCols)
		y := rand.Intn(cnts.GridRows)

		if !sim.World().IsOccupied(x, y) {
			// deer = entity.NewDeer(x, y, sim.World())
			hit = true
		}
	}

	g.sim.AddVillager(villager)
	g.sim.AddStorage(storage)
	// g.sim.AddDeer(deer)
	return g
}

func (g *Game) Update() {
	dt := float64(rl.GetFrameTime()) * 1000.0
	ticks := g.clock.Advance(dt)

	for i := 0; i < ticks; i++ {
		g.sim.Tick()
		g.UI.Draw()
	}
}

func (g *Game) Save() {
	fmt.Printf("[SAVE] Game saved to %s\n", save.GetSavePath())
}
