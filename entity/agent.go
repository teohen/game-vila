package entity

import (
	"fmt"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"log"
)

type PlanType string

const (
	PlanTypeChopTrees PlanType = "chop_tree"
)

type IAgent interface {
	SetPlan(planType PlanType, entities *[]Entity, target *Entity)
	ExecuteAction()
}

type Agent struct {
	Movement
	Lumberjack
	StartPlanState *goap.State
	Goals          []goap.State
	Actions        []goap.Action
	plan           Plan
	ActionIdx      int
}

func NewAgent(x, y int, w *world.World) Agent {
	a := Agent{
		Goals:          make([]goap.State, 0),
		Actions:        make([]goap.Action, 0),
		ActionIdx:      0,
		StartPlanState: goap.StateOf("!near_tree"),
		Movement:       NewMovement(x, y, w),
		Lumberjack:     NewLumberjack(),
	}

	a.Actions = append(a.Actions, NewAction("move_to", "!near_tree", "near_tree"))

	return a
}

func (a *Agent) SetPlan(planType PlanType, entities *[]Entity, target Entity) {
	switch planType {
	// TODO: fix this mess
	case PlanTypeChopTrees:
		trees := make([]Tree, 0)
		t, ok := target.(*Tree)

		if !ok {
			log.Fatal("target is not a tree")
		}
		for _, ent := range *entities {
			if ent.GetType() == EntityTypeTree {
				tree, ok := target.(*Tree)
				if !ok {
					log.Fatal("target is not a tree")
				}
				trees = append(trees, *tree)
			}
		}
		a.planChopTrees(t)
	}
}

func (a *Agent) planChopTrees(tree *Tree) {
	a.StartPlanState.Add(fmt.Sprintf("%s_health=%d", tree.ID, tree.Health))
	goapActions := append(a.Actions, NewAction("chop_tree", "near_tree", fmt.Sprintf("%s_health-20", tree.ID)))
	goal := goap.StateOf(fmt.Sprintf("%s_health=0", tree.ID))

	plan, err := goap.Plan(a.StartPlanState, goal, goapActions)

	if err != nil {
		log.Fatal("ERRRO", err.Error())
	}
	finalActions := make([]*Action, 0)

	for _, act := range plan {
		action := act.(*Action)
		action.target = tree
		finalActions = append(finalActions, action)

	}

	a.plan = Plan{
		goal:      goal,
		actions:   finalActions,
		TargetPos: tree.Pos(),
	}
}

func (a *Agent) ExecuteAction() bool {
	if a.ActionIdx > len(a.plan.actions) {
		return true
	}
	action := a.plan.actions[a.ActionIdx]

	switch action.name {
	case "move_to":
		if a.Movement.MovementState == StateMovementIdle {
			a.Movement.SetTarget(a.plan.TargetPos)
		} else {
			a.Movement.Update()
			if a.Movement.MovementState == StateMovementArrived {
				a.nextAction()
			}
		}
	case "chop_tree":
		t, ok := action.target.(*Tree)
		if !ok {
			log.Fatal("TREE CONVERTION NOT WOTK")
		}
		if !a.Lumberjack.IsHitting() {
			a.Lumberjack.Start(t)
		} else {
			_, done := a.Lumberjack.Update()
			a.nextAction()
			if done {
				a.clearPlan()
				return true
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
