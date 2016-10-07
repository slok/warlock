package engine

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/slok/warlock/log"
)

// File file lock will implement a distributed lock using a shared filesystem
type File struct {
	Key    string
	Path   string
	TTL    time.Duration
	Expire bool
	ticker *time.Ticker

	pathKey string
	waiter  chan struct{}
}

// Lock will lock using a simple file
func (f *File) Lock() error {
	// Check locked first
	locked, err := f.Locked()
	if err != nil {
		return err
	}
	if locked {
		return fmt.Errorf("already locked")
	}

	// Lock by creating the key and setting the TTL on the file
	if err = f.renew(); err != nil {
		return err
	}

	// If don't expire then we need to renew the key before the TTL
	if !f.Expire {
		// every half of the time of the TTL renew
		f.ticker = time.NewTicker(f.TTL / 2)
		go func() {
			for range f.ticker.C {
				f.renew()
			}
		}()
	}

	return nil
}

// renew will renew the ttl of the lock
func (f *File) renew() error {
	pathKey := f.getPathKey()
	now := time.Now().UTC()
	t := now.Add(f.TTL)
	b := []byte(fmt.Sprintf("%d", t.UnixNano()))
	if err := ioutil.WriteFile(pathKey, b, 0644); err != nil {
		return err
	}
	return nil
}

// Unlock unlocks a defined key
func (f *File) Unlock() error {
	pathKey := f.getPathKey()
	// Check locked first
	locked, err := f.Locked()
	if err != nil {
		return err
	}
	if !locked {
		return fmt.Errorf("not locked previously")
	}

	// Unlock removing the key
	if err := os.Remove(pathKey); err != nil {
		return err
	}
	// Stop the renewer
	if !f.Expire {
		f.ticker.Stop()
	}
	return nil
}

// Locked checks if the key is locked
func (f *File) Locked() (bool, error) {
	pathKey := f.getPathKey()
	if _, err := os.Stat(pathKey); os.IsNotExist(err) {
		return false, nil
	}

	// Check TTL on file
	d, err := ioutil.ReadFile(pathKey)
	if err != nil {
		return true, err
	}

	// If empty means that the file was created by some other process and the
	// timestamp set wasn't finished writing, to ensure atomicity in this
	// operation we wait random time and check again
	if len(string(d)) == 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)
		d, err = ioutil.ReadFile(pathKey)
		if err != nil {
			return true, err
		}
	}

	i, err := strconv.Atoi(string(d))
	if err != nil {
		return true, err
	}
	t := time.Unix(0, int64(i))
	now := time.Now().UTC()
	if now.Before(t) {
		return true, nil
	}

	return false, nil
}

// Wait will return a channel that will be blocked until
func (f *File) Wait() <-chan struct{} {
	// If already waiting return the same channel
	if f.waiter != nil {
		return f.waiter
	}

	// Create a goroutine checking the status and return the waiting channel
	f.waiter = make(chan struct{})
	go func() {
		// check every TTL
		t := time.NewTicker(f.TTL)

		for range t.C {
			locked, err := f.Locked()
			if err != nil {
				log.Logger.Error(err.Error())
			}
			// Stop ticker and close channel to free the wait signal
			if !locked {
				t.Stop()
				close(f.waiter)
				f.waiter = nil
			}
		}
		log.Logger.Info("4")
	}()

	return f.waiter
}

// Internal method so we don't generate the path key on every check
func (f *File) getPathKey() string {
	if f.pathKey == "" {
		f.pathKey = path.Join(f.Path, f.Key)
	}

	return f.pathKey
}
