package entity

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tree struct {
	ID        string
	X, Y      int
	Health    int
	WoodYield int
}

func NewTree(id string, x, y, health, woodYield int) Tree {
	return Tree{
		ID:        id,
		X:         x,
		Y:         y,
		Health:    health,
		WoodYield: woodYield,
	}
}

func (t *Tree) Draw() {
	x, y := constants.WorldToScreen(t.X, t.Y)
	src := rl.NewRectangle(448, 192, 32, 32)
	dst := rl.NewRectangle(x, y, constants.TileSize, constants.TileSize)
	rl.DrawTexturePro(spritebank.Terrain, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}
