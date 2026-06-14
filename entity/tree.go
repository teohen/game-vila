package entity

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type Tree struct {
	id        string
	pos       cnts.Point
	Health    int
	WoodYield int
	marked    bool
}

func NewTree(x, y, health, woodYield int) *Tree {
	return &Tree{
		id:        strings.ReplaceAll(uuid.NewString(), "-", ""),
		pos:       cnts.Point{X: x, Y: y},
		Health:    health,
		WoodYield: woodYield,
	}
}

func (t *Tree) Tick() {
}

func (t *Tree) Pos() cnts.Point {
	return t.pos
}

func (t *Tree) Draw() {
	x, y := cnts.WorldToScreen(t.pos.X, t.pos.Y)
	src := rl.NewRectangle(448, 192, 32, 32)
	dst := rl.NewRectangle(x, y, cnts.TileSize, cnts.TileSize)
	rl.DrawTexturePro(spritebank.Terrain, src, dst, rl.NewVector2(0, 0), 0, rl.White)
	if t.marked {
		rl.DrawText("X", dst.ToInt32().X+16, dst.ToInt32().Y+8, 20, rl.Red)
	}
	if cnts.DEBUGGING {
		rl.DrawText(fmt.Sprintf("%s", t.ID()[27:]), dst.ToInt32().X+8, dst.ToInt32().Y+8, 10, rl.Black)
	}
}

func (t *Tree) ID() string {
	return t.id
}

func (t *Tree) Type() EntityType {
	return EntityTypeTree
}

func (t *Tree) Mark(v bool) {
	t.marked = v
}
