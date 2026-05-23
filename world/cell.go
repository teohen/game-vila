package world

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CellType int

const (
	Empty CellType = 0
	Grass CellType = 1
	Dirt  CellType = 2
	Water CellType = 3

	SpriteSize = 16
)

func (ct CellType) Walkable() bool {
	return ct != Water
}

var terrainSources = map[CellType]rl.Rectangle{
	Empty: rl.NewRectangle(0, 0, 0, 0),
	Grass: rl.NewRectangle(432, 96, SpriteSize, SpriteSize),
	Dirt:  rl.NewRectangle(224, 16, SpriteSize, SpriteSize),
	Water: rl.NewRectangle(176, 304, SpriteSize, SpriteSize),
}

type Cell struct {
	Type        CellType
	row         int
	col         int
	size        float32
	BorderWidth float32
	BorderColor rl.Color
}

func newTile(t CellType, r, c int) Cell {
	tile := Cell{
		Type:        t,
		row:         r,
		col:         c,
		size:        constants.TileSize,
		BorderWidth: 1,
		BorderColor: rl.NewColor(255, 255, 255, 80),
	}

	return tile
}

func getSource(t *Cell) (rl.Rectangle, rl.Rectangle) {
	dst := rl.NewRectangle(
		float32(t.col)*t.size,
		float32(t.row)*t.size,
		t.size,
		t.size,
	)
	src := terrainSources[t.Type]
	return src, dst
}

func (t *Cell) Walkable() bool {
	return t.Type.Walkable()
}

func (t *Cell) Draw() {
	src, dst := getSource(t)

	rl.DrawRectangleRec(dst, t.bgColor())
	rl.DrawTexturePro(spritebank.Terrain, src, dst, rl.NewVector2(0, 0), 0, rl.White)
	rl.DrawRectangleLinesEx(dst, t.BorderWidth, t.BorderColor)
}

func (t *Cell) bgColor() rl.Color {
	switch t.Type {
	case Grass:
		return rl.NewColor(200, 230, 200, 255)
	case Dirt:
		return rl.NewColor(210, 190, 160, 255)
	case Water:
		return rl.NewColor(180, 200, 230, 255)
	default:
		return rl.NewColor(220, 220, 220, 255)
	}
}
