package game

import (
	"math/rand"

	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/npc"
	"github/teohen/mgm-tto/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	sim      *simulation.Simulation
	Commands []string
	clock    Clock
}

func New() Game {
	sim := simulation.New(cnts.GridCols, 0)
	g := Game{
		sim:      sim,
		clock:    newClock(),
		Commands: make([]string, 0),
	}

	hit := false
	var villager *npc.Villager
	var storage *building.Storage
	var deer *npc.Deer
	for !hit {
		x := rand.Intn(cnts.GridCols)
		y := rand.Intn(cnts.GridRows)

		if !sim.World().IsOccupied(x, y) {
			villager = npc.NewVillager(x, y, sim.World(), sim.CollectResourceAt)
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
			deer = npc.NewDeer(x, y, sim.World())
			hit = true
		}
	}

	g.sim.AddVillager(villager)
	g.sim.AddStorage(storage)
	g.sim.AddDeer(deer)
	return g
}

func (g *Game) Update() *simulation.Simulation {
	dt := float64(rl.GetFrameTime()) * 1000.0
	ticks := g.clock.Advance(dt)

	for range ticks {
		g.sim.Tick()
	}

	return g.sim
}
