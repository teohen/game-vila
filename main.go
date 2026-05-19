package main

import (
	"fmt"
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/grid"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	running = true
	g       grid.Grid

	npc entity.NPC
)

func drawScene() {
	g.Draw()

	npc.Draw()
}

func input() {
	mousePos := rl.GetMousePosition()

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		g.IsDragging = true
		g.DragStart = mousePos
		g.DragEnd = mousePos
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) && g.IsDragging {
		g.DragEnd = mousePos
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) && g.IsDragging {
		g.IsDragging = false
		g.SelectTilesInBox()
	}
}

func update() {
	running = !rl.WindowShouldClose()
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)
	drawScene()
	rl.EndDrawing()
}

func init() {
	rl.InitWindow(constants.ScreenW, constants.ScreenH, "mgm-tto")
	rl.SetExitKey(rl.KeyEscape)
	rl.SetTargetFPS(60)
	g = grid.NewGrid(constants.GridCols, constants.GridRows)
	npc = entity.New("teo", "teo", 10, 10)
}

func quit() {
	rl.CloseWindow()
}

func main() {
	defer quit()

	for running {
		input()
		update()
		render()
	}
}

func pr() {
	fmt.Println("")
}
