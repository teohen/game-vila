package entity

import (
	"github/teohen/mgm-tto/resource"
	"testing"
)

var wood = 0

func reset() {
	wood = 0
}

func TestCutDownTree(t *testing.T) {

	tree := resource.NewTree(0, 0, 100, 100)
	lj := NewLumberjack()
	lj.hit = 50

	end := false
	actionCounter := 0

	for !end {
		if final := lj.ExecuteAction(tree); final {
			if lj.State != StateLumberjackIdle {
				t.Fatalf("expect lj.State to be %s. got=%s", StateLumberjackIdle, lj.State)
			}

			if lj.tree != nil {
				t.Fatalf("expect lj.tree to be nill. got=defined")
			}
			end = true
			break
		}

		if lj.State != StateLumberjackHitting {
			t.Fatalf("expect lj.State to be %s. got=%s", StateLumberjackHitting, lj.State)
		}

		if lj.tree == nil {
			t.Fatalf("expect lj.tree to be defined. got=nil")
		}

		if lj.tree.Health != 100-(actionCounter*lj.hit) {
			t.Fatalf("expect lj.tree.Health to be %d. got=%d", 100-(actionCounter*lj.hit), lj.tree.Health)
		}
		actionCounter += 1
	}

	if actionCounter != 2 {
		t.Fatalf("expect counter to be %d. got=%d", 2, actionCounter)
	}
	reset()
}
