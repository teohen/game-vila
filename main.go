package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github/teohen/mgm-tto/cnts"
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
	debug := flag.Bool("debug", false, "toggle debug info")
	flag.Parse()

	rl.InitWindow(cnts.ScreenW, cnts.ScreenH, "mgm-tto")
	rl.SetExitKey(rl.KeyEscape)
	rl.SetTargetFPS(60)

	spritebank.LoadAll()

	if *loadPath != "" {
		_, err := save.LoadFromFile(*loadPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		log.Fatal("LOAD FROM SAVE NOT IMPLEMENTED")
		return
	}

	if *debug == true {
		cnts.DEBUGGING = true
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
		g.UI.Input()
		g.Update()
		running = !rl.WindowShouldClose()

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		g.UI.Draw()
		rl.EndDrawing()
	}
}
