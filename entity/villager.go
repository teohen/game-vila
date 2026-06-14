package entity

import (
	"fmt"
	"github/teohen/mgm-tto/agent"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
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
	movement   *Movement
	lumberjack *Lumberjack
	storager   *Storager
	id         string
	State      VillagerState
	w          *world.World
}

func NewVillager(x, y int, w *world.World) *Villager {
	v := Villager{}
	id := fmt.Sprintf("villager_%d_%d", x, y)
	movement := NewMovement(x, y, w)
	storager := NewStorager(100)
	lumberjack := NewLumberjack(storager.IncrementInventory)

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

func (v *Villager) Tick(w *world.World, entities *[]Entity, buildings *building.BuildingsList) {
	if v.State == StateVillagerIdle {
		v.agent.AddCollectTreeGoal(w, v.movement.pos)
	}
	if v.storager.isOverweighted() {
		v.agent.SetState("overweighted")
		v.agent.AddStorageGoal(w, v.Pos(), v.storager.inventory)
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

func (v *Villager) ID() string {
	return v.id
}

func (v *Villager) Draw() {
	x, y := cnts.WorldToScreen(v.Pos().X, v.Pos().Y)
	src := rl.NewRectangle(41, 21, 16, 19)
	dst := rl.NewRectangle(x, y, cnts.TileSize, cnts.TileSize)

	rl.DrawTexturePro(spritebank.Human, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (v *Villager) Type() EntityType {
	return EntityTypeVillager
}
