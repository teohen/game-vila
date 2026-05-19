package grid

import (
	"fmt"
	"github/teohen/mgm-tto/constants"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Grid struct {
	cells       [][]Cell
	texture     rl.Texture2D
	tileTexture rl.Texture2D

	IsDragging bool
	DragStart  rl.Vector2
	DragEnd    rl.Vector2
	Selected   map[rl.Vector2]bool
}

func NewGrid(rows, cols int) Grid {
	g := Grid{
		cells:   make([][]Cell, rows),
		texture: rl.LoadTexture("./res/assets/spr_terrains.png"),
	}

	for i, _ := range g.cells {
		g.cells[i] = make([]Cell, cols)
	}

	g.Generate()

	return g
}

func GridToScreen(col, row int) (x, y float32) {
	return float32(col)*constants.TileSize + constants.GridOffsetX, float32(row)*constants.TileSize + constants.GridOffsetY
}

func ScreenToGrid(mx, my int32) (int, int, bool) {
	gx := int(mx-constants.GridOffsetX) / constants.TileSize
	gy := int(my-constants.GridOffsetY) / constants.TileSize
	if gx < 0 || gx >= constants.GridCols || gy < 0 || gy >= constants.GridRows {
		return 0, 0, false
	}
	return gx, gy, true
}

func (g *Grid) Draw() {
	//fmt.Println("draw grid", g.cells[3][3].col)
	for _, row := range g.cells {
		for _, tile := range row {
			tile.draw()
		}
	}

	// Draw selection box while dragging
	if g.IsDragging {
		w := g.DragEnd.X - g.DragStart.X
		h := g.DragEnd.Y - g.DragStart.Y
		rl.DrawRectangleLines(
			int32(g.DragStart.X),
			int32(g.DragStart.Y),
			int32(w),
			int32(h),
			rl.Blue,
		)
	}

	// Draw highlight border on selected tiles
	for pos := range g.Selected {
		col, row := int(pos.X), int(pos.Y)
		x, y := GridToScreen(col, row)
		rl.DrawRectangleLines(
			int32(x), int32(y),
			int32(constants.TileSize), int32(constants.TileSize),
			rl.Red,
		)
	}
}

func (g *Grid) destroy() {
	rl.UnloadTexture(g.texture)
}

func (g *Grid) SelectTilesInBox() {
	minX := math.Min(float64(g.DragStart.X), float64(g.DragEnd.X))
	maxX := math.Max(float64(g.DragStart.X), float64(g.DragEnd.X))
	minY := math.Min(float64(g.DragStart.Y), float64(g.DragEnd.Y))
	maxY := math.Max(float64(g.DragStart.Y), float64(g.DragEnd.Y))

	g.Selected = make(map[rl.Vector2]bool)

	for row := 0; row < len(g.cells); row++ {
		for col := 0; col < len(g.cells[row]); col++ {
			x, y := GridToScreen(col, row)
			tileCenterX := float64(x + constants.TileSize/2)
			tileCenterY := float64(y + constants.TileSize/2)

			if tileCenterX >= minX && tileCenterX <= maxX &&
				tileCenterY >= minY && tileCenterY <= maxY {
				g.Selected[rl.NewVector2(float32(col), float32(row))] = true
			}
		}
	}
}

func pr2() {
	fmt.Println("")
}

func (g *Grid) Generate() {
	for i, _ := range g.cells {
		for j, _ := range g.cells[i] {
			if i > 6 && i < 11 && j > 3 && j < 8 {
				g.cells[i][j] = newTile(Grass, i, j, g.texture)
			} else {
				g.cells[i][j] = newTile(Empty, i, j, g.texture)
			}
		}
	}
}
