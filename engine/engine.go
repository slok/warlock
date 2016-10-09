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
}
