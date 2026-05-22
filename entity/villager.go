package entity

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WaitDuration = 5
	MaxRetries   = 10
)

type VillagerState int

const (
	StateIdle    VillagerState = 0
	StateMoving  VillagerState = 1
	StateWaiting VillagerState = 2
	StateArrived VillagerState = 3
)

type VillagerType int

const (
	Human VillagerType = 1
)

type Villager struct {
	ID   string
	name string
	Type VillagerType
	X    int
	Y    int

	State     VillagerState
	TargetX   int
	TargetY   int
	Waypoints []pathfinding.Point
	WaitTicks int
	WaitCount int
}

func NewVillager(id, name string, x, y int) Villager {
	return Villager{
		ID:        id,
		name:      name,
		X:         x,
		Y:         y,
		Type:      Human,
		State:     StateIdle,
		WaitTicks: 0,
		WaitCount: 0,
	}
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
