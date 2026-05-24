package game

import (
	"fmt"
	"strconv"
	"strings"

	"github/teohen/mgm-tto/entity"
)

func ExecuteCommand(sim *Simulation, cmd string) {
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
		sim.ClearTrees()
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
		sim.ClearJobs()
		fmt.Println("all jobs cleared")

	default:
		fmt.Printf("unknown command: %s\n", parts[0])
	}
}
