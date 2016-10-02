package engine

// Engine describes the interface needed to implement by the engines able to
// be locks
type Engine interface {
	// Lock locks a defined key
	Lock(key string) error

	// Unlock unlocks a defined key
	Unlock(key string) error

	// Locked checks if the key is locked
	Locked(key string) (bool, error)
}
