// +build integration

package engine

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	testPath    = "/tmp"
	testKey     = "warlock_test"
	testPathKey = "/tmp/warlock_test"
)

func TestMain(m *testing.M) {
	flag.Parse()
	// Setup
	ec := m.Run()
	// Teardown
	//Delete key
	os.Remove(testPathKey)

	os.Exit(ec)
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func TestLockNoPreviousLock(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	f := File{
		Key:  testKey,
		Path: testPath,
		TTL:  1 * time.Second,
	}
	err := f.Lock()
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}

	if !fileExists(testPathKey) {
		t.Errorf("File should exist")
	}
}

func TestLockMultipleLocks(t *testing.T) {
	pathKeys := []string{}
	for i := 0; i < 10; i++ {
		pathKeys = append(pathKeys, fmt.Sprintf("%s-%d", testPathKey, i))
	}

	defer func() {
		for _, p := range pathKeys {
			os.Remove(p)
		}
	}()

	for i, p := range pathKeys {
		f := File{
			Key:  fmt.Sprintf("%s-%d", testKey, i),
			Path: testPath,
			TTL:  1 * time.Second,
		}
		err := f.Lock()
		if err != nil {
			t.Errorf("Lock shouldn't return an error: %v", err)
		}

		if !fileExists(p) {
			t.Errorf("File %s should exist", p)
		}
	}
}

func TestLockPreviousLock(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	f := File{
		Key:  testKey,
		Path: testPath,
		TTL:  1 * time.Second,
	}
	err := f.Lock()
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
	err = f.Lock()
	if err == nil {
		t.Errorf("Lock should return an error")
	}
}

func TestLockExpire(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	f := File{
		Key:    testKey,
		Path:   testPath,
		TTL:    10 * time.Millisecond,
		Expire: true,
	}
	err := f.Lock()
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
	if !fileExists(testPathKey) {
		t.Errorf("File should exist")
	}
	time.Sleep(f.TTL)
	err = f.Lock()
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
}

func TestLockNotExpire(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	f := File{
		Key:  testKey,
		Path: testPath,
		TTL:  50 * time.Millisecond,
	}
	err := f.Lock()
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
	if !fileExists(testPathKey) {
		t.Errorf("File should exist")
	}
	time.Sleep(f.TTL * 2)
	err = f.Lock()
	if err == nil {
		t.Errorf("Lock should return an error")
	}
}

func TestUnLockPreviousLock(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	f := File{
		Key:  testKey,
		Path: testPath,
		TTL:  1 * time.Second,
	}
	err := f.Lock()
	if err != nil {
		t.Fatalf("Lock shouldn't return an error: %v", err)
	}
	err = f.Unlock()
	if err != nil {
		t.Errorf("Unlock shouldn't return an error: %v", err)
	}

	if fileExists(testPathKey) {
		t.Errorf("File shouldn't exist")
	}
}

func TestUnLockNoPreviousLock(t *testing.T) {
	f := File{
		Key:  testKey,
		Path: testPath,
		TTL:  1 * time.Second,
	}
	err := f.Unlock()
	if err == nil {
		t.Errorf("Unlock should return an error")
	}
}

func TestLockWait(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	// Create one lock
	f := File{
		Key:  testKey,
		Path: testPath,
		TTL:  10 * time.Millisecond,
	}
	err := f.Lock()
	if err != nil {
		t.Fatalf("Lock shouldn't return an error: %v", err)
	}
	// Create a 2nd lock
	f2 := File{
		Key:  testKey,
		Path: testPath,
		TTL:  10 * time.Millisecond,
	}
	err = f2.Lock()
	if err == nil {
		t.Fatalf("Lock should return an error")
	}

	var unlocked bool
	go func() {
		// Wait until it unlocks
		<-f2.Wait()
		unlocked = true
	}()

	// Check we didn't received while blocked by f (the one with the lock)
	if unlocked {
		t.Errorf("The unlock signal shouldn't be received, it did")
	}

	f.Unlock()
	time.Sleep(f2.TTL * 5)

	if !unlocked {
		t.Errorf("The unlock signal should be received, it didn't")
	}
}
