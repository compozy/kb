//go:build darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd

package corpus

import (
	"errors"
	"os"
	"syscall"
)

func lockStateFile(file *os.File) error {
	return flockState(file, syscall.LOCK_EX)
}

func unlockStateFile(file *os.File) error {
	return flockState(file, syscall.LOCK_UN)
}

func flockState(file *os.File, operation int) error {
	for {
		err := syscall.Flock(int(file.Fd()), operation)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}
