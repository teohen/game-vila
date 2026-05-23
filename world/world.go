package world

type World struct {
	cells    [][]Cell
	occupied [][]bool
}

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

func NewWorldFromCells(cells [][]CellType) *World {
	rows := len(cells)
	if rows == 0 {
		return &World{}
	}
	cols := len(cells[0])
	w := &World{
		cells:    make([][]Cell, rows),
		occupied: make([][]bool, rows),
	}
	for r := range cells {
		w.cells[r] = make([]Cell, cols)
		w.occupied[r] = make([]bool, cols)
		for c := range cells[r] {
			w.cells[r][c] = newTile(cells[r][c], r, c)
		}
	}
	return w
}

func (w *World) Draw() {
	for _, row := range w.cells {
		for _, cell := range row {
			cell.Draw()
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

const (
	terrainFrequency = 0.035
	waterThreshold   = -0.15
	dirtThreshold    = 0.05
)

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
