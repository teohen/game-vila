package entity

import (
	"fmt"
	"sort"
	"strings"
)

type GoapState map[string]int

func (s GoapState) Clone() GoapState {
	c := make(GoapState, len(s))
	for k, v := range s {
		c[k] = v
	}
	return c
}

func (s GoapState) Satisfies(goal GoapState) bool {
	for k, v := range goal {
		if s[k] != v {
			return false
		}
	}
	return true
}

func (s GoapState) Key() string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s=%d,", k, s[k])
	}
	return b.String()
}

type GoapAction struct {
	Name      string
	Check     func(GoapState) bool
	Apply     func(GoapState) GoapState
	Cost      func(GoapState) int
	BuildStep func() PlanStep
}

type node struct {
	state GoapState
	steps []PlanStep
	cost  int
}

type GoapPlanner struct{}

func NewGoapPlanner() *GoapPlanner {
	return &GoapPlanner{}
}

func (p *GoapPlanner) Plan(current GoapState, goal GoapState, actions []GoapAction) []PlanStep {
	frontier := []node{{state: current.Clone(), steps: nil, cost: 0}}
	visited := make(map[string]bool)

	for len(frontier) > 0 {
		best := 0
		for i, n := range frontier {
			if n.cost < frontier[best].cost {
				best = i
			}
		}
		cur := frontier[best]
		frontier = append(frontier[:best], frontier[best+1:]...)

		if cur.state.Satisfies(goal) {
			return cur.steps
		}

		key := cur.state.Key()
		if visited[key] {
			continue
		}
		visited[key] = true

		for _, a := range actions {
			if !a.Check(cur.state) {
				continue
			}
			newState := a.Apply(cur.state.Clone())
			newSteps := append([]PlanStep(nil), cur.steps...)
			newSteps = append(newSteps, a.BuildStep())
			newCost := cur.cost + a.Cost(cur.state)
			frontier = append(frontier, node{newState, newSteps, newCost})
		}
	}

	return nil
}
