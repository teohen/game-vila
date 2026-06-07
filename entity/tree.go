package entity

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tree struct {
	ID        string
	pos       cnts.Point
	Health    int
	WoodYield int
}

func NewTree(id string, x, y, health, woodYield int) *Tree {
	return &Tree{
		ID:        id,
		pos:       cnts.Point{X: x, Y: y},
		Health:    health,
		WoodYield: woodYield,
	}
}

func (t *Tree) Tick(ent *[]Entity, w *world.World) {
}

func (t *Tree) Pos() cnts.Point {
	return t.pos
}

func (t *Tree) Draw() {
	x, y := cnts.WorldToScreen(t.pos.X, t.pos.Y)
	src := rl.NewRectangle(448, 192, 32, 32)
	dst := rl.NewRectangle(x, y, cnts.TileSize-8, cnts.TileSize-8)
	rl.DrawTexturePro(spritebank.Terrain, src, dst, rl.NewVector2(0, 0), 0, rl.White)
	rl.DrawText(fmt.Sprintf("%s", t.ID), dst.ToInt32().X+8, dst.ToInt32().Y+8, 10, rl.Black)
}

func (t *Tree) GetID() string {
	return t.ID
}

func (t *Tree) GetType() EntityType {
	return EntityTypeTree
}
