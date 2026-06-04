package entity

import (
	"fmt"
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/spritebank"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type VillagerType int

const (
	Human VillagerType = 1
)

type Villager struct {
	Movement
	ID       string
	name     string
	Type     VillagerType
	Goals    []goap.State
	Actions  []goap.Action
	plan     []goap.Action
	VilState *goap.State
}

func NewVillager(id, name string, x, y int) *Villager {
	v := &Villager{
		Movement: Movement{
			X: x,
			Y: y,
		},
		ID:      id,
		name:    name,
		Type:    Human,
		Goals:   make([]goap.State, 0),
		Actions: make([]goap.Action, 0),
	}

	v.VilState = goap.StateOf("!near_tree")

	v.Actions = append(
		v.Actions,
		NewAction("move_to", "!near_tree", "near_tree"),
	)

	return v
}

func (v *Villager) Tick(entities []Entity) {
	if len(v.plan) < 1 {
		v.setPlan(entities)
	} else {
		for _, act := range v.plan {
			action := act.(*Action)
			// TODO: now the villager should act on the plan that is already set
		}
	}
}

func (v *Villager) setPlan(entities []Entity) {
	if len(GetJobQueue().jobs) > 0 {
		job := GetJobQueue().Pop()
		switch job.Type {
		case JobTypeChopTrees:
			t := getTreeFrom(job.TargetID, entities)
			v.VilState.Add(fmt.Sprintf("%s_health=%d", t.ID, t.Health))
			actions := append(v.Actions, NewAction("chopTree", "near_tree", fmt.Sprintf("%s_health-20", t.ID)))
			goal := goap.StateOf(fmt.Sprintf("%s_health=0", t.ID))
			plan, err := goap.Plan(v.VilState, goal, actions)

			if err != nil {
				log.Fatal("ERRRO", err.Error())
			}

			v.plan = plan
		}
	}
}

func getTreeFrom(id string, entities []Entity) *Tree {
	for _, e := range entities {
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
