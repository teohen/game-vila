package building

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/spritebank"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

type BuildingType string

const (
	StorageType BuildingType = "Storage"
)

type Storage struct {
	id        string
	Wood      int
	pos       cnts.Point
	typeBuild BuildingType
}

func NewStorage(x, y int) *Storage {
	return &Storage{
		id:        strings.ReplaceAll(uuid.NewString(), "-", ""),
		pos:       cnts.Point{X: x, Y: y},
		Wood:      50,
		typeBuild: StorageType,
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

func (s *Storage) ID() string {
	return s.id
}

func (s *Storage) Type() BuildingType {
	return StorageType
}
