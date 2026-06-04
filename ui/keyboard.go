package ui

import (
	"fmt"
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
}
