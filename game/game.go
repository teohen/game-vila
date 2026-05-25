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

const debugPrintInterval = 60

type Game struct {
	sim    *simulation.Simulation
	UI     *ui.UI
	camera rl.Camera2D

	debugMode       bool
	debugFrameCount int

	clock Clock
}

func newGameCamera() rl.Camera2D {
	return rl.Camera2D{
		Target:   rl.NewVector2(float32(constants.GridCols)*constants.TileSize/2, float32(constants.GridRows)*constants.TileSize/2),
		Offset:   rl.NewVector2(constants.ScreenW/2, constants.ScreenH/2),
		Rotation: 0,
		Zoom:     1.0,
	}
}

func New() Game {
	sim := simulation.New()
	cam := newGameCamera()
	g := Game{
		sim:    sim,
		UI:     ui.New(sim, cam),
		camera: newGameCamera(),
		clock:  newClock(),
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
	cam := newGameCamera()
	if s.Camera.Zoom != 0 {
		cam.Target.X = float32(s.Camera.TargetX)
		cam.Target.Y = float32(s.Camera.TargetY)
		cam.Zoom = float32(s.Camera.Zoom)
	}

	sim := simulation.NewFromSave(s)
	ui := ui.New(sim, cam)
	return Game{
		sim:    sim,
		camera: cam,
		UI:     ui,
		clock:  newClock(),
	}
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
	s := g.sim.ToSave()
	s.Camera = save.CameraSave{
		TargetX: float64(g.camera.Target.X),
		TargetY: float64(g.camera.Target.Y),
		Zoom:    float64(g.camera.Zoom),
	}
	/*if err := save.SaveToFile(savePath, s); err != nil {
		fmt.Printf("[ERROR] Save failed: %v\n", err)
		return
	}
	*/
	savePath := save.GetSavePath()
	fmt.Printf("[SAVE] Game saved to %s\n", savePath)
}

func (g *Game) Load() {
	s, err := save.LoadFromFile(save.GetSavePath())
	if err != nil {
		fmt.Printf("[ERROR] Load failed: %v\n", err)
		return
	}

	g.sim = simulation.NewFromSave(s)
	g.camera.Target.X = float32(s.Camera.TargetX)
	g.camera.Target.Y = float32(s.Camera.TargetY)
	g.camera.Zoom = float32(s.Camera.Zoom)
	g.sim.Selected = make(map[rl.Vector2]bool)
	fmt.Printf("[SAVE] Game loaded from %s\n", save.GetSavePath())
}
