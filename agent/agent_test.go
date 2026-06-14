package agent

import (
	"fmt"
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/job"
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

func testPlan(t *testing.T, plan *Plan, goalType GoalType, goalDesired string, actionsLen, currAction int) bool {
	if !testGoal(t, plan.goal, goalType, goalDesired) {
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

func testGoal(t *testing.T, goal IGoal, typeG GoalType, desired string) bool {
	if goal.DesiredState().String() != desired {
		t.Errorf("expect goal.DesiredState() to be %s. got=%s", desired, goal.DesiredState().String())
		return false
	}

	if goal.Type() != typeG {
		t.Errorf("expect goal.Type() to be %s. got=%s", typeG, goal.Type())
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

	if testPlan(t, a.plan, GoalStoreInventoryType, "{STO_ID_wood=100}", 2, 0) {
		return
	}
}

func TestGetGoals(t *testing.T) {
	a := Agent{
		movement:   &ActorTest{},
		storager:   &ActorTest{},
		lumberjack: &ActorTest{},
	}
	g1 := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "T1"), &TargetTest{id: "T1"})
	g2 := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "T2"), &TargetTest{id: "T2"})
	g3 := NewGoalStoreInventory(fmt.Sprintf("%s_wood=%d", "STO", 100), &TargetTest{id: "STO"})

	a.AddGoal(g1)
	a.AddGoal(g2)
	a.AddGoal(g3)

	goals := a.GetGoals()
	if len(goals) != 3 {
		t.Fatalf("expected 3 goals, got %d", len(goals))
	}
	if goals[0].ID() != g1.ID() {
		t.Errorf("expected goals[0].ID()=%s, got=%s", g1.ID(), goals[0].ID())
	}
	if goals[1].ID() != g2.ID() {
		t.Errorf("expected goals[1].ID()=%s, got=%s", g2.ID(), goals[1].ID())
	}
	if goals[2].ID() != g3.ID() {
		t.Errorf("expected goals[2].ID()=%s, got=%s", g3.ID(), goals[2].ID())
	}
}

func TestRemoveGoal(t *testing.T) {
	a := Agent{
		movement:   &ActorTest{},
		storager:   &ActorTest{},
		lumberjack: &ActorTest{},
	}
	g1 := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "T1"), &TargetTest{id: "T1"})
	g2 := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "T2"), &TargetTest{id: "T2"})
	g3 := NewGoalStoreInventory(fmt.Sprintf("%s_wood=%d", "STO", 100), &TargetTest{id: "STO"})

	a.AddGoal(g1)
	a.AddGoal(g2)
	a.AddGoal(g3)

	a.RemoveGoal(g2.ID())

	goals := a.GetGoals()
	if len(goals) != 2 {
		t.Fatalf("expected 2 goals after removal, got %d", len(goals))
	}
	if goals[0].ID() != g1.ID() {
		t.Errorf("expected goals[0].ID()=%s, got=%s", g1.ID(), goals[0].ID())
	}
	if goals[1].ID() != g3.ID() {
		t.Errorf("expected goals[1].ID()=%s, got=%s", g3.ID(), goals[1].ID())
	}
}

func TestRemoveGoalNotInList(t *testing.T) {
	a := Agent{
		movement:   &ActorTest{},
		storager:   &ActorTest{},
		lumberjack: &ActorTest{},
	}
	g := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "T1"), &TargetTest{id: "T1"})
	a.AddGoal(g)

	a.RemoveGoal("non-existent-id")

	goals := a.GetGoals()
	if len(goals) != 1 {
		t.Fatalf("expected 1 goal, got %d", len(goals))
	}
	if goals[0].ID() != g.ID() {
		t.Errorf("expected goals[0].ID()=%s, got=%s", g.ID(), goals[0].ID())
	}
}

func TestChooseGoalReturnsFalseWhenNoGoalRelevant(t *testing.T) {
	w := world.NewWorld(10, 10)
	a := Agent{
		movement:   &ActorTest{},
		storager:   &ActorTest{},
		lumberjack: &ActorTest{},
		State:      goap.StateOf("walkable", "overweighted"),
	}
	g := NewGoalCollectTree(fmt.Sprintf("%s_health=0", "T1"), &TargetTest{})
	a.AddGoal(g)

	if a.ChooseGoal(&w, cnts.Point{X: 0, Y: 0}) {
		t.Fatal("expected ChooseGoal to return false when no goal is relevant")
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

	if !testPlan(t, a.plan, GoalCollectTreeType, "{ID_TREE_health=0}", 2, 0) {
		return
	}

	end := false
	for end == false {
		end = a.ExecutePlan()
	}
	if end {
		if a.plan.goal != nil || len(a.plan.actions) > 0 {
			t.Fatalf("a.plan expected to be=%s, got=%s", "zeroed", "properties")
			return
		}
	}
}

func TestShouldAddCollectTreeGoal(t *testing.T) {
	w := world.NewWorld(10, 10)
	mv := ActorTest{}
	lj := ActorTest{}
	sto := ActorTest{}
	a := Agent{
		movement:   &mv,
		storager:   &sto,
		lumberjack: &lj,
	}

	job.GetJobQueue().Push(*job.NewJob(
		job.JobChopTreeType, &TargetTest{id: "TREE_ID"},
	))

	a.AddCollectTreeGoal(&w, cnts.Point{X: 0, Y: 0})

	if len(a.Goals) != 1 {
		t.Fatalf("expect len(a.Goal) to be %d. got=%d", 1, len(a.Goals))
	}

	if !testGoal(t, a.Goals[0], GoalCollectTreeType, "{TREE_ID_health=0}") {
		return
	}

	if len(job.GetJobQueue().Jobs) != 0 {
		t.Fatalf("expect job.GetJobQueue().Jobs to be %d. got=%d", 0, len(job.GetJobQueue().Jobs))
	}
}
