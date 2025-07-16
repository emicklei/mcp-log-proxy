//go:build windows

package filelock

import (
	"errors"

	"github.com/gofrs/flock"
)

type lockType uint32

const (
	readLock  lockType = 0
	writeLock lockType = 1
)

const (
	reserved = 0
	allBytes = ^uint32(0)
)

var currentLock *flock.Flock

func lock(f File, _ lockType) error {
	lockName := f.Name() + ".lock"
	fileLock := flock.New(lockName)
	locked, err := fileLock.TryLock()
	if err != nil {
		return err
	}
	if locked {
		currentLock = fileLock
		return nil // we got it
	}
	return errors.New("file is already locked: " + f.Name())
}

func unlock(_ File) error {
	if currentLock == nil {
		return errors.New("no lock to unlock")
	}
	defer func() {
		currentLock = nil // reset current lock after unlocking
	}()
	return currentLock.Unlock()
}
