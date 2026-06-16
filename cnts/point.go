package cnts

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y int
}

func (p *Point) Near(dist int, other Point) bool {
	distance := math.Abs(float64(p.X-other.X)) + math.Abs(float64(p.Y-other.Y))
	return int(distance) <= dist
}

func (p *Point) Away(p2 Point) Point {
	dx := p2.X - p.X
	dy := p2.Y - p.Y

	if dx == 0 && dy == 0 {
		return Point{X: p2.X + 1, Y: p2.Y}
	}

	var stepX, stepY int

	if dx > 0 {
		stepX = 1
	} else if dx < 0 {
		stepX = -1
	}

	if dy > 0 {
		stepY = 1
	} else if dy < 0 {
		stepY = -1
	}

	if math.Abs(float64(dx)) >= math.Abs(float64(dy)) {
		return Point{X: p2.X + stepX, Y: p2.Y}
	} else {
		return Point{X: p2.X, Y: p2.Y + stepY}
	}
}

func (p *Point) Equals(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

func (p *Point) String() string {
	return fmt.Sprintf("%d_%d", p.X, p.Y)
}

type Pin struct {
	Id       string
	Position Point
}

func (p *Pin) Pos() Point {
	return p.Position
}

func (p *Pin) ID() string {
	return p.Id
}
