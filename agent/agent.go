package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/job"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
	"strings"
)

type Target interface {
	Pos() cnts.Point
	ID() string
}

type IAgent interface {
	AddGoal(goal IGoal)
	RemoveGoal(ID string)
	ChooseGoal(w *world.World, pos cnts.Point) bool
	ExecutePlan() bool
	GetGoals() []IGoal
	AddStorageGoal(w *world.World, from cnts.Point, inventory int)
	AddCollectTreeGoal(w *world.World, from cnts.Point)
	SetState(state string)
	RegisterActor(typeAction ActionType, actor Actor)
}

type Actor interface {
	ExecuteAction(target Target) bool
}

type Agent struct {
	Owner   string
	Actors  map[ActionType]Actor
	State   *goap.State
	Goals   []IGoal
	Actions []goap.Action
	plan    *Plan
}

func NewAgent(x, y int, w *world.World, own string, initialState string) IAgent {
	a := Agent{
		Owner:   own,
		Goals:   make([]IGoal, 0),
		Actions: make([]goap.Action, 0),
		State:   goap.StateOf(strings.Split(initialState, ",")...),
		Actors:  make(map[ActionType]Actor),
	}

	return &a
}

func (a *Agent) AddGoal(g IGoal) {
	a.Goals = append(a.Goals, g)
}

func (a *Agent) RemoveGoal(id string) {
	for i, goal := range a.Goals {
		if id == goal.ID() {
			a.Goals = append(a.Goals[:i], a.Goals[i+1:]...)
		}
	}
}

func (ag *Agent) AddAction(a IAction) {
	ag.Actions = append(ag.Actions, a)
}

func (a *Agent) ChooseGoal(w *world.World, pos cnts.Point) bool {
	a.plan = NewPlan()
	cheapest := float32(10_000)

	for _, goal := range a.Goals {
		startPlan := a.State
		if !goal.IsRelevant(pos, a.State) {
			continue
		}

		for _, act := range goal.Actions() {
			act.Update(goal.Target(), pos)
		}

		actions, err := goap.Plan(startPlan, goal.DesiredState(), goal.GetGoapActions())
		if err != nil {
			if cnts.DEBUGGING {
				fmt.Println(err.Error())
			}
			continue
		}

		canditateCost := float32(0)
		candidateActions := make([]IAction, 0)
		for _, act := range actions {
			action := act.(IAction)
			action.SetTarget(action.Target())
			candidateActions = append(candidateActions, action)
			canditateCost += action.Cost()
		}

		if canditateCost < cheapest {
			a.plan.goal = goal
			a.plan.SetActions(candidateActions)
			cheapest = canditateCost
		}
	}

	if a.plan.IsSet() {
		return true
	}

	return false
}

func (a *Agent) ExecutePlan() bool {
	hasAct := a.plan.hasAction()
	if !hasAct {
		a.RemoveGoal(a.plan.goal.ID())
		a.Actions = make([]goap.Action, 0)
		a.plan.Clear()
		return true
	}

	action := a.plan.GetCurrentAction()
	typACt := action.Type()
	switch typACt {
	case ActionMoveType:
		actor := a.Actors[ActionMoveType]
		if done := actor.ExecuteAction(action.Target()); done {
			a.plan.nextAction()
		}
	case ActionChopTreeType:
		actor := a.Actors[ActionChopTreeType]
		if done := actor.ExecuteAction(action.Target()); done {
			a.plan.nextAction()
		}

	case ActionPutIntoType:
		actor := a.Actors[ActionPutIntoType]
		if done := actor.ExecuteAction(action.Target()); done {
			a.State.Apply(goap.StateOf("!overweighted"))
			a.plan.nextAction()
		}
	}

	return false
}

func (a *Agent) GetGoals() []IGoal {
	return a.Goals
}

func (a *Agent) GetGoalsOf(goalType GoalType) []IGoal {
	goals := make([]IGoal, 0)
	for _, g := range a.Goals {
		if g.Type() == goalType {
			goals = append(goals, g)
		}
	}
	return goals
}

func (a *Agent) SetState(state string) {
	a.State.Apply(goap.StateOf(state))
}

// TODO: remover essa função e adicionar na entity e remover essa logica do Agent
func (a *Agent) AddStorageGoal(w *world.World, from cnts.Point, inventory int) {
	if len(a.GetGoalsOf(GoalStoreInventoryType)) > 0 {
		return
	}

	g := NewGoalStoreInventory(w, nil, from)

	a.AddGoal(g)
}

// TODO: remover essa função e adicionar na entity e remover essa logica do Agent
func (a *Agent) AddCollectTreeGoal(w *world.World, from cnts.Point) {
	if len(job.GetJobQueue().Jobs) < 1 {
		return
	}

	targets := make([]Target, 0)
	for _, j := range job.GetJobQueue().Jobs {
		targets = append(targets, j.GetObject())
	}

	closest := pathfinding.FindClosest(w, from, targets)
	for _, j := range job.GetJobQueue().Jobs {
		if j.GetObject().Pos() == closest {
			g := NewGoalCollectTree(w, j.Object, from)
			a.AddGoal(g)
			job.GetJobQueue().Remove(j.Name(), j.Object.ID())
		}
	}
}

func (a *Agent) RegisterActor(typeAction ActionType, actor Actor) {
	a.Actors[typeAction] = actor
}
