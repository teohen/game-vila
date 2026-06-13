package agent

import (
	"fmt"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/job"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
	"log"
)

type IncrementWood func(amount int)

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
	UpdateState(storage, overweighted bool)
}

type Actor interface {
	ExecuteAction(target Target) bool
}

type Agent struct {
	movement   Actor
	lumberjack Actor
	storager   Actor
	State      *goap.State
	Goals      []IGoal
	Actions    []goap.Action
	plan       *Plan
}

// TODO: REMOVE movement and lumberjack dependencies
func NewAgent(x, y int, w *world.World, movement, lumberjack, storager Actor) IAgent {
	a := Agent{
		Goals:      make([]IGoal, 0),
		Actions:    make([]goap.Action, 0),
		State:      goap.StateOf("walkable", "!has_storage", "!overweighted"),
		movement:   movement,
		lumberjack: lumberjack,
		storager:   storager,
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

func (a *Agent) UpdateState(storage, overweighted bool) {
	a.State = goap.StateOf("walkable", "!overweighted", "!has_storage")
	if storage {
		a.State.Apply(goap.StateOf("has_storage"))
	}
	if overweighted {
		a.State.Apply(goap.StateOf("overweighted"))
	}
}

func (a *Agent) ChooseGoal(w *world.World, pos cnts.Point) bool {
	a.plan = NewPlan()
	cheapest := float32(10_000)

	for _, goal := range a.Goals {
		startPlan := a.State
		if !goal.IsRelevant(a.State) {
			continue
		}

		mv := NewActionMove("walkable", "near", goal.Target(), w, pos)
		cp := NewActionChopTree("near", goal.Target())
		pi := NewActionPutInto("near", goal.DesiredState(), goal.Target())
		goal.SetActions(mv, cp, pi)
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
		if done := a.movement.ExecuteAction(action.Target()); done {
			a.plan.nextAction()
		}
	case ActionChopTreeType:
		if done := a.lumberjack.ExecuteAction(action.Target()); done {
			a.plan.nextAction()
		}

	case ActionPutIntoType:
		if done := a.storager.ExecuteAction(action.Target()); done {
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

func (a *Agent) AddStorageGoal(w *world.World, from cnts.Point, inventory int) {
	if len(a.GetGoalsOf(GoalStoreInventoryType)) > 0 {
		return
	}

	if near := pathfinding.FindClosest(w, from, building.Get().GetBuildingsOf(building.StorageType)); near.X != -1 {
		b := building.Get().GetBuildingAt(near)
		storage, ok := b.(*building.Storage)
		if !ok {
			log.Fatal("FOUND A BUILDING DIFERENT THAN A STORAGE")
		}

		desired := fmt.Sprintf("%s_wood=%d", storage.ID(), (storage.Wood + inventory))
		a.AddGoal(NewGoalStoreInventory(desired, storage))
	}
}

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
			a.AddGoal(NewGoalCollectTree(fmt.Sprintf("%s_health=0", j.Object.ID()), j.Object))
			job.GetJobQueue().Remove(j.Name(), j.Object.ID())
		}
	}
}
