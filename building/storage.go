package building

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BuildingType string

const (
	StorageType BuildingType = "Storage"
)

type Storage struct {
	ID   string
	Wood int
	pos  cnts.Point
	Type BuildingType
}

func NewStorage(x, y int) *Storage {
	return &Storage{
		ID:   "dalskdjalskdls",
		pos:  cnts.Point{X: x, Y: y},
		Wood: 50,
		Type: StorageType,
	}
}

func (s *Storage) Insert(amount int) {
	s.Wood += amount
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

func (s *Storage) GetID() string {
	return s.ID
}

func (s *Storage) GetType() BuildingType {
	return StorageType
}
