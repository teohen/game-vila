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

type Villager struct {
	Movement
	ID   string
	name string
	Type VillagerType
}

func NewVillager(id, name string, x, y int) *Villager {
	return &Villager{
		Movement: Movement{
			X: x,
			Y: y,
		},
		ID:   id,
		name: name,
		Type: Human,
	}

}

func (v *Villager) Tick(w *world.World) MovementEvent {
	if len(GetJobQueue().jobs) > 0 {
		moveAction := NewAction("move_to", "!near_tree", "near_tree")
		job := GetJobQueue().Pop()
		chopTreeAction := NewAction("chopTree", "near_tree", fmt.Sprintf("%s_health-20", job.TargetID))

		goal := goap.StateOf(fmt.Sprintf("%s_health=0", job.TargetID))
		fmt.Println(job.WorldState.String())
		fmt.Println(goal.String())
		fmt.Println(moveAction.outcome.String())
		fmt.Println(chopTreeAction.outcome.String())

		plan, err := goap.Plan(job.WorldState, goal, []goap.Action{moveAction, chopTreeAction})
		if err != nil {
			log.Fatal("ERRRO", err.Error())
		}

		for i, action := range plan {
			fmt.Printf("%2d. %s\n", i+1, action.(*Action).String())
		}

	}
	return v.Movement.Update(w)
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
