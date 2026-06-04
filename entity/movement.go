package entity

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
)

const (
	WaitDuration = 5
	MaxRetries   = 10
)

type MovementState string

const (
	StateMovementIdle    MovementState = "idle"
	StateMovementMoving  MovementState = "moving"
	StateMovementWaiting MovementState = "waiting"
	StateMovementArrived MovementState = "arrived"
)

type Movement struct {
	pos           cnts.Point
	TargetPos     cnts.Point
	Waypoints     []cnts.Point
	MovementState MovementState
	WaitTicks     int
	WaitCount     int
	w             *world.World
}

func (m *Movement) SetTarget(target cnts.Point) {
	path := pathfinding.FindPath(m.w, m.pos, target)

	if path == nil {
		return
	}
	m.Waypoints = path
	m.MovementState = StateMovementMoving
}

func (m *Movement) Update() {
	switch m.MovementState {
	case StateMovementIdle:
		m.MovementState = StateMovementIdle

	case StateMovementMoving:
		if len(m.Waypoints) == 0 {
			m.MovementState = StateMovementArrived
		}
		next := m.Waypoints[0]
		if next == m.TargetPos {
			if m.w.IsOccupied(m.TargetPos.X, m.TargetPos.Y) {
				m.Waypoints = m.Waypoints[1:]
				m.MovementState = StateMovementArrived

			}
			m.w.Vacate(m.pos.X, m.pos.Y)
			m.pos.X = next.X
			m.pos.Y = next.Y
			m.w.Occupy(m.pos.X, m.pos.Y)
			m.Waypoints = m.Waypoints[1:]
			m.MovementState = StateMovementArrived

		}
		if m.w.IsOccupied(next.X, next.Y) {
			m.MovementState = StateMovementWaiting
			m.WaitTicks = 0
			m.WaitCount++

		}
		m.w.Vacate(m.pos.X, m.pos.Y)
		m.pos.X = next.X
		m.pos.Y = next.Y
		m.w.Occupy(m.pos.X, m.pos.Y)
		m.Waypoints = m.Waypoints[1:]
		m.WaitCount = 0

	case StateMovementWaiting:
		m.WaitTicks++
		if m.WaitTicks >= WaitDuration {
			if m.WaitCount >= MaxRetries {
				m.MovementState = StateMovementIdle
				m.WaitCount = 0
				m.TargetPos = cnts.Point{X: 0, Y: 0}
				m.Waypoints = nil
			}
			from := m.pos
			to := m.TargetPos
			path := pathfinding.FindPath(m.w, from, to)
			if len(path) == 0 {
				m.MovementState = StateMovementIdle
				m.WaitCount = 0

			}
			m.Waypoints = path
			m.MovementState = StateMovementMoving
		}

	case StateMovementArrived:
		m.MovementState = StateMovementIdle
		m.WaitCount = 0
		m.TargetPos = cnts.Point{X: -1, Y: -1}
		m.Waypoints = nil
	}
}

func (m *Movement) Pos() cnts.Point {
	return m.pos
}

func NewMovement(x, y int, w *world.World) Movement {
	return Movement{
		pos:       cnts.Point{X: x, Y: y},
		TargetPos: cnts.Point{X: -1, Y: -1},
		Waypoints: make([]cnts.Point, 0),
		WaitTicks: 0,
		WaitCount: 0,
		w:         w,
	}
}
