package main

import (
	"fmt"
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/game"
	"github/teohen/mgm-tto/spritebank"
	"math/rand"

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
	for i := 0; i < 10; i++ {
		x := rand.Intn(constants.GridCols)
		y := rand.Intn(constants.GridRows)
		villager := entity.NewVillager(fmt.Sprintf("teo-%s", i), fmt.Sprintf("teo-%s", i), x, y)
		g.AddVillager(villager)
	}
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
