package pathfinding

import (
	"github/teohen/mgm-tto/cnts"
	"math"
)

type Positioner interface {
	Pos() cnts.Point
}

type WalkableGrid interface {
	IsWalkable(col, row int) bool
	IsOccupied(col, row int) bool
}

type cellInfo struct {
	g       int
	h       int
	f       int
	parent  cnts.Point
	visited bool
	inOpen  bool
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func heuristic(a, b cnts.Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func FindPath(grid WalkableGrid, from, to cnts.Point) []cnts.Point {

	if from == to {
		return nil
	}

	var info [cnts.GridRows][cnts.GridCols]cellInfo
	for r := 0; r < cnts.GridRows; r++ {
		for c := 0; c < cnts.GridCols; c++ {
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

	openList := []cnts.Point{from}

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

		dirs := []cnts.Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, d := range dirs {
			nx, ny := current.X+d.X, current.Y+d.Y
			if nx < 0 || nx >= cnts.GridRows || ny < 0 || ny >= cnts.GridCols {
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
				ni.h = heuristic(cnts.Point{nx, ny}, to)
				ni.f = g + ni.h
				ni.parent = current
				openList = append(openList, cnts.Point{nx, ny})
			} else if g < ni.g {
				ni.g = g
				ni.f = g + ni.h
				ni.parent = current
			}
		}
	}

	return nil
}

func reconstructPath(info [cnts.GridRows][cnts.GridCols]cellInfo, from, to cnts.Point) []cnts.Point {
	var path []cnts.Point
	current := to
	for current != from {
		path = append([]cnts.Point{current}, path...)
		current = info[current.Y][current.X].parent
	}
	return path
}

func FindClosest[T Positioner](w WalkableGrid, from cnts.Point, options []T) cnts.Point {
	closest := cnts.Point{X: -1, Y: -1}
	nearest := 10_000
	for _, o := range options {
		path := FindPath(w, from, o.Pos())
		if len(path) > 0 && len(path) < nearest {
			closest = o.Pos()
		}
	}
	return closest
}
