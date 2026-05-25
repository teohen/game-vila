package ui

import (
	"fmt"
	"github/teohen/mgm-tto/debug"
	"github/teohen/mgm-tto/simulation"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func handleKeyboard(ui *UI) {
	if rl.IsKeyPressed(rl.KeyOne) {
		if ui.simulation.ActiveTool == simulation.ToolAxe {
			ui.simulation.ActiveTool = simulation.ToolSelect
			fmt.Println("[TOOL] Selection")
		} else {
			ui.simulation.ActiveTool = simulation.ToolAxe
			fmt.Println("[TOOL] Axe")
		}
	}

	// if rl.IsKeyPressed(rl.KeyF9) {
	// 	g.Save()
	// }

	// if rl.IsKeyPressed(rl.KeyF10) {
	// 	g.Load()
	// }

	if rl.IsKeyPressed(rl.KeyF5) {
		debug.ToggleDebug()
		if debug.IsDebugActive() {
			fmt.Printf("[DEBUG] Debug mode enabled — %s\n", debug.ActiveString())
		} else {
			fmt.Printf("[DEBUG] Debug mode disabled\n")
		}
	}

	if debug.IsDebugActive() {
		switch {
		case rl.IsKeyPressed(rl.KeyZero):
			debug.Toggle(debug.Sim)
		case rl.IsKeyPressed(rl.KeyTwo):
			debug.Toggle(debug.Move)
		case rl.IsKeyPressed(rl.KeyThree):
			debug.Toggle(debug.Path)
		case rl.IsKeyPressed(rl.KeyFour):
			debug.Toggle(debug.Clock)
		case rl.IsKeyPressed(rl.KeyFive):
			debug.Toggle(debug.Job)
		case rl.IsKeyPressed(rl.KeySix):
			debug.Toggle(debug.World)
		}
	}
}
