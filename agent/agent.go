package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"strings"
)

type Target interface {
	Pos() cnts.Point
	ID() string
}

type InteractionPositioner interface {
	InteractionPos(w *world.World, from cnts.Point) cnts.Point
}

type IAgent interface {
	AddGoal(goal IGoal)
	RemoveGoal(ID string)
	ChooseGoal(w *world.World, pos cnts.Point) bool
	ExecutePlan() bool
	GetGoals() []IGoal
	SetState(state string)
	RegisterActor(typeAction ActionType, actor Actor)
	GetGoalsOf(goalType GoalType) []IGoal
}

type Actor interface {
	ExecuteAction(target Target) bool
}

type Agent struct {
	Owner  string
	Actors map[ActionType]Actor
	State  *goap.State
	Goals  []IGoal
	plan   *Plan
}

func NewAgent(x, y int, w *world.World, own string, initialState string) IAgent {
	a := Agent{
		Owner:  own,
		Goals:  make([]IGoal, 0),
		State:  goap.StateOf(strings.Split(initialState, ",")...),
		Actors: make(map[ActionType]Actor),
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

func (a *Agent) ChooseGoal(w *world.World, pos cnts.Point) bool {
	a.plan = NewPlan()
	cheapest := float32(10_000)

	for _, goal := range a.Goals {
		if !goal.IsRelevant(pos, a.State) {
			continue
		}

		for _, act := range goal.Actions() {
			act.Update(goal.Target(), pos)
		}

		actions, err := goap.Plan(a.State, goal.DesiredState(), goal.GetGoapActions())
		if err != nil {
			if cnts.DEBUGGING && goal.Type() == GoalRunAwayType {
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

	case ActionPickUpType:
		actor := a.Actors[ActionPickUpType]
		if done := actor.ExecuteAction(action.Target()); done {
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

func (a *Agent) RegisterActor(typeAction ActionType, actor Actor) {
	a.Actors[typeAction] = actor
}
