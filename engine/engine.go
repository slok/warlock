package engine

// Engine describes the interface needed to implement by the engines able to
// be locks
type Engine interface {
	// Lock locks a defined key
	Lock() error

	// Unlock unlocks a defined key
	Unlock() error

	// Locked checks if the key is locked
	Locked() (bool, error)

	// Wait returns a channel that will receive a signal when the lock is released
	Wait() <-chan struct{}
}
