package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/job"
	"github/teohen/mgm-tto/world"
)

type IncrementWood func(amount int)

type Target interface {
	Pos() cnts.Point
	GetID() string
}

type IAgent interface {
	UpdateGoals(w *world.World, pos cnts.Point, target Target, task Task)
	AddGoal(goal IGoal)
	RemoveGoal(goal *Goal)
	ChooseGoal() bool
	ExecutePlan() bool
}

type Actor interface {
	// TODO: the argument used to be a entity
	ExecuteAction(target Target) bool
}

type Task interface {
	Type() string
	Name() string
}

type Agent struct {
	movement       Actor
	lumberjack     Actor
	StartPlanState *goap.State
	Goals          []IGoal
	Actions        []goap.Action
	plan           *Plan
	ActionIdx      int
	CurrentGoal    *Goal
}

func NewAgent(x, y int, w *world.World, ic IncrementWood, movement, lumberjack Actor) IAgent {
	a := Agent{
		Goals:          make([]IGoal, 0),
		Actions:        make([]goap.Action, 0),
		ActionIdx:      0,
		StartPlanState: goap.StateOf("!near_tree"),
		movement:       movement,
		lumberjack:     lumberjack,
	}

	return &a
}

func (a *Agent) AddGoal(g IGoal) {
	a.Goals = append(a.Goals, g)
}

func (a *Agent) RemoveGoal(g *Goal) {
	for i, v := range a.Goals {
		if v == g {
			a.Goals = append(a.Goals[:i], a.Goals[i+1:]...)
		}
	}
}

func (ag *Agent) AddAction(a IAction) {
	ag.Actions = append(ag.Actions, a)
}

func (a *Agent) ChooseGoal() bool {
	a.plan = NewPlan()

	for _, candidate := range a.Goals {
		actions, err := goap.Plan(a.StartPlanState, candidate.DesiredState(), candidate.GetGoapActions())
		if err != nil {
			fmt.Println("not possible", err.Error())
			continue
		}
		a.plan.goal = candidate

		for _, act := range actions {
			action := act.(IAction)
			action.SetTarget(action.Target())
			a.plan.AppendActions(action)
		}
	}

	if a.plan.IsSet() {
		return true
	}

	return false
}

// TODO: UpdateAllJobsAsAPossibility
func (a *Agent) UpdateGoals(w *world.World, pos cnts.Point, target Target, task Task) {
	switch task.Type() {
	case "ChopTreeJob":
		mv := NewActionMove("walkable", "near_tree", target, w, pos)
		cp := NewActionChopTree("near_tree", target)
		g := NewGoal(task.Name(), fmt.Sprintf("%s_health=0", target.GetID()), target)
		g.AddActions(mv, cp)
		a.AddGoal(g)
	}
}

func (a *Agent) ExecutePlan() bool {
	if a.plan.hasAction() {
		if a.plan.job != nil {
			job.GetJobQueue().Remove(a.plan.job)
		}

		a.plan.Clear()
		return true
	}

	action := a.plan.GetCurrentAction()

	switch action.Type() {
	case ActionMoveType:
		if done := a.movement.ExecuteAction(action.Target()); done {
			a.plan.nextAction()
		}
	case ActionChopTreeType:
		if done := a.lumberjack.ExecuteAction(action.Target()); done {
			a.plan.nextAction()
		}
	}

	return false
}

func (a *Agent) nextAction() {
	a.ActionIdx += 1
}

func (a *Agent) clearPlan() {
	a.ActionIdx = 0
}
