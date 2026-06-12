package agent

import (
	"fmt"
	"github/teohen/mgm-tto/building"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/pathfinding"
	"github/teohen/mgm-tto/world"
	"log"
)

type IncrementWood func(amount int)

type Target interface {
	Pos() cnts.Point
	GetID() string
}

type IAgent interface {
	AddGoal(goal IGoal)
	RemoveGoal(ID string)
	ChooseGoal(w *world.World, pos cnts.Point) bool
	ExecutePlan() bool
	GetGoals() []IGoal
	AddStorageGoal(w *world.World, from cnts.Point, inventory int)
}

type Actor interface {
	ExecuteAction(target Target) bool
}

type Agent struct {
	movement       Actor
	lumberjack     Actor
	storager       Actor
	StartPlanState *goap.State
	Goals          []IGoal
	Actions        []goap.Action
	plan           *Plan
	ActionIdx      int
	CurrentGoal    *Goal
}

// TODO: REMOVE movement and lumberjack dependencies
func NewAgent(x, y int, w *world.World, movement, lumberjack, storager Actor) IAgent {
	a := Agent{
		Goals:          make([]IGoal, 0),
		Actions:        make([]goap.Action, 0),
		ActionIdx:      0,
		StartPlanState: goap.StateOf("walkable"),
		movement:       movement,
		lumberjack:     lumberjack,
		storager:       storager,
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
			fmt.Println("job removed", goal.GetID())
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
		mv := NewActionMove("walkable", "near", candidate.Target(), w, pos)
		cp := NewActionChopTree("near", candidate.Target())
		pi := NewActionPutInto("near", candidate.DesiredState(), candidate.Target())
		candidate.AddActions(mv, cp, pi)

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
		if g.GetType() == goalType {
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

		desired := fmt.Sprintf("%s_wood=%d", storage.GetID(), (storage.Wood + inventory))
		a.AddGoal(NewGoalStoreInventory(desired, storage))
		fmt.Println("StorageGoal added", near, storage.Wood, inventory)
	}
}
