package entity

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"log"
	"sort"
)

type Plan struct {
	goal      *goap.State
	actions   []*Action
	TargetPos cnts.Point
}

type PlanType string

type IAgent interface {
	UpdateGoals(entities *[]Entity, w *world.World)
	OrderGoals()
	AddGoal(goal *Goal)
	RemoveGoal(goal *Goal)
	AddAction(action IAction)
	ChooseGoal() bool
	ExecuteAction() bool
	Movement() *Movement
	Lumberjack() *Lumberjack
}

type Agent struct {
	movement       Movement
	lumberjack     Lumberjack
	StartPlanState *goap.State
	Goals          []*Goal
	Actions        []goap.Action
	plan           Plan
	ActionIdx      int
	CurrentGoal    *Goal
}

func NewAgent(x, y int, w *world.World, ic IncrementWood) IAgent {
	a := Agent{
		Goals:          make([]*Goal, 0),
		Actions:        make([]goap.Action, 0),
		ActionIdx:      0,
		StartPlanState: goap.StateOf("!near_tree"),
		movement:       NewMovement(x, y, w),
		lumberjack:     NewLumberjack(ic),
	}

	return &a
}

func (a *Agent) AddGoal(g *Goal) {
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
	i := 1
	a.plan = Plan{}
	for len(a.Goals) > 0 && i <= len(a.Goals) {
		goal := a.Goals[len(a.Goals)-i]
		desired := goal.DesiredState()
		a.StartPlanState.Add("walkable")
		a.StartPlanState.Add(fmt.Sprintf("%s_health=%d", goal.Target().GetID(), 100))
		actions, err := goap.Plan(a.StartPlanState, desired, a.Actions)

		if err != nil {
			i += 1
			fmt.Println("not", err)
			goal.priority -= 5
			continue
		}

		for _, act := range actions {
			action := act.(*Action)
			action.target = goal.Target()
			a.plan.actions = append(a.plan.actions, action)
		}
		a.CurrentGoal = goal
		return true
	}
	return false
}

func (a *Agent) UpdateGoals(entities *[]Entity, w *world.World) {
	job := GetJobQueue().Peek()
	if job != nil {
		switch job.Type {
		case JobTypeChopTrees:
			ent := getEntityFrom(job.TargetID, entities)
			desired := fmt.Sprintf("%s_health=0", ent.GetID())
			a.AddGoal(NewGoal("ChopTree", desired, 50, ent))
			mv := NewAction("move_to", "walkable", "near_tree", &ent)
			mv.from = a.Movement().pos
			mv.world = w
			cp := NewAction("chop_tree", "near_tree", fmt.Sprintf("%s_health=0", ent.GetID()), &ent)
			a.AddAction(mv)
			a.AddAction(cp)
		}
	}
}

func (a *Agent) OrderGoals() {
	sort.Slice(a.Goals, func(i, j int) bool {
		gi := *a.Goals[i]
		gj := *a.Goals[j]
		return gi.Priority() > gj.Priority()
	})
}

func (a *Agent) ExecuteAction() bool {
	if a.ActionIdx >= len(a.plan.actions) {
		//TODO: Remove job from the job queue
		GetJobQueue().Remove()
		a.clearPlan()
		return true
	}
	action := a.plan.actions[a.ActionIdx]

	switch action.Name() {
	case "move_to":
		if a.Movement().State == StateMovementIdle {
			a.Movement().SetTarget(action.Target().Pos())
		} else {
			a.Movement().Update()
			if a.Movement().State == StateMovementArrived {
				a.nextAction()
			}
		}
	case "chop_tree":
		t, ok := action.target.(*Tree)
		if !ok {
			log.Fatal("TREE CONVERTION NOT WOTK")
		}
		if a.Lumberjack().State == StateLumberjackIdle {
			a.Lumberjack().Start(t)
		} else {
			if done := a.Lumberjack().Update(); done {
				a.Actions = make([]goap.Action, 0)
				a.RemoveGoal(a.CurrentGoal)
				a.nextAction()
			}
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

func (a *Agent) Movement() *Movement {
	return &a.movement
}

func (a *Agent) Lumberjack() *Lumberjack {
	return &a.lumberjack
}
