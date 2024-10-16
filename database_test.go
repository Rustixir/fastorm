package fastorm_test

import (
	"fastorm"
	"testing"
)

// Test Cases
func TestDb_SetGetCommit(t *testing.T) {
	db := fastorm.NewDatabase()
	// Start a new transaction
	tx := db.Begin(fastorm.ReadCommitted)

	// Test Set and Get in active transaction
	err := tx.Set("key1", "value1")
	if err != nil {
		t.Fatalf("expected no error on Set, got %v", err)
	}

	val, err := tx.Get("key1")
	if err != nil {
		t.Fatalf("expected no error on Get, got %v", err)
	}
	if val != "value1" {
		t.Fatalf("expected value1, got %v", val)
	}

	// Test Commit
	err = tx.Commit()
	if err != nil {
		t.Fatalf("expected no error on Commit, got %v", err)
	}

	// Check if data was committed
	tx2 := db.Begin(fastorm.ReadCommitted)
	val, err = tx2.Get("key1")
	if err != nil {
		t.Fatalf("expected no error on Get after Commit, got %v", err)
	}
	if val != "value1" {
		t.Fatalf("expected value1 after commit, got %v", val)
	}
}

func TestDb_Rollback(t *testing.T) {
	db := fastorm.NewDatabase()

	// Start a new transaction
	tx := db.Begin(fastorm.ReadCommitted)

	// Test Set
	err := tx.Set("key2", "value2")
	if err != nil {
		t.Fatalf("expected no error on Set, got %v", err)
	}

	// Test Rollback
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("expected no error on Rollback, got %v", err)
	}

	// Verify that data was not committed
	_, err = tx.Get("key2")
	if err == nil {
		t.Fatalf("expected error on Get after Rollback, got no error")
	}
}

func TestDb_Delete(t *testing.T) {
	db := fastorm.NewDatabase()

	// Start a new transaction
	tx := db.Begin(fastorm.ReadCommitted)

	// Commit some initial data
	err := tx.Set("key3", "value3")
	if err != nil {
		t.Fatalf("expected no error on Set, got %v", err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatalf("expected no error on Commit, got %v", err)
	}

	// Start a new transaction to delete the key
	tx2 := db.Begin(fastorm.ReadCommitted)

	// Test Delete
	err = tx2.Delete("key3")
	if err != nil {
		t.Fatalf("expected no error on Delete, got %v", err)
	}

	// Commit to delete
	err = tx2.Commit()
	if err != nil {
		t.Fatalf("expected no error on Commit after Delete, got %v", err)
	}

	// Verify the key was deleted
	_, err = tx.Get("key3")
	if err == nil {
		t.Fatalf("expected error on Get after Delete and Commit, got no error")
	}
}

func TestDb_Range(t *testing.T) {
	db := fastorm.NewDatabase()
	txn := db.Begin(fastorm.ReadCommitted)

	// Commit some initial data
	_ = txn.Set("keyA", "valueA")
	_ = txn.Set("keyB", "valueB")
	txn.Commit()

	// Start a new transaction
	tx := db.Begin(fastorm.ReadCommitted)

	// Test Range
	count := 0
	err := tx.Range(func(key fastorm.Key, value fastorm.Value) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error on Range, got %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 keys in Range, got %d", count)
	}
}
