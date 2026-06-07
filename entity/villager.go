package entity

import (
	"fmt"
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/job"
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
	agent      agent.IAgent
	movement   *Movement
	lumberjack *Lumberjack
	ID         string
	State      VillagerState
	//TODO: transform inventory into a Trait
	inventory      int
	maxCarryWeight int
	weight         int
	w              *world.World
}

func NewVillager(x, y int, w *world.World) *Villager {
	v := Villager{}
	id := fmt.Sprintf("villager_%d_%d", x, y)
	movement := NewMovement(x, y, w)
	Lumberjack := NewLumberjack(v.incrementWood)

	v.ID = id
	v.State = StateVillagerIdle
	v.movement = movement
	v.lumberjack = Lumberjack
	v.agent = agent.NewAgent(x, y, nil, v.incrementWood, movement, Lumberjack)
	v.maxCarryWeight = 100
	v.weight = 10
	v.w = w

	return &v
}

// TODO: this goes to the inventory Trait later
func (v *Villager) incrementWood(amount int) {
	v.inventory += amount
	v.weight += amount * 5
}

func (v *Villager) Tick(entities *[]Entity, w *world.World) {
	if j := job.GetJobQueue().Pop(); j != nil && v.State == StateVillagerIdle {
		v.agent.UpdateGoals(v.w, v.movement.pos, j.Object, j.Name())
	}

	switch v.State {
	case StateVillagerIdle:
		if found := v.agent.ChooseGoal(v.w, v.Pos()); found {
			v.State = StateVillagerBusy
		}
	case StateVillagerBusy:
		if finalAction := v.agent.ExecutePlan(); finalAction {
			v.State = StateVillagerIdle
		}
	}

}

func (v *Villager) Pos() cnts.Point {
	return v.movement.pos
}

func (v *Villager) GetID() string {
	return v.ID
}

// TODO: move to Renderer
func getSource(v *Villager) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	x, y := cnts.WorldToScreen(v.Pos().X, v.Pos().Y)
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

func GetEntityFrom(id string, entities *[]Entity) Entity {
	for _, e := range *entities {
		if e.GetID() == id {
			return e
		}
	}
	return nil
}
