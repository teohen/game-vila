package entity

import (
	"fmt"
	"github/teohen/mgm-tto/constants"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type NPCType int

const (
	Human NPCType = 1
)

type NPC struct {
	ID      string
	name    string
	Type    NPCType
	x       int
	y       int
	texture rl.Texture2D
}

func New(id, name string, x, y int) NPC {
	n := NPC{
		ID:      id,
		name:    name,
		x:       x,
		y:       y,
		Type:    Human,
		texture: rl.LoadTexture("./res/assets/player_anims.png"),
	}

	return n
}

func getSource(n *NPC) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	switch n.Type {
	case Human:
		x, y := constants.GridToScreen(n.x, n.y)
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
	rl.DrawTexturePro(n.texture, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func pr() {
	fmt.Println("")
}

