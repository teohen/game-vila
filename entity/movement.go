package entity

import (
	"github/teohen/mgm-tto/agent"
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
	pos       cnts.Point
	TargetPos cnts.Point
	Waypoints []cnts.Point
	State     MovementState
	WaitTicks int
	WaitCount int
	w         *world.World
}

func NewMovement(x, y int, w *world.World) *Movement {
	return &Movement{
		pos:       cnts.Point{X: x, Y: y},
		TargetPos: cnts.Point{X: -1, Y: -1},
		Waypoints: make([]cnts.Point, 0),
		WaitTicks: 0,
		WaitCount: 0,
		w:         w,
		State:     StateMovementIdle,
	}
}

func (m *Movement) SetTarget(target cnts.Point) bool {
	path := pathfinding.FindPath(m.w, m.pos, target)

	if path == nil {
		return true
	}
	m.Waypoints = path
	m.State = StateMovementMoving
	m.TargetPos = target
	return false
}

func (m *Movement) Update() {
	switch m.State {
	case StateMovementIdle:
		m.State = StateMovementIdle
		return

	case StateMovementMoving:
		if len(m.Waypoints) == 0 {
			m.State = StateMovementArrived
			return
		}
		next := m.Waypoints[0]
		if next == m.TargetPos {
			if m.w.IsOccupied(m.TargetPos.X, m.TargetPos.Y) {
				m.Waypoints = m.Waypoints[1:]
				m.State = StateMovementArrived
				return
			}
			m.w.Vacate(m.pos.X, m.pos.Y)
			m.pos.X = next.X
			m.pos.Y = next.Y
			m.w.Occupy(m.pos.X, m.pos.Y)
			m.Waypoints = m.Waypoints[1:]
			m.State = StateMovementArrived
			return
		}
		if m.w.IsOccupied(next.X, next.Y) {
			m.State = StateMovementWaiting
			m.WaitTicks = 0
			m.WaitCount++
			return
		}
		m.w.Vacate(m.pos.X, m.pos.Y)
		m.pos.X = next.X
		m.pos.Y = next.Y
		m.w.Occupy(m.pos.X, m.pos.Y)
		m.Waypoints = m.Waypoints[1:]
		m.WaitCount = 0
		return

	case StateMovementWaiting:
		m.WaitTicks++
		if m.WaitTicks >= WaitDuration {
			if m.WaitCount >= MaxRetries {
				m.State = StateMovementIdle
				m.WaitCount = 0
				m.TargetPos = cnts.Point{X: 0, Y: 0}
				m.Waypoints = nil
				return
			}
			from := m.pos
			to := m.TargetPos
			path := pathfinding.FindPath(m.w, from, to)
			if len(path) == 0 {
				m.State = StateMovementIdle
				m.WaitCount = 0
				return
			}
			m.Waypoints = path
			m.State = StateMovementMoving
		}
		return
	case StateMovementArrived:
		m.State = StateMovementIdle
		m.WaitCount = 0
		m.TargetPos = cnts.Point{X: -1, Y: -1}
		m.Waypoints = nil
		return
	}
}

func (m *Movement) Pos() cnts.Point {
	return m.pos
}

func (m *Movement) ExecuteAction(target agent.Target) bool {
	if m.State == StateMovementIdle {
		m.SetTarget(target.Pos())
	} else {
		m.Update()
		if m.State == StateMovementArrived {
			return true
		}
	}
	return false
}

func (m *Movement) IsActor(actionType agent.ActionType) bool {
	return actionType == agent.ActionMoveType
}
