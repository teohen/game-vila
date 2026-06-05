package entity

import (
	"github/teohen/mgm-tto/goap"
	"strings"
)

type Goal struct {
	name         string
	desiredState *goap.State
	priority     int

	target Entity
}

type IGoal interface {
	Priority() int
	EvaluatePriority() int
	DesiredState() *goap.State
	Target() Entity
}

func NewGoal(name, desired string, p int, t Entity) *Goal {
	g := Goal{
		name:         name,
		desiredState: goap.StateOf(strings.Split(desired, ",")...),
		priority:     p,
		target:       t,
	}

	return &g
}

func (g *Goal) Priority() int {
	return g.priority
}

func (g *Goal) DesiredState() *goap.State {
	return g.desiredState
}

func (g *Goal) EvaluatePriority() int {
	return 1
}

func (g *Goal) Target() Entity {
	return g.target
}
