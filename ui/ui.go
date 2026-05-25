package ui

import (
	"github/teohen/mgm-tto/constants"
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
	camera     rl.Camera2D
}

func New(sim *simulation.Simulation, cam rl.Camera2D) *UI {
	ui := UI{
		simulation: sim,
		camera:     cam,
	}
	return &ui
}

func (ui *UI) Input() {
	ui.console.handleConsole(ui)
	handleMouse(ui)
	handleKeyboard(ui)
}

func (ui *UI) Draw() {
	rl.BeginMode2D(ui.camera)
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
		col, row := int(pos.X), int(pos.Y)
		x, y := constants.WorldToScreen(col, row)
		rl.DrawRectangleLines(
			int32(x), int32(y),
			int32(constants.TileSize), int32(constants.TileSize),
			rl.Red,
		)
	}
}
