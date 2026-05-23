package pathfinding

import (
	"github/teohen/mgm-tto/constants"
	"math"
)

type Point struct {
	X, Y int
}

type WalkableGrid interface {
	IsWalkable(col, row int) bool
	IsOccupied(col, row int) bool
}

type cellInfo struct {
	g       int
	h       int
	f       int
	parent  Point
	visited bool
	inOpen  bool
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func heuristic(a, b Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func FindPath(grid WalkableGrid, from, to Point) []Point {
	if from == to {
		return nil
	}

	var info [constants.GridRows][constants.GridCols]cellInfo
	for r := 0; r < constants.GridRows; r++ {
		for c := 0; c < constants.GridCols; c++ {
			info[r][c] = cellInfo{g: math.MaxInt32}
		}
	}

	startH := heuristic(from, to)
	info[from.Y][from.X] = cellInfo{
		g:      0,
		h:      startH,
		f:      startH,
		inOpen: true,
	}

	openList := []Point{from}

	for len(openList) > 0 {
		bestIdx := 0
		for i, p := range openList {
			if info[p.Y][p.X].f < info[openList[bestIdx].Y][openList[bestIdx].X].f {
				bestIdx = i
			}
		}

		current := openList[bestIdx]
		if current == to {
			return reconstructPath(info, from, to)
		}

		openList = append(openList[:bestIdx], openList[bestIdx+1:]...)
		info[current.Y][current.X].inOpen = false
		info[current.Y][current.X].visited = true

		dirs := []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, d := range dirs {
			nx, ny := current.X+d.X, current.Y+d.Y
			if nx < 0 || nx >= constants.GridRows || ny < 0 || ny >= constants.GridCols {
				continue
			}

			if info[ny][nx].visited {
				continue
			}

			isTarget := nx == to.X && ny == to.Y
			if !isTarget {
				if !grid.IsWalkable(nx, ny) || grid.IsOccupied(nx, ny) {
					continue
				}
			}

			g := info[current.Y][current.X].g + 1
			ni := &info[ny][nx]
			if !ni.inOpen {
				ni.inOpen = true
				ni.g = g
				ni.h = heuristic(Point{nx, ny}, to)
				ni.f = g + ni.h
				ni.parent = current
				openList = append(openList, Point{nx, ny})
			} else if g < ni.g {
				ni.g = g
				ni.f = g + ni.h
				ni.parent = current
			}
		}
	}

	return nil
}

func reconstructPath(info [constants.GridRows][constants.GridCols]cellInfo, from, to Point) []Point {
	var path []Point
	current := to
	for current != from {
		path = append([]Point{current}, path...)
		current = info[current.Y][current.X].parent
	}
	return path
}
