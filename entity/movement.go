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
}

func (m *Movement) SetTarget(target cnts.Point, w *world.World) {
	path := pathfinding.FindPath(w, m.pos, target)

	if path == nil {
		return
	}
	m.Waypoints = path
	m.MovementState = StateMovementMoving
}

func (m *Movement) Update(w *world.World) MovementEvent {
	switch m.MovementState {
	case StateMovementIdle:
		return EventIdle

	case StateMovementMoving:
		if len(m.Waypoints) == 0 {
			m.MovementState = StateMovementArrived
			return EventArrived
		}
		next := m.Waypoints[0]
		if next == m.TargetPos {
			if w.IsOccupied(m.TargetPos.X, m.TargetPos.Y) {
				m.Waypoints = m.Waypoints[1:]
				m.MovementState = StateMovementArrived
				return EventArrived
			}
			w.Vacate(m.pos.X, m.pos.Y)
			m.pos.X = next.X
			m.pos.Y = next.Y
			w.Occupy(m.pos.X, m.pos.Y)
			m.Waypoints = m.Waypoints[1:]
			m.MovementState = StateMovementArrived
			return EventArrived
		}
		if w.IsOccupied(next.X, next.Y) {
			m.MovementState = StateMovementWaiting
			m.WaitTicks = 0
			m.WaitCount++
			return EventNone
		}
		w.Vacate(m.pos.X, m.pos.Y)
		m.pos.X = next.X
		m.pos.Y = next.Y
		w.Occupy(m.pos.X, m.pos.Y)
		m.Waypoints = m.Waypoints[1:]
		m.WaitCount = 0
		return EventNone

	case StateMovementWaiting:
		m.WaitTicks++
		if m.WaitTicks >= WaitDuration {
			if m.WaitCount >= MaxRetries {
				m.MovementState = StateMovementIdle
				m.WaitCount = 0
				m.TargetPos = cnts.Point{X: 0, Y: 0}
				m.Waypoints = nil
				return EventStuck
			}
			from := m.pos
			to := m.TargetPos
			path := pathfinding.FindPath(w, from, to)
			if len(path) == 0 {
				m.MovementState = StateMovementIdle
				m.WaitCount = 0
				return EventStuck
			}
			m.Waypoints = path
			m.MovementState = StateMovementMoving
		}
		return EventNone

	case StateMovementArrived:
		m.MovementState = StateMovementIdle
		m.WaitCount = 0
		m.TargetPos = cnts.Point{X: 0, Y: 0}
		m.Waypoints = nil
		return EventArrived
	}
	return EventNone
}

func (m *Movement) Pos() cnts.Point {
	return m.pos
}

func NewMovement(x, y int) Movement {
	return Movement{
		pos:       cnts.Point{X: x, Y: y},
		TargetPos: cnts.Point{X: -1, Y: -1},
		Waypoints: make([]cnts.Point, 0),
		WaitTicks: 0,
		WaitCount: 0,
	}
}
