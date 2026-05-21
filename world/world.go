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

	w.Generate()

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

func (w *World) Generate() {
	for i := range w.cells {
		for j := range w.cells[i] {
			if i > 6 && i < 11 && j > 3 && j < 8 {
				w.cells[i][j] = newTile(Grass, i, j)
			} else if i == 20 && j == 21 {
				w.cells[i][j] = newTile(Dirt, i, j)
			} else if i == 30 && j == 21 {
				w.cells[i][j] = newTile(Water, i, j)
			} else {
				w.cells[i][j] = newTile(Empty, i, j)
			}
		}
	}
}
