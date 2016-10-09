package warlock

import (
	"fmt"

	"github.com/slok/warlock/engine"
)

// Warlock reprensents the lock object that holds the lock
type Warlock struct {
	// The lock key that will identify the lock
	Key string

	// Engine will reprenset the locks engine
	Engine engine.Engine
}

// Lock locks the lock
func (w *Warlock) Lock() error {
	// Check if is already locked
	l, err := w.Engine.Locked()
	if err != nil {
		return err
	}
	if l {
		return fmt.Errorf("already locked")
	}

	// Lock
	if err = w.Engine.Lock(); err != nil {
		return err
	}

	return nil
}

// Unlock unlocks the lock
func (w *Warlock) Unlock() error {
	// If not locked then can't be unlocked
	l, err := w.Engine.Locked()
	if err != nil {
		return err
	}
	if !l {
		return fmt.Errorf("not locked")
	}

	// Unlock
	if err = w.Engine.Unlock(); err != nil {
		return err
	}

	return nil
}
