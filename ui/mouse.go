package ui

import (
	"github/teohen/mgm-tto/constants"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *UI) pan(dx, dy float32) {
	g.camera.Target.X -= dx / g.camera.Zoom
	g.camera.Target.Y -= dy / g.camera.Zoom
}

func (g *UI) zoom(factor float32) {
	newZoom := g.camera.Zoom * factor
	if newZoom < constants.CameraZoomMin {
		newZoom = constants.CameraZoomMin
	}
	if newZoom > constants.CameraZoomMax {
		newZoom = constants.CameraZoomMax
	}
	g.camera.Zoom = newZoom
}

func (ui *UI) selectCellsInBox() {
	ui.simulation.Selected = make(map[rl.Vector2]bool)

	minX := math.Min(float64(ui.dragStart.X), float64(ui.dragEnd.X))
	maxX := math.Max(float64(ui.dragStart.X), float64(ui.dragEnd.X))
	minY := math.Min(float64(ui.dragStart.Y), float64(ui.dragEnd.Y))
	maxY := math.Max(float64(ui.dragStart.Y), float64(ui.dragEnd.Y))

	if maxX-minX < 4 && maxY-minY < 4 {
		col := int(ui.dragEnd.X) / constants.TileSize
		row := int(ui.dragEnd.Y) / constants.TileSize
		if col >= 0 && col < constants.GridCols && row >= 0 && row < constants.GridRows {
			ui.simulation.Selected[rl.NewVector2(float32(col), float32(row))] = true
		}
		return
	}

	for row := 0; row < ui.simulation.World().Rows(); row++ {
		for col := 0; col < ui.simulation.World().Cols(); col++ {
			x, y := constants.WorldToScreen(col, row)
			tileCenterX := float64(x + constants.TileSize/2)
			tileCenterY := float64(y + constants.TileSize/2)

			if tileCenterX >= minX && tileCenterX <= maxX &&
				tileCenterY >= minY && tileCenterY <= maxY {
				ui.simulation.Selected[rl.NewVector2(float32(col), float32(row))] = true
			}
		}
	}
}

func handleMouse(ui *UI) {
	screenPos := rl.GetMousePosition()
	worldPos := screenToWorld(screenPos, ui.camera)

	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		delta := rl.GetMouseDelta()
		ui.pan(delta.X, delta.Y)
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		ui.zoom(1.0 + wheel*0.1)
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		ui.isDragging = true
		ui.dragStart = worldPos
		ui.dragEnd = worldPos
	}

	if rl.IsMouseButtonDown(rl.MouseLeftButton) && ui.isDragging {
		ui.dragEnd = worldPos
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) && ui.isDragging {
		ui.dragEnd = worldPos
		ui.isDragging = false
		ui.selectCellsInBox()
		ui.simulation.OnSelectionComplete()
	}
}

func screenToWorld(screenPos rl.Vector2, camera rl.Camera2D) rl.Vector2 {
	return rl.GetScreenToWorld2D(screenPos, camera)
}
