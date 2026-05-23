package entity

import (
	"fmt"

	"github/teohen/mgm-tto/debug"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
)

const (
	WaitDuration = 5
	MaxRetries   = 10
)

type MovementState int

const (
	StateIdle    MovementState = 0
	StateMoving  MovementState = 1
	StateWaiting MovementState = 2
	StateArrived MovementState = 3
)

type Movement struct {
	X, Y      int
	TargetX   int
	TargetY   int
	Waypoints []pathfinding.Point
	State     MovementState
	WaitTicks int
	WaitCount int
}

func (m *Movement) SetTarget(x, y int, w *world.World) {
	m.TargetX = x
	m.TargetY = y
	from := pathfinding.Point{X: m.X, Y: m.Y}
	to := pathfinding.Point{X: x, Y: y}
	path := pathfinding.FindPath(w, from, to)
	m.debugMovement(path, x, y)

	if path == nil {
		return
	}
	m.Waypoints = path
	m.State = StateMoving
}

func (m *Movement) Update(w *world.World) MovementEvent {
	defer m.debugMovement(nil, 0, 0)
	switch m.State {
	case StateIdle:
		return EventIdle

	case StateMoving:
		if len(m.Waypoints) == 0 {
			m.State = StateArrived
			return EventArrived
		}
		next := m.Waypoints[0]
		if next.X == m.TargetX && next.Y == m.TargetY {
			if w.IsOccupied(m.TargetX, m.TargetY) {
				m.Waypoints = m.Waypoints[1:]
				m.State = StateArrived
				return EventArrived
			}
			w.Vacate(m.X, m.Y)
			m.X = next.X
			m.Y = next.Y
			w.Occupy(m.X, m.Y)
			m.Waypoints = m.Waypoints[1:]
			m.State = StateArrived
			return EventArrived
		}
		if w.IsOccupied(next.X, next.Y) {
			m.State = StateWaiting
			m.WaitTicks = 0
			m.WaitCount++
			return EventNone
		}
		w.Vacate(m.X, m.Y)
		m.X = next.X
		m.Y = next.Y
		w.Occupy(m.X, m.Y)
		m.Waypoints = m.Waypoints[1:]
		m.WaitCount = 0
		return EventNone

	case StateWaiting:
		m.WaitTicks++
		if m.WaitTicks >= WaitDuration {
			if m.WaitCount >= MaxRetries {
				m.State = StateIdle
				m.WaitCount = 0
				m.TargetX = 0
				m.TargetY = 0
				m.Waypoints = nil
				return EventStuck
			}
			from := pathfinding.Point{X: m.X, Y: m.Y}
			to := pathfinding.Point{X: m.TargetX, Y: m.TargetY}
			path := pathfinding.FindPath(w, from, to)
			if len(path) == 0 {
				m.State = StateIdle
				m.WaitCount = 0
				return EventStuck
			}
			m.Waypoints = path
			m.State = StateMoving
		}
		return EventNone

	case StateArrived:
		m.State = StateIdle
		m.WaitCount = 0
		m.TargetX = 0
		m.TargetY = 0
		m.Waypoints = nil
		return EventArrived
	}
	return EventNone
}

func (m *Movement) Pos() (int, int) {
	return m.X, m.Y
}

func (m *Movement) debugMovement(path []pathfinding.Point, x, y int) {
	if debug.IsEnabled(debug.Move) {
		if path != nil {
			fmt.Printf("[MOVEMENT - TARGER] Move path found len=%d from (%d,%d) to (%d,%d)\n", len(path), m.X, m.Y, x, y)
		} else if x != 0 || y != 0 {
			fmt.Printf("[MOVEMENT - TARGET] Move no path from (%d,%d) to (%d,%d)\n", m.X, m.Y, x, y)
		}

		fmt.Printf("[DEBUG] Move state=%d pos=(%d,%d) target=(%d,%d) waypoints=%d\n",
			m.State, m.X, m.Y, m.TargetX, m.TargetY, len(m.Waypoints))
	}
}
