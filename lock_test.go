package warlock

import (
	"fmt"
	"testing"
	"time"
)

const (
	key = "test_key"
)

// TestEngine is an engine only for testing purposes
type TestEngine struct {
	Key   string
	locks map[string]interface{}
	waitT time.Duration
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

func (t *TestEngine) Wait() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		time.Sleep(t.waitT)
		close(c)
	}()
	return c
}

// Tests

func TestLock(t *testing.T) {
	l := Warlock{
		Engine: newTestEngine(key),
	}
	if err := l.Lock(); err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
}

func TestLockBeingLocked(t *testing.T) {
	e := newTestEngine(key)
	l1 := Warlock{
		Engine: e,
	}

	l1.Lock()
	l2 := Warlock{
		Engine: e,
	}
	if err := l2.Lock(); err == nil {
		t.Errorf("Lock should return an error, it didn't")
	}
}

func TestUnlock(t *testing.T) {
	l := Warlock{
		Engine: newTestEngine(key),
	}
	l.Lock()
	if err := l.Unlock(); err != nil {
		t.Errorf("Unlock shouldn't return an error: %v", err)
	}
}

func TestUnlockWithoutLock(t *testing.T) {
	l := Warlock{
		Engine: newTestEngine(key),
	}
	if err := l.Unlock(); err == nil {
		t.Errorf("Unlock should return an error, it didn't")
	}
}

func TestLockWait(t *testing.T) {
	e := newTestEngine(key)
	e.waitT = 10 * time.Millisecond
	l := Warlock{
		Engine: e,
	}
	if err := l.Lock(); err != nil {
		t.Fatalf("Lock shouldn't return an error, it did: %v", err)
	}

	l2 := Warlock{
		Engine: e,
	}
	time.Sleep(1 * time.Millisecond)
	if err := l2.Lock(); err == nil {
		t.Fatalf("Lock should return an error, it didn't")
	}

	var unlocked bool
	go func() {
		// Wait until it unlocks
		<-l2.Wait()
		unlocked = true
	}()

	// Check we didn't received while blocked by f (the one with the lock)
	if unlocked {
		t.Errorf("The unlock signal shouldn't be received, it did")
	}

	l.Unlock()
	time.Sleep(11 * time.Millisecond)

	if !unlocked {
		t.Errorf("The unlock signal should be received, it didn't")
	}

}
