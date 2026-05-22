package main

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/game"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	running = true
	g       game.Game
)

func init() {
	rl.InitWindow(constants.ScreenW, constants.ScreenH, "mgm-tto")
	rl.SetExitKey(rl.KeyEscape)
	rl.SetTargetFPS(60)

	spritebank.LoadAll()
	g = game.New()
	villager := entity.NewVillager("teo", "teo", 10, 10)
	g.AddVillager(villager)

	tree := entity.NewTree("tree-1", 30, 11, 3, 5)
	g.AddTree(tree)
}

func quit() {
	spritebank.UnloadAll()
	rl.CloseWindow()
}

func main() {
	defer quit()

	for running {
		g.Input()
		g.Update()
		running = !rl.WindowShouldClose()

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		g.Draw()
		rl.EndDrawing()
	}
}
