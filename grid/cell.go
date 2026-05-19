package grid

import (
	"fmt"
	"github/teohen/mgm-tto/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileType int

const (
	Empty TileType = 0
	Grass TileType = 1
	Dirt  TileType = 2
	Water TileType = 3

	SpriteSize = 16
)

type Cell struct {
	Type        TileType
	texture     rl.Texture2D
	row         int
	col         int
	size        float32
	Selected    bool
	BorderWidth float32
	BorderColor rl.Color
}

func newTile(t TileType, r, c int, texture rl.Texture2D) Cell {
	tile := Cell{
		Type:        t,
		texture:     texture,
		row:         r,
		col:         c,
		size:        constants.TileSize,
		BorderWidth: 1,
		BorderColor: rl.NewColor(50, 50, 50, 255),
	}

	return tile
}

func getSource(t *Cell) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)
	switch t.Type {
	case Grass:
		dst.X = (float32(t.col) * t.size)
		dst.Y = (float32(t.row) * t.size)
		dst.Width = t.size
		dst.Height = t.size

		src.X = 16
		src.Y = 16
		src.Width = SpriteSize
		src.Height = SpriteSize
	default:
		dst.X = (float32(t.col) * t.size)
		dst.Y = (float32(t.row) * t.size)
		dst.Width = t.size
		dst.Height = t.size
	}

	return src, dst
}

func (t *Cell) draw() {
	src, dst := getSource(t)

	rl.DrawRectangleRec(dst, t.bgColor())

	// Draw border
	borderClr := t.BorderColor
	if t.Selected {
		borderClr = rl.Blue
	}
	rl.DrawRectangleLinesEx(dst, t.BorderWidth, borderClr)

	rl.DrawTexturePro(t.texture, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (t *Cell) bgColor() rl.Color {
	switch t.Type {
	case Grass:
		return rl.NewColor(200, 230, 200, 255) // light green
	case Dirt:
		return rl.NewColor(210, 190, 160, 255) // light brown
	case Water:
		return rl.NewColor(180, 200, 230, 255) // light blue
	default: // Empty
		return rl.NewColor(220, 220, 220, 255) // light gray
	}
}

func pr() {
	fmt.Println("")
}
