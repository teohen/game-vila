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
	goal    *goap.State
	actions []*Action
	TargetX *int
	TargetY *int
}

type Villager struct {
	Movement
	Lumberjack
	Agent
	ID    string
	name  string
	Type  VillagerType
	State VillagerState
	World *world.World
}

func NewVillager(id, name string, x, y int) *Villager {
	v := &Villager{
		ID:         id,
		name:       name,
		Type:       Human,
		State:      "Idle",
		Movement:   NewMovement(x, y),
		Lumberjack: NewLumberjack(),
		Agent:      NewAgent(),
	}

	return v
}

func (v *Villager) Tick(entities *[]Entity, w *world.World) {
	v.World = w
	switch v.State {
	case StateVillagerIdle:
		v.SetPlan(entities)

	case StateVillagerBusy:
		finalAction := v.Agent.ExecuteAction()
		if finalAction {
			v.State = "Idle"
		}
	}

}

func (v *Villager) Name() string {
	return v.name
}

func (v *Villager) Pos() cnts.Point {
	return v.Movement.Pos()
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
