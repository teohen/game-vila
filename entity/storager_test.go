package entity

import (
	"github/teohen/mgm-tto/building"
	"testing"
)

func TestStoreInventory(t *testing.T) {
	storage := building.NewStorage(5, 5)
	sto := NewStorager(100)
	sto.Inventory = 50

	done := sto.ExecuteAction(storage)
	if !done {
		t.Fatal("expected ExecuteAction to return true")
	}

	if sto.Inventory != 0 {
		t.Fatalf("expect sto.Inventory to be %d. got=%d", 0, sto.Inventory)
	}

	if sto.CalculateWeight() != 0 {
		t.Fatalf("expect sto.Weight to be %d. got=%d", 0, sto.CalculateWeight())
	}

	if sto.storage != nil {
		t.Fatal("expect sto.storage to be nil")
	}
}
