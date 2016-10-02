package engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

// File file lock will implement a distributed lock using a shared filesystem
type File struct {
	Path   string
	TTL    time.Duration
	Expire bool
}

// Lock will lock using a simple file
func (f *File) Lock(key string) error {
	// Check locked first
	locked, err := f.Locked(key)
	if err != nil {
		return err
	}
	if locked {
		return fmt.Errorf("already locked")
	}

	// Lock by creating the key and setting the TTL on the file
	pathKey := path.Join(f.Path, key)
	now := time.Now().UTC()
	t := now.Add(f.TTL)
	b := []byte(fmt.Sprintf("%d", t.Unix()))
	if err := ioutil.WriteFile(pathKey, b, 0644); err != nil {
		return err
	}
	return nil
}

// Unlock unlocks a defined key
func (f *File) Unlock(key string) error {
	// Check locked first
	locked, err := f.Locked(key)
	if err != nil {
		return err
	}
	if !locked {
		return fmt.Errorf("not locked previously")
	}

	// Unlock removing the key
	pathKey := path.Join(f.Path, key)
	if err := os.Remove(pathKey); err != nil {
		return err
	}
	return nil
}

// Locked checks if the key is locked
func (f *File) Locked(key string) (bool, error) {
	pathKey := path.Join(f.Path, key)
	if _, err := os.Stat(pathKey); os.IsNotExist(err) {
		return false, nil
	}

	// Check TTL on file
	d, err := ioutil.ReadFile(pathKey)
	if err != nil {
		return true, err
	}
	i, err := strconv.Atoi(string(d))
	if err != nil {
		return true, err
	}
	t := time.Unix(int64(i), 0)
	now := time.Now().UTC()
	if now.Before(t) {
		return true, nil
	}

	return false, nil
}
