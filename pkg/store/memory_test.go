package store

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	memoryStore := NewMemoryStore()

	// Test Set and Get methods
	key := "test"
	value := int64(100)
	memoryStore.Set(key, value)
	if storedValue := memoryStore.Get(key); storedValue != value {
		t.Errorf("Set and Get failed: expected %d but got %d", value, storedValue)
	}

	// Test Add method
	delta := int64(50)
	expectedSum := value + delta
	memoryStore.Add(key, delta)
	if storedValue := memoryStore.Get(key); storedValue != expectedSum {
		t.Errorf("Add failed: expected %d but got %d", expectedSum, storedValue)
	}

	// Test Delete method
	memoryStore.Del(key)
	if storedValue := memoryStore.Get(key); storedValue != 0 {
		t.Errorf("Del failed: expected 0 but got %d", storedValue)
	}
}
