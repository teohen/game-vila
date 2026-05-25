package game

import (
	"fmt"
	"math/rand"

	"github/teohen/mgm-tto/constants"
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

	for i := 0; i < 1; i++ {
		x := rand.Intn(constants.GridCols)
		y := rand.Intn(constants.GridRows)
		villager := entity.NewVillager(fmt.Sprintf("teo-%d", i), fmt.Sprintf("teo-%d", i), x, y)
		g.sim.AddVillager(villager)
	}
	return g
}

func NewFromSave(s save.Save) Game {
	sim := simulation.NewFromSave(s)
	return Game{
		sim:   sim,
		UI:    ui.NewFromSave(sim, s),
		clock: newClock(),
	}
}

func (g *Game) Update() {
	dt := float64(rl.GetFrameTime()) * 1000.0
	ticks := g.clock.Advance(dt)

	for i := 0; i < ticks; i++ {
		g.sim.Tick()
	}
}

func (g *Game) Save() {
	/*if err := save.SaveToFile(save.GetSavePath(), g.sim.ToSave()); err != nil {
		fmt.Printf("[ERROR] Save failed: %v\n", err)
		return
	}
	*/
	fmt.Printf("[SAVE] Game saved to %s\n", save.GetSavePath())
}

func (g *Game) Load() {
	s, err := save.LoadFromFile(save.GetSavePath())
	if err != nil {
		fmt.Printf("[ERROR] Load failed: %v\n", err)
		return
	}

	g.sim = simulation.NewFromSave(s)
	g.sim.Selected = make(map[[2]int]bool)
	fmt.Printf("[SAVE] Game loaded from %s\n", save.GetSavePath())
}
