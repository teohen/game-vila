package entity

import (
	"fmt"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type Tree struct {
	id        string
	pos       cnts.Point
	Health    int
	WoodYield int
}

func NewTree(x, y, health, woodYield int) *Tree {
	return &Tree{
		id:        strings.ReplaceAll(uuid.NewString(), "-", ""),
		pos:       cnts.Point{X: x, Y: y},
		Health:    health,
		WoodYield: woodYield,
	}
}

func (t *Tree) Tick(w *world.World, ent *[]Entity, buildings *building.BuildingsList) {
}

func (t *Tree) Pos() cnts.Point {
	return t.pos
}

func (t *Tree) Draw() {
	x, y := cnts.WorldToScreen(t.pos.X, t.pos.Y)
	src := rl.NewRectangle(448, 192, 32, 32)
	dst := rl.NewRectangle(x, y, cnts.TileSize-8, cnts.TileSize-8)
	rl.DrawTexturePro(spritebank.Terrain, src, dst, rl.NewVector2(0, 0), 0, rl.White)
	if cnts.DEBUGGING {
		rl.DrawText(fmt.Sprintf("%s", t.ID()), dst.ToInt32().X+8, dst.ToInt32().Y+8, 10, rl.Black)
	}
}

func (t *Tree) ID() string {
	return t.id
}

func (t *Tree) Type() EntityType {
	return EntityTypeTree
}
