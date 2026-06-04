package entity

import (
	"fmt"
	"github/teohen/mgm-tto/events"
	"github/teohen/mgm-tto/goap"
	"log"
)

type IAgent interface {
	SetPlan(entities *[]Entity)
	ExecuteAction()
}

type Agent struct {
	StartPlanState *goap.State
	Goals          []goap.State
	Actions        []goap.Action
	plan           Plan
	ActionIdx      int
}

func NewAgent() Agent {
	a := Agent{
		Goals:          make([]goap.State, 0),
		Actions:        make([]goap.Action, 0),
		ActionIdx:      0,
		StartPlanState: goap.StateOf("!near_tree"),
	}

	a.Actions = append(a.Actions, NewAction("move_to", "!near_tree", "near_tree"))

	return a
}

func (a *Agent) SetPlan(entities *[]Entity) {
	if len(GetJobQueue().Jobs) > 0 {
		job := GetJobQueue().Pop()
		switch job.Type {
		case JobTypeChopTrees:
			t := getTreeFrom(job.TargetID, entities)
			v.VilState.Add(fmt.Sprintf("%s_health=%d", t.ID, t.Health))
			goapActions := append(v.Actions, NewAction("chop_tree", "near_tree", fmt.Sprintf("%s_health-20", t.ID)))
			goal := goap.StateOf(fmt.Sprintf("%s_health=0", t.ID))
			plan, err := goap.Plan(v.VilState, goal, goapActions)

			if err != nil {
				log.Fatal("ERRRO", err.Error())
			}
			finalActions := make([]*Action, 0)

			for _, act := range plan {
				action := act.(*Action)
				action.target = t
				finalActions = append(finalActions, action)

			}

			v.plan = Plan{
				goal:    goal,
				actions: finalActions,
				TargetX: &job.TargetX,
				TargetY: &job.TargetY,
			}
			v.State = "Planning"
		}
	}
}

func (a *Agent) ExecutePlan() bool {
	if v.ActionIdx >= len(v.plan.actions) {
		return true
	}
	v.executeAction()
	return false
}

func (a *Agent) ExecuteAction() bool {
	action := v.plan.actions[v.ActionIdx]
	switch action.name {
	case "move_to":
		if v.Movement.TargetX < 1 && v.Movement.TargetY < 1 {
			v.Movement.SetTarget(*v.plan.TargetX, *v.plan.TargetY, v.World)
		} else {
			movState := v.Movement.Update(v.World)
			if movState == EventArrived {
				v.ActionIdx += 1
			}
		}
	case "chop_tree":
		t, ok := action.target.(*Tree)
		if !ok {
			log.Fatal("TREE CONVERTION NOT WOTK")
		}
		if !v.Lumberjack.IsHitting() {
			v.Lumberjack.Start(t)
		} else {
			_, done := v.Lumberjack.Update(v.World)
			if done {
				events.Emit(events.GameEvent{
					Type: events.EventTreeCut,
					Payload: map[string]interface{}{
						"treeX": t.X,
						"treeY": t.Y,
					},
				})
				v.ActionIdx += 1
			}
		}
	}
}

func getTreeFrom(id string, entities *[]Entity) *Tree {
	for _, e := range *entities {
		tree, ok := e.(*Tree)
		if !ok {
			continue
		}

		if tree.ID == id {
			return tree
		}
	}
	return nil
}
