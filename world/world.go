package world

import (
	"github/teohen/mgm-tto/cnts"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type World struct {
	cells          [][]Cell
	occupied       [][]bool
	terrainTexture rl.RenderTexture2D
	terrainReady   bool
}

const (
	terrainFrequency = 0.035
	waterThreshold   = -0.15
	dirtThreshold    = 0.05
)

func NewWorld(rows, cols int) World {
	w := World{
		cells:    make([][]Cell, rows),
		occupied: make([][]bool, rows),
	}

	for i := range w.cells {
		w.cells[i] = make([]Cell, cols)
		w.occupied[i] = make([]bool, cols)
	}

	return w
}

func (w *World) renderTerrain() {
	tw := int32(cnts.GridCols * cnts.TileSize)
	th := int32(cnts.GridRows * cnts.TileSize)
	w.terrainTexture = rl.LoadRenderTexture(tw, th)
	rl.BeginTextureMode(w.terrainTexture)
	for _, row := range w.cells {
		for _, cell := range row {
			cell.Draw()
		}
	}
	rl.EndTextureMode()
	w.terrainReady = true
}

func (w *World) Draw() {
	if !w.terrainReady {
		w.renderTerrain()
	}
	src := rl.NewRectangle(0, 0, float32(w.terrainTexture.Texture.Width), float32(w.terrainTexture.Texture.Height))
	dst := rl.NewRectangle(0, 0, float32(cnts.GridCols*cnts.TileSize), float32(cnts.GridRows*cnts.TileSize))
	rl.DrawTexturePro(w.terrainTexture.Texture, src, dst, rl.NewVector2(0, 0), 0, rl.White)
	if cnts.DEBUGGING {
		for _, row := range w.cells {
			for _, cell := range row {
				cell.DrawDebug()
			}
		}
	}
}

func (w *World) GetCell(col, row int) *Cell {
	if col < 0 || col >= len(w.cells[0]) || row < 0 || row >= len(w.cells) {
		return nil
	}
	return &w.cells[row][col]
}

func (w *World) Rows() int {
	return len(w.cells)
}

func (w *World) Cols() int {
	if len(w.cells) == 0 {
		return 0
	}
	return len(w.cells[0])
}

func (w *World) Occupy(col, row int) bool {
	if col < 0 || col >= w.Cols() || row < 0 || row >= w.Rows() {
		return false
	}
	if w.occupied[row][col] {
		return false
	}
	w.occupied[row][col] = true
	return true
}

func (w *World) Vacate(col, row int) {
	if col < 0 || col >= w.Cols() || row < 0 || row >= w.Rows() {
		return
	}
	w.occupied[row][col] = false
}

func (w *World) IsOccupied(col, row int) bool {
	if col < 0 || col >= w.Cols() || row < 0 || row >= w.Rows() {
		return true
	}
	return w.occupied[row][col]
}

func (w *World) IsWalkable(col, row int) bool {
	cell := w.GetCell(col, row)
	if cell == nil {
		return false
	}
	return cell.Walkable()
}

func (w *World) Generate(seed int64) {
	n := NewNoise(seed)
	for r := range w.cells {
		for c := range w.cells[r] {
			v := n.Noise2D(float64(c)*terrainFrequency, float64(r)*terrainFrequency)
			switch {
			case v < waterThreshold:
				w.cells[r][c] = newTile(Water, r, c)
			case v < dirtThreshold:
				w.cells[r][c] = newTile(Dirt, r, c)
			default:
				w.cells[r][c] = newTile(Grass, r, c)
			}
		}
	}
}
