package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"testing"
)

type TargetTest struct {
	id string
}

func (tt *TargetTest) Pos() cnts.Point {
	return cnts.Point{X: 3, Y: 4}
}

func (tt *TargetTest) ID() string {
	return tt.id
}

type ActorTest struct {
}

func (at *ActorTest) ExecuteAction(target Target) bool {
	return true
}

func testPlan(t *testing.T, plan *Plan, goalID string, actionsLen, currAction int) bool {
	if plan.goal.ID() != goalID {
		t.Errorf("expect plan.goal.ID() to be %s. got=%s", goalID, plan.goal.ID())
		return false
	}

	if len(plan.actions) != actionsLen {
		t.Errorf("expect length of plan.actions to be %d. got=%d", actionsLen, len(plan.actions))
		return false
	}

	if plan.currAction != currAction {
		t.Errorf("expect plan.currAction to be %d. got=%d", currAction, plan.currAction)
		return false
	}

	return true
}

func TestShouldStorageIfOverweighted(t *testing.T) {
	w := world.NewWorld(10, 10)
	mv := ActorTest{}
	lj := ActorTest{}
	sto := ActorTest{}
	a := Agent{
		movement:   &mv,
		storager:   &sto,
		lumberjack: &lj,
	}

	collectGoal := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "ID_TREE"), &TargetTest{})
	storeGoal := NewGoalStoreInventory(fmt.Sprintf("%s_wood=%d", "STO_ID", 100), &TargetTest{id: "STO_ID"})
	a.State = goap.StateOf("walkable", "overweighted", "has_storage")

	a.AddGoal(storeGoal)
	a.AddGoal(collectGoal)
	planSet := a.ChooseGoal(&w, cnts.Point{X: 0, Y: 0})
	if !planSet {
		t.Fatalf("expect a.ChooseGoal return to be %t. got=%t", true, planSet)
	}

	if testPlan(t, a.plan, storeGoal.ID(), 2, 0) {
		return
	}
}

func TestShouldExecuteCollectTreePlan(t *testing.T) {
	w := world.NewWorld(10, 10)
	mv := ActorTest{}
	lj := ActorTest{}
	sto := ActorTest{}
	a := Agent{
		movement:   &mv,
		storager:   &sto,
		lumberjack: &lj,
	}

	collectGoal := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "ID_TREE"), &TargetTest{})
	a.State = goap.StateOf("walkable", "!overweighted")

	a.AddGoal(collectGoal)
	if planSet := a.ChooseGoal(&w, cnts.Point{X: 0, Y: 0}); !planSet {
		t.Fatalf("expect a.ChooseGoal return to be %t. got=%t", true, planSet)
	}

	if !testPlan(t, a.plan, collectGoal.ID(), 2, 0) {
		return
	}

	end := false
	for end == false {
		end = a.ExecutePlan()
		// TODO: test plan steps
	}
	if end {
		if a.plan.goal != nil || len(a.plan.actions) > 0 {
			t.Fatalf("a.plan expected to be=%s, got=%s", "zeroed", "properties")
			return
		}
	}
}
