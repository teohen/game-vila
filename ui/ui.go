package ui

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/save"
	"github/teohen/mgm-tto/simulation"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UI struct {
	isDragging bool
	dragStart  rl.Vector2
	dragEnd    rl.Vector2
	console    Console
	simulation *simulation.Simulation
	camera     *rl.Camera2D
}

func newGameCamera() *rl.Camera2D {
	cam := rl.Camera2D{
		Target:   rl.NewVector2(float32(constants.GridCols)*constants.TileSize/2, float32(constants.GridRows)*constants.TileSize/2),
		Offset:   rl.NewVector2(constants.ScreenW/2, constants.ScreenH/2),
		Rotation: 0,
		Zoom:     1.0,
	}
	return &cam
}

func New(sim *simulation.Simulation) *UI {
	return &UI{
		simulation: sim,
		camera:     newGameCamera(),
	}
}

func NewFromSave(sim *simulation.Simulation, s save.Save) *UI {
	cam := newGameCamera()
	if s.Camera.Zoom != 0 {
		cam.Target.X = float32(s.Camera.TargetX)
		cam.Target.Y = float32(s.Camera.TargetY)
		cam.Zoom = float32(s.Camera.Zoom)
	}
	return &UI{
		simulation: sim,
		camera:     cam,
	}
}

func (ui *UI) Input() {
	ui.console.handleConsole(ui)
	if ui.console.IsOpen() {
		return
	}

	handleMouse(ui)
	handleKeyboard(ui)
}

func (ui *UI) Draw() {
	cam := *ui.camera
	rl.BeginMode2D(cam)
	ui.simulation.World().Draw()
	for _, e := range ui.simulation.Entities() {
		e.Draw()
	}

	drawSelectionRectangle(ui)
	drawSelectedCells(ui)
	rl.EndMode2D()
	if ui.console.IsOpen() {
		ui.console.DrawConsole()
	}
}

func drawSelectionRectangle(ui *UI) {
	if ui.isDragging {
		x := float32(math.Min(float64(ui.dragStart.X), float64(ui.dragEnd.X)))
		y := float32(math.Min(float64(ui.dragStart.Y), float64(ui.dragEnd.Y)))
		w := float32(math.Abs(float64(ui.dragEnd.X - ui.dragStart.X)))
		h := float32(math.Abs(float64(ui.dragEnd.Y - ui.dragStart.Y)))
		rl.DrawRectangleLines(int32(x), int32(y), int32(w), int32(h), rl.Blue)
	}
}

func drawSelectedCells(ui *UI) {
	for pos := range ui.simulation.Selected {
		col, row := pos[0], pos[1]
		x, y := constants.WorldToScreen(col, row)
		rl.DrawRectangleLines(
			int32(x), int32(y),
			int32(constants.TileSize), int32(constants.TileSize),
			rl.Red,
		)
	}
}
