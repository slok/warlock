// +build integration

package engine

import (
	"flag"
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
	if _, err := os.Stat(testPathKey); os.IsNotExist(err) {
		return false
	}
	return true
}

func TestLockNoPreviousLock(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	e := File{
		Path: testPath,
	}
	err := e.Lock(testKey)
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}

	if !fileExists(testPathKey) {
		t.Errorf("File should exist")
	}
}

func TestLockPreviousLock(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	e := File{
		Path: testPath,
		TTL:  1 * time.Second,
	}
	err := e.Lock(testKey)
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
	err = e.Lock(testKey)
	if err == nil {
		t.Errorf("Lock should return an error")
	}
}

func TestLockExpire(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	e := File{
		Path:   testPath,
		TTL:    10 * time.Millisecond,
		Expire: true,
	}
	err := e.Lock(testKey)
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
	if !fileExists(testPathKey) {
		t.Errorf("File should exist")
	}
	time.Sleep(e.TTL)
	err = e.Lock(testKey)
	if err != nil {
		t.Errorf("Lock shouldn't return an error: %v", err)
	}
}

func TestUnLockPreviousLock(t *testing.T) {
	defer func() { os.Remove(testPathKey) }()
	e := File{
		Path: testPath,
		TTL:  1 * time.Second,
	}
	e.Lock(testKey)
	err := e.Unlock(testKey)
	if err != nil {
		t.Errorf("Unlock shouldn't return an error: %v", err)
	}

	if fileExists(testPathKey) {
		t.Errorf("File shouldn't exist")
	}
}

func TestUnLockNoPreviousLock(t *testing.T) {
	e := File{
		Path: testPath,
	}
	err := e.Unlock(testKey)
	if err == nil {
		t.Errorf("Unlock should return an error")
	}
}
