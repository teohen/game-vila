package resource

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Wood struct {
	id     string
	amount int
	pos    cnts.Point
	weight float32
}

func NewWood(pos cnts.Point, amount int) IResource {
	w := Wood{
		id:     cnts.NewID(),
		amount: amount,
		pos:    pos,
	}

	return &w
}

func (w *Wood) ID() string {
	return w.id
}
func (w *Wood) Type() ResourceType {
	return ResourceWoodType
}
func (w *Wood) Draw() {
	x, y := cnts.WorldToScreen(w.Pos().X, w.Pos().Y)
	src := rl.NewRectangle(96, 144, 16, 16)
	dst := rl.NewRectangle(x+5, y+5, cnts.TileSize-10, cnts.TileSize-10)
	rl.DrawTexturePro(spritebank.Structures, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (w *Wood) Pos() cnts.Point {
	return w.pos
}

func (w *Wood) Collectable() bool {
	return true
}

func (w *Wood) Amount() int {
	return w.amount
}
