package building

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/inventory"
	"github/teohen/mgm-tto/spritebank"
	"github/teohen/mgm-tto/world"
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
	pos       cnts.Point
	typeBuild BuildingType
}

func NewStorage(x, y int) *Storage {
	return &Storage{
		id:        strings.ReplaceAll(uuid.NewString(), "-", ""),
		pos:       cnts.Point{X: x, Y: y},
		typeBuild: StorageType,
	}
}

func (s *Storage) Insert(amount int) {
	inventory.Get().AddToInventory(inventory.Wood, amount)
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

func (s *Storage) InteractionPos(w *world.World, from cnts.Point) cnts.Point {
	dirs := []cnts.Point{
		{X: s.pos.X, Y: s.pos.Y - 1},
		{X: s.pos.X, Y: s.pos.Y + 1},
		{X: s.pos.X - 1, Y: s.pos.Y},
		{X: s.pos.X + 1, Y: s.pos.Y},
	}
	for _, p := range dirs {
		if p.X >= 0 && p.X < w.Cols() && p.Y >= 0 && p.Y < w.Rows() {
			if w.IsWalkable(p.X, p.Y) && !w.IsOccupied(p.X, p.Y) {
				return p
			}
		}
	}
	return s.pos
}
