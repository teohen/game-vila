package entity

import (
	"github/teohen/mgm-tto/constants"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type VillagerType int

const (
	Human VillagerType = 1
)

type Villager struct {
	Movement
	ID   string
	name string
	Type VillagerType
}

func NewVillager(id, name string, x, y int) *Villager {
	return &Villager{
		Movement: Movement{
			X: x,
			Y: y,
		},
		ID:   id,
		name: name,
		Type: Human,
	}
}

func (v *Villager) Tick(w *world.World) MovementEvent {
	return v.Movement.Update(w)
}

func (v *Villager) Pos() (int, int) {
	return v.Movement.Pos()
}

func getSource(v *Villager) (rl.Rectangle, rl.Rectangle) {
	src := rl.NewRectangle(0, 0, 0, 0)
	dst := rl.NewRectangle(0, 0, 0, 0)

	switch v.Type {
	case Human:
		x, y := constants.WorldToScreen(v.X, v.Y)
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

func (v *Villager) Draw() {
	src, dst := getSource(v)
	rl.DrawTexturePro(spritebank.Human, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}
