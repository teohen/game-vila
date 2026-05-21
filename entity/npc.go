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

type NPCState int

const (
	StateIdle    NPCState = 0
	StateMoving  NPCState = 1
	StateWaiting NPCState = 2
	StateArrived NPCState = 3
)

type NPCType int

const (
	Human NPCType = 1
)

type NPC struct {
	ID   string
	name string
	Type NPCType
	X    int
	Y    int

	State     NPCState
	TargetX   int
	TargetY   int
	Waypoints []pathfinding.Point
	WaitTicks int
	WaitCount int
}

func New(id, name string, x, y int) NPC {
	return NPC{
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

func getSource(n *NPC) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	switch n.Type {
	case Human:
		x, y := constants.WorldToScreen(n.X, n.Y)
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

func (n *NPC) Draw() {
	src, dst := getSource(n)
	rl.DrawTexturePro(spritebank.Human, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}
