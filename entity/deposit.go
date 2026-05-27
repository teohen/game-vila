package entity

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Deposit struct {
	X, Y int
}

func NewDeposit(x, y int) *Deposit {
	return &Deposit{X: x, Y: y}
}

func (d *Deposit) Tick(w *world.World) MovementEvent {
	return EventNone
}

func (d *Deposit) Pos() (int, int) {
	return d.X, d.Y
}

func (d *Deposit) Draw() {
	x, y := constants.WorldToScreen(d.X, d.Y)
	rl.DrawRectangle(int32(x), int32(y), int32(constants.TileSize), int32(constants.TileSize), rl.Red)
}
