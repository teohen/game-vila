package entity

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type VillagerType int
type VillagerState string

const (
	Human VillagerType = 1
)

const (
	StateVillagerIdle VillagerState = "idle"
	StateVillagerBusy VillagerState = "busy"
)

type Plan struct {
	goal      *goap.State
	actions   []*Action
	TargetPos cnts.Point
}

type Villager struct {
	Agent
	ID    string
	name  string
	Type  VillagerType
	State VillagerState
}

func NewVillager(id, name string, x, y int) *Villager {
	v := &Villager{
		ID:    id,
		name:  name,
		Type:  Human,
		State: StateVillagerIdle,
		Agent: NewAgent(x, y, nil),
	}

	return v
}

func (v *Villager) Tick(entities *[]Entity, w *world.World) {
	v.Agent.Movement.w = w
	switch v.State {
	case StateVillagerIdle:
		if job := GetJobQueue().Pop(); job != nil {
			switch job.Type {
			case JobTypeChopTrees:
				v.SetPlan(PlanTypeChopTrees, entities, getEntityFrom(job.TargetID, entities))
				v.State = StateVillagerBusy
			}
		}

	case StateVillagerBusy:
		finalAction := v.Agent.ExecuteAction()
		if finalAction {
			v.State = StateVillagerIdle
		}
	}

}

func (v *Villager) Name() string {
	return v.name
}

func (v *Villager) Pos() cnts.Point {
	return v.Agent.Movement.Pos()
}

func (v *Villager) GetID() string {
	return v.ID
}

// TODO: move to Renderer
func getSource(v *Villager) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	switch v.Type {
	case Human:
		x, y := cnts.WorldToScreen(v.pos.X, v.pos.Y)
		dst.X = x
		dst.Y = y
		dst.Width = cnts.TileSize
		dst.Height = cnts.TileSize
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
