package buildings

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Storage struct {
	ID   string
	wood int
	pos  cnts.Point
}

func NewStorage(x, y int) *Storage {
	id := fmt.Sprintf("storage_%d_%d", x, y)
	return &Storage{
		ID:  id,
		pos: cnts.Point{X: x, Y: y},
	}
}

func (s *Storage) Insert(amount int) {
	s.wood += amount
}

func (s *Storage) Draw() {
	x, y := cnts.WorldToScreen(s.pos.X, s.pos.Y)
	src := rl.NewRectangle(176, 112, 32, 16)
	dst := rl.NewRectangle(x, y, cnts.TileSize, cnts.TileSize)
	rl.DrawTexturePro(spritebank.Structures, src, dst, rl.NewVector2(0, 0), 0, rl.White)
}

func (s *Storage) Pos() cnts.Point {
	return s.pos
}
