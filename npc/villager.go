package npc

import (
	"fmt"
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/entity"
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
	StateVillagerOverWeighted VillagerState = "overweighted"
)

type Villager struct {
	agent      agent.IAgent
	movement   *entity.Movement
	lumberjack *entity.Lumberjack
	storager   *entity.Storager
	id         string
	State      VillagerState
	w          *world.World
}

func NewVillager(x, y int, w *world.World) *Villager {
	v := Villager{}
	id := fmt.Sprintf("villager_%d_%d", x, y)
	movement := entity.NewMovement(x, y, w)
	storager := entity.NewStorager(100)
	lumberjack := entity.NewLumberjack(storager.IncrementInventory)

	v.id = id
	v.State = StateVillagerIdle
	v.movement = movement
	v.storager = storager
	v.lumberjack = lumberjack
	v.agent = agent.NewAgent(x, y, w, "Villager", "walkable,!overweighted")
	v.agent.RegisterActor(agent.ActionMoveType, movement)
	v.agent.RegisterActor(agent.ActionChopTreeType, lumberjack)
	v.agent.RegisterActor(agent.ActionPutIntoType, storager)
	v.w = w
	return &v
}

func (v *Villager) Tick() {
	if v.State == StateVillagerIdle {
		v.AddCollectTreeGoal()
	}
	if v.storager.IsOverweighted() {
		v.AddStorageGoal()
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
	return v.movement.Pos()
}

func (v *Villager) ID() string {
	return v.id
}

func (v *Villager) Draw() {
	x, y := cnts.WorldToScreen(v.Pos().X, v.Pos().Y)
	src := rl.NewRectangle(41, 21, 16, 19)
	dst := rl.NewRectangle(x+5, y+5, cnts.TileSize-10, cnts.TileSize-10)

	rl.DrawTexturePro(spritebank.Human, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (v *Villager) Type() NPCType {
	return VillagerNPCType
}

func (v *Villager) AddStorageGoal() {
	if len(v.agent.GetGoalsOf(agent.GoalStoreInventoryType)) > 0 {
		return
	}

	near := pathfinding.FindClosest(v.w, v.Pos(), building.Get().GetBuildingsOf(building.StorageType))
	if near.X == -1 {
		return
	}

	b := building.Get().GetBuildingAt(near)
	storage, ok := b.(*building.Storage)
	if !ok {
		return
	}

	v.agent.AddGoal(agent.NewGoalStoreInventory(v.w, storage, v.Pos()))
	v.agent.SetState("overweighted")
}

func (v *Villager) AddCollectTreeGoal() {
	if len(job.GetJobQueue().Jobs) < 1 {
		return
	}

	targets := make([]agent.Target, 0)
	for _, j := range job.GetJobQueue().Jobs {
		targets = append(targets, j.GetObject())
	}

	closest := pathfinding.FindClosest(v.w, v.Pos(), targets)
	for _, j := range job.GetJobQueue().Jobs {
		if j.GetObject().Pos() == closest {
			v.agent.AddGoal(agent.NewGoalCollectTree(v.w, j.Object, v.Pos()))
			job.GetJobQueue().Remove(j.Name(), j.Object.ID())
		}
	}
}
