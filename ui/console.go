package ui

import (
	"fmt"
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/entity"
	"github/teohen/mgm-tto/simulation"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Console struct {
	open  bool
	input []rune
}

func (console *Console) handleConsole(ui *UI) {
	if rl.IsKeyPressed(rl.KeyGrave) {
		console.toggle()
	}

	if console.IsOpen() {
		cmd := console.readInput()
		if cmd != "" {
			ExecuteCommand(ui.simulation, cmd)
		}
		return
	}
}

func (c *Console) toggle() {
	c.open = !c.open
	if !c.open {
		c.input = c.input[:0]
	}
}

func (c *Console) IsOpen() bool {
	return c.open
}

func (c *Console) Input() string {
	return string(c.input)
}

func (c *Console) readInput() string {
	for ch := rl.GetCharPressed(); ch != 0; ch = rl.GetCharPressed() {
		if ch == '`' || ch == '~' {
			continue
		}
		c.input = append(c.input, ch)
	}

	if rl.IsKeyPressed(rl.KeyBackspace) && len(c.input) > 0 {
		c.input = c.input[:len(c.input)-1]
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		cmd := string(c.input)
		c.input = c.input[:0]
		return cmd
	}

	return ""
}

func ExecuteCommand(sim *simulation.Simulation, cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "help":
		fmt.Println("Commands: help, spawnvillager, addtree, removetree, cleartrees, addjob, clearjobs")

	case "spawnvillager":
		if len(parts) < 4 {
			fmt.Println("usage: spawnvillager <name> <x> <y>")
			return
		}
		x, err1 := strconv.Atoi(parts[2])
		y, err2 := strconv.Atoi(parts[3])
		if err1 != nil || err2 != nil {
			fmt.Println("invalid coordinates")
			return
		}
		sim.AddVillager(entity.NewVillager(parts[1], parts[1], x, y))
		fmt.Printf("spawned villager %s at (%d,%d)\n", parts[1], x, y)

	case "addtree":
		if len(parts) < 3 {
			fmt.Println("usage: addtree <x> <y>")
			return
		}
		x, err1 := strconv.Atoi(parts[1])
		y, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			fmt.Println("invalid coordinates")
			return
		}
		id := fmt.Sprintf("tree-%d-%d", x, y)
		sim.AddTree(entity.NewTree(id, x, y, 3, 5))
		fmt.Printf("added tree at (%d,%d)\n", x, y)

	case "removetree":
		if len(parts) < 3 {
			fmt.Println("usage: removetree <x> <y>")
			return
		}
		x, err1 := strconv.Atoi(parts[1])
		y, err2 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil {
			fmt.Println("invalid coordinates")
			return
		}
		if sim.RemoveTree(x, y) {
			fmt.Printf("removed tree at (%d,%d)\n", x, y)
		} else {
			fmt.Printf("no tree at (%d,%d)\n", x, y)
		}

	case "cleartrees":
		// sim.ClearTrees()
		fmt.Println("all trees removed")

	case "addjob":
		if len(parts) < 4 {
			fmt.Println("usage: addjob <move|chop> <x> <y>")
			return
		}
		x, err1 := strconv.Atoi(parts[2])
		y, err2 := strconv.Atoi(parts[3])
		if err1 != nil || err2 != nil {
			fmt.Println("invalid coordinates")
			return
		}
		var jt entity.JobType
		switch parts[1] {
		case "move":
			jt = entity.JobTypeMove
		case "chop":
			jt = entity.JobTypeChopTrees
		default:
			fmt.Println("unknown job type, use move or chop")
			return
		}
		sim.PushJob(entity.Job{Type: jt, TargetX: x, TargetY: y})
		fmt.Printf("added %s job at (%d,%d)\n", parts[1], x, y)

	case "clearjobs":
		// sim.ClearJobs()
		fmt.Println("all jobs cleared")

	default:
		fmt.Printf("unknown command: %s\n", parts[0])
	}
}

func (c *Console) DrawConsole() {
	const barH = 30
	barY := constants.ScreenH - barH
	rl.DrawRectangle(0, int32(barY), constants.ScreenW, barH, rl.NewColor(0, 0, 0, 200))
	line := "> " + c.Input() + "|"
	rl.DrawText(line, 8, int32(barY)+5, 18, rl.White)
}
