package entity

import (
	"fmt"
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type VillagerType int

const (
	Human VillagerType = 1
)

type Plan struct {
	goal    *goap.State
	actions []*Action
	TargetX *int
	TargetY *int
}

type Villager struct {
	Movement
	Lumberjack
	ID        string
	name      string
	Type      VillagerType
	Goals     []goap.State
	Actions   []goap.Action
	plan      Plan
	VilState  *goap.State
	State     string
	ActionIdx int
	World     *world.World
}

func NewVillager(id, name string, x, y int) *Villager {
	v := &Villager{
		Movement: Movement{
			X: x,
			Y: y,
		},
		ID:         id,
		name:       name,
		Type:       Human,
		Goals:      make([]goap.State, 0),
		Actions:    make([]goap.Action, 0),
		ActionIdx:  0,
		State:      "Idle",
		Lumberjack: Lumberjack{hit: 20, state: LumberjackIdle},
	}

	v.VilState = goap.StateOf("!near_tree")

	v.Actions = append(
		v.Actions,
		NewAction("move_to", "!near_tree", "near_tree"),
	)

	return v
}

func (v *Villager) Tick(entities *[]Entity, w *world.World) {
	v.World = w
	switch v.State {
	case "Idle":
		v.setPlan(entities)
	case "Planning":
		fmt.Println("PLANNING")
		v.State = "Executing"
	case "Executing":
		finalAction := v.executePlan(entities)
		if finalAction {
			v.State = "Idle"
		}
	}

}

func (v *Villager) setPlan(entities *[]Entity) {
	if len(GetJobQueue().jobs) > 0 {
		job := GetJobQueue().Pop()
		switch job.Type {
		case JobTypeChopTrees:
			t := getTreeFrom(job.TargetID, entities)
			v.VilState.Add(fmt.Sprintf("%s_health=%d", t.ID, t.Health))
			goapActions := append(v.Actions, NewAction("chop_tree", "near_tree", fmt.Sprintf("%s_health-20", t.ID)))
			goal := goap.StateOf(fmt.Sprintf("%s_health=0", t.ID))
			plan, err := goap.Plan(v.VilState, goal, goapActions)

			if err != nil {
				log.Fatal("ERRRO", err.Error())
			}
			finalActions := make([]*Action, 0)

			for _, act := range plan {
				action := act.(*Action)
				action.target = t
				finalActions = append(finalActions, action)

			}

			v.plan = Plan{
				goal:    goal,
				actions: finalActions,
				TargetX: &job.TargetX,
				TargetY: &job.TargetY,
			}
			v.State = "Planning"
		}
	}
}

func (v *Villager) executePlan(entities *[]Entity) bool {
	if v.ActionIdx >= len(v.plan.actions) {
		return true
	}
	v.executeAction(entities)
	return false
}

func (v *Villager) executeAction(entities *[]Entity) {
	action := v.plan.actions[v.ActionIdx]
	switch action.name {
	case "move_to":
		if v.Movement.TargetX < 1 && v.Movement.TargetY < 1 {
			v.Movement.SetTarget(*v.plan.TargetX, *v.plan.TargetY, v.World)
		} else {
			movState := v.Movement.Update(v.World)
			if movState == EventArrived {
				v.ActionIdx += 1
			}
		}
	case "chop_tree":
		t, ok := action.target.(*Tree)
		if !ok {
			log.Fatal("TREE CONVERTION NOT WOTK")
		}
		if !v.Lumberjack.IsHitting() {
			v.Lumberjack.Start(t)
		} else {
			_, done := v.Lumberjack.Update(v.World)
			if done {
				v.ActionIdx += 1
			}
		}
	}
}

func getTreeFrom(id string, entities *[]Entity) *Tree {
	for _, e := range *entities {
		tree, ok := e.(*Tree)
		if !ok {
			continue
		}

		if tree.ID == id {
			return tree
		}
	}
	return nil
}

func (v *Villager) Name() string {
	return v.name
}

func (v *Villager) Pos() (int, int) {
	return v.Movement.Pos()
}

func getSource(v *Villager) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	switch v.Type {
	case Human:
		x, y := constants.WorldToScreen(v.X, v.Y)
		dst.X = x
		dst.Y = y
		dst.Width = constants.TileSize
		dst.Height = constants.TileSize
		src.X = 41
		src.Y = 21
		src.Width = 16
		src.Height = 19
	}

	return src, dst
}

func (v *Villager) Draw() {
	src, dst := getSource(v)
	rl.DrawTexturePro(spritebank.Human, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}
