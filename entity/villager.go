package entity

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type VillagerState string

const (
	StateVillagerIdle VillagerState = "idle"
	StateVillagerBusy VillagerState = "busy"
)

type Villager struct {
	IAgent
	ID    string
	State VillagerState
	//TODO: transform inventory into a Trait
	inventory      int
	maxCarryWeight int
	weight         int
}

func NewVillager(x, y int) *Villager {
	id := fmt.Sprintf("villager_%d_%d", x, y)
	v := Villager{}

	v.ID = id
	v.State = StateVillagerIdle
	v.IAgent = NewAgent(x, y, nil, v.incrementWood)
	v.maxCarryWeight = 100
	v.weight = 10

	return &v
}

// TODO: this goes to the inventory Trait later
func (v *Villager) incrementWood(amount int) {
	v.inventory += amount
	v.weight += amount * 5
}

func (v *Villager) Tick(entities *[]Entity, w *world.World) {
	v.IAgent.Movement().w = w
	// v.UpdateGoals(entities, w)
	switch v.State {
	case StateVillagerIdle:
		if found := v.ChooseGoal(); found {
			v.State = StateVillagerBusy
		}

		/*if job := GetJobQueue().Pop(); job != nil {
			switch job.Type {
			case JobTypeChopTrees:
				v.SetPlan(PlanTypeChopTrees, entities, getEntityFrom(job.TargetID, entities))
				v.State = StateVillagerBusy
			}
		}
		*/

	case StateVillagerBusy:
		finalAction := v.ExecuteAction()
		if finalAction {
			v.State = StateVillagerIdle
		}
	}

}

func (v *Villager) Pos() cnts.Point {
	return v.IAgent.Movement().Pos()
}

func (v *Villager) GetID() string {
	return v.ID
}

// TODO: move to Renderer
func getSource(v *Villager) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	x, y := cnts.WorldToScreen(v.Movement().pos.X, v.Movement().pos.Y)
	dst.X = x
	dst.Y = y
	dst.Width = cnts.TileSize
	dst.Height = cnts.TileSize
	src.X = 41
	src.Y = 21
	src.Width = 16
	src.Height = 19

	return src, dst
}

func (v *Villager) Draw() {
	src, dst := getSource(v)
	rl.DrawTexturePro(spritebank.Human, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (v *Villager) GetType() EntityType {
	return EntityTypeVillager
}

func getEntityFrom(id string, entities *[]Entity) Entity {
	for _, e := range *entities {
		if e.GetID() == id {
			return e
		}
	}
	return nil
}
