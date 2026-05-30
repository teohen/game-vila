package entity

import (
	"github/teohen/mgm-tto/goap"
	"strings"
)

func NewAction(name, require, outcome string) *Action {
	return &Action{
		name:    name,
		require: goap.StateOf(strings.Split(require, ",")...),
		outcome: goap.StateOf(strings.Split(outcome, ",")...),
	}
}

type Action struct {
	name    string
	cost    int
	require *goap.State
	outcome *goap.State
}

func (a *Action) Simulate(current *goap.State) (*goap.State, *goap.State) {
	return a.require, a.outcome
}

func (a *Action) Cost() float32 {
	return 1
}

func (a *Action) String() string {
	return a.name
}
