package warlock

import (
	"fmt"
	"testing"
)

// TestEngine is an engine only for testing purposes
type TestEngine struct {
	locks map[string]interface{}
}

func newTestEngine() *TestEngine {
	return &TestEngine{
		locks: make(map[string]interface{}),
	}
}

func (t *TestEngine) Lock(key string) error {
	if _, ok := t.locks[key]; ok {
		return fmt.Errorf("already locked")
	}

	t.locks[key] = nil

	return nil
}

func (t *TestEngine) Unlock(key string) error {
	if _, ok := t.locks[key]; !ok {
		return fmt.Errorf("not locked")
	}

	delete(t.locks, key)

	return nil
}

func (t *TestEngine) Locked(key string) (bool, error) {
	if _, ok := t.locks[key]; ok {
		return true, nil
	}

	return false, nil
}

// Tests

func TestLock(t *testing.T) {
	key := "test_key"
	l := Warlock{
		Key:    key,
		Engine: newTestEngine(),
	}
	if err := l.Lock(); err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
}

func TestLockBeingLocked(t *testing.T) {
	key := "test_key"
	e := newTestEngine()
	l1 := Warlock{
		Key:    key,
		Engine: e,
	}

	l1.Lock()
	l2 := Warlock{
		Key:    key,
		Engine: e,
	}
	if err := l2.Lock(); err == nil {
		t.Errorf("Lock should return an error, it didn't")
	}
}

func TestUnlock(t *testing.T) {
	key := "test_key"
	l := Warlock{
		Key:    key,
		Engine: newTestEngine(),
	}
	l.Lock()
	if err := l.Unlock(); err != nil {
		t.Errorf("Unlock shouldn't return an error: %v", err)
	}
}

func TestUnlockWithoutLock(t *testing.T) {
	key := "test_key"
	l := Warlock{
		Key:    key,
		Engine: newTestEngine(),
	}
	if err := l.Unlock(); err == nil {
		t.Errorf("Unlock should return an error, it didn't")
	}
}
