package main

import (
	"flag"
	"fmt"
	"os"

	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/game"
	"github/teohen/mgm-tto/save"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	running = true
	g       game.Game
)

func init() {
	loadPath := flag.String("load", "", "start from a save file")
	flag.Parse()

	rl.InitWindow(constants.ScreenW, constants.ScreenH, "mgm-tto")
	rl.SetExitKey(rl.KeyEscape)
	rl.SetTargetFPS(60)

	spritebank.LoadAll()

	if *loadPath != "" {
		s, err := save.LoadFromFile(*loadPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		g = game.NewFromSave(s)
		return
	}

	g = game.New()
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
