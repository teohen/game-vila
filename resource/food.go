package resource

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Food struct {
	id     string
	amount int
	pos    cnts.Point
	weight float32
}

func NewFood(pos cnts.Point, amount int) IResource {
	f := Food{
		id:     cnts.NewID(),
		amount: amount,
		pos:    pos,
	}

	return &f
}

func (f *Food) ID() string {
	return f.id
}
func (f *Food) Type() ResourceType {
	return ResourceWoodType
}
func (f *Food) Draw() {
	x, y := cnts.WorldToScreen(f.Pos().X, f.Pos().Y)
	src := rl.NewRectangle(96, 144, 16, 16)
	dst := rl.NewRectangle(x+5, y+5, cnts.TileSize-10, cnts.TileSize-10)
	rl.DrawTexturePro(spritebank.Structures, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (f *Food) Pos() cnts.Point {
	return f.pos
}

func (f *Food) Collectable() bool {
	return true
}

func (f *Food) Amount() int {
	return f.amount
}
