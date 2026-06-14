package entity

import (
	"github/teohen/mgm-tto/building"
	"testing"
)

func TestStoreInventory(t *testing.T) {
	storage := building.NewStorage(5, 5)
	sto := NewStorager(100)
	sto.inventory = 50
	sto.weight = 100

	done := sto.ExecuteAction(storage)
	if !done {
		t.Fatal("expected ExecuteAction to return true")
	}

	if sto.inventory != 0 {
		t.Fatalf("expect sto.inventory to be %d. got=%d", 0, sto.inventory)
	}

	if sto.weight != 0 {
		t.Fatalf("expect sto.weight to be %d. got=%d", 0, sto.weight)
	}

	if sto.storage != nil {
		t.Fatal("expect sto.storage to be nil")
	}

	if storage.Wood != 100 {
		t.Fatalf("expect storage.Wood to be %d. got=%d", 100, storage.Wood)
	}
}
