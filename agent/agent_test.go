package agent

import (
	"github/teohen/mgm-tto/cnts"
	"github/teohen/mgm-tto/goap"
	"github/teohen/mgm-tto/world"
	"testing"
)

type TargetTest struct {
}

func (tt *TargetTest) Pos() cnts.Point {
	return cnts.Point{X: 3, Y: 4}
}

func (tt *TargetTest) ID() string {
	return ""
}

type ActorTest struct {
}

func (at *ActorTest) ExecuteAction(target Target) bool {
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

	highP := NewGoalCollectTree("", &TargetTest{})
	g2 := NewGoalStoreInventory("g2", &TargetTest{})
	a.State = goap.StateOf("walkable", "!overweighted", "!has_storage")

	a.ChooseGoal(&w, cnts.Point{X: 0, Y: 0})

	if a.plan.goal.ID() != highP.ID() {
		t.Fatalf("expect plan.goal.ID() to be %s. got=%s", highP.ID(), a.plan.goal.ID())
	}
}
