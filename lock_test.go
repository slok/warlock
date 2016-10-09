package warlock

import (
	"fmt"
	"testing"
)

const (
	key = "test_key"
)

// TestEngine is an engine only for testing purposes
type TestEngine struct {
	Key   string
	locks map[string]interface{}
}

func newTestEngine(key string) *TestEngine {
	return &TestEngine{
		Key:   key,
		locks: make(map[string]interface{}),
	}
}

func (t *TestEngine) Lock() error {
	if _, ok := t.locks[t.Key]; ok {
		return fmt.Errorf("already locked")
	}

	t.locks[t.Key] = nil

	return nil
}

func (t *TestEngine) Unlock() error {
	if _, ok := t.locks[t.Key]; !ok {
		return fmt.Errorf("not locked")
	}

	delete(t.locks, t.Key)

	return nil
}

func (t *TestEngine) Locked() (bool, error) {
	if _, ok := t.locks[t.Key]; ok {
		return true, nil
	}

	return false, nil
}

// Tests

func TestLock(t *testing.T) {
	l := Warlock{
		Key:    key,
		Engine: newTestEngine(key),
	}
	if err := l.Lock(); err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
}

func TestLockBeingLocked(t *testing.T) {
	e := newTestEngine(key)
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
	l := Warlock{
		Key:    key,
		Engine: newTestEngine(key),
	}
	l.Lock()
	if err := l.Unlock(); err != nil {
		t.Errorf("Unlock shouldn't return an error: %v", err)
	}
}

func TestUnlockWithoutLock(t *testing.T) {
	l := Warlock{
		Key:    key,
		Engine: newTestEngine(key),
	}
	if err := l.Unlock(); err == nil {
		t.Errorf("Unlock should return an error, it didn't")
	}
}
