package entity

import (
	"fmt"
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/job"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type VillagerState string

const (
	StateVillagerIdle         VillagerState = "idle"
	StateVillagerBusy         VillagerState = "busy"
	StateVillagerOverWeighted VillagerState = "pverweighted"
)

// TODO: create a villager Interface
type Villager struct {
	agent      agent.IAgent
	movement   *Movement
	lumberjack *Lumberjack
	storager   *Storager
	ID         string
	State      VillagerState
	w          *world.World
}

func NewVillager(x, y int, w *world.World) *Villager {
	v := Villager{}
	id := fmt.Sprintf("villager_%d_%d", x, y)
	movement := NewMovement(x, y, w)
	storager := NewStorager(100)
	lumberjack := NewLumberjack(storager.incrementWood)

	v.ID = id
	v.State = StateVillagerIdle
	v.movement = movement
	v.storager = storager
	v.lumberjack = lumberjack
	v.agent = agent.NewAgent(x, y, nil, movement, lumberjack, storager)
	v.w = w
	return &v
}

func (v *Villager) Tick(w *world.World, entities *[]Entity, buildings []*building.Storage) {
	if j := job.GetJobQueue().Pop(); j != nil && v.State == StateVillagerIdle {
		v.agent.AddGoal(agent.NewGoalCollectTree(fmt.Sprintf("%s_health=0", j.Object.GetID()), j.Object))
		job.GetJobQueue().Remove(j.Name(), j.Object.GetID())
	}

	// over := v.storager.isOverweighted()
	// if over {
	// 	// fmt.Println("is overweighted", v.storager.inventory)
	// } else {
	// 	// fmt.Println("is light", v.storager.weight)
	// }

	// if over {

	// 	storage := v.findNearestStorage(w, buildings)
	// 	hasGoal := false
	// 	for _, g := range v.agent.GetGoals() {
	// 		if g.GetType() == agent.GoalStoreInventoryType && g.Target().GetID() == storage.ID {
	// 			hasGoal = true
	// 		}
	// 	}

	// 	if !hasGoal {
	// 		desired := fmt.Sprintf("%s_wood=%d", storage.ID, (storage.Wood + v.storager.inventory))
	// 		v.agent.AddGoal(agent.NewGoalStoreInventory(desired, storage))
	// 	}
	// }

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

func (v *Villager) findNearestStorage(w *world.World, storages []*building.Storage) *building.Storage {
	var storage *building.Storage
	nearest := make([]cnts.Point, 10_000)
	for _, b := range storages {
		path := pathfinding.FindPath(w, v.movement.pos, b.Pos())
		if len(path) > 0 && len(path) < len(nearest) {
			storage = b
		}
	}
	return storage
}
