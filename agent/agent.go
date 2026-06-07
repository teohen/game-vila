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
	UpdateGoals(w *world.World, pos cnts.Point, target Target, name string)
	AddGoal(goal IGoal)
	RemoveGoal(ID string)
	ChooseGoal(w *world.World, pos cnts.Point) bool
	ExecutePlan() bool
}

type Actor interface {
	ExecuteAction(target Target) bool
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

// TODO: REMOVE movement and lumberjack dependencies
func NewAgent(x, y int, w *world.World, ic IncrementWood, movement, lumberjack Actor) IAgent {
	a := Agent{
		Goals:          make([]IGoal, 0),
		Actions:        make([]goap.Action, 0),
		ActionIdx:      0,
		StartPlanState: goap.StateOf("walkable"),
		movement:       movement,
		lumberjack:     lumberjack,
	}

	return &a
}

func (a *Agent) AddGoal(g IGoal) {
	a.Goals = append(a.Goals, g)
}

func (a *Agent) RemoveGoal(id string) {
	for i, goal := range a.Goals {
		if id == goal.GetID() {
			a.Goals = append(a.Goals[:i], a.Goals[i+1:]...)
		}
	}
}

func (ag *Agent) AddAction(a IAction) {
	ag.Actions = append(ag.Actions, a)
}

func (a *Agent) ChooseGoal(w *world.World, pos cnts.Point) bool {
	a.plan = NewPlan()

	for _, candidate := range a.Goals {
		startPlan := a.StartPlanState
		mv := NewActionMove("walkable", "near_tree", candidate.Target(), w, pos)
		cp := NewActionChopTree("near_tree", candidate.Target())
		candidate.AddActions(mv, cp)

		actions, err := goap.Plan(startPlan, candidate.DesiredState(), candidate.GetGoapActions())
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
func (a *Agent) UpdateGoals(w *world.World, pos cnts.Point, target Target, name string) {
	g := NewGoal(name, fmt.Sprintf("%s_health=0", target.GetID()), target)
	a.AddGoal(g)
	job.GetJobQueue().Remove(name, target.GetID())

}

func (a *Agent) ExecutePlan() bool {
	hasAct := a.plan.hasAction()
	if !hasAct {

		a.RemoveGoal(a.plan.goal.GetID())
		a.Actions = make([]goap.Action, 0)
		a.plan.Clear()
		return true
	}

	action := a.plan.GetCurrentAction()
	typACt := action.Type()
	switch typACt {
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
