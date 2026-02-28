//go:build windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// dupFd duplicates a file descriptor using Windows DuplicateHandle.
func dupFd(fd uintptr) (uintptr, error) {
	proc, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, err
	}
	var dup syscall.Handle
	err = syscall.DuplicateHandle(
		proc,
		syscall.Handle(fd),
		proc,
		&dup,
		0,
		false,
		syscall.DUPLICATE_SAME_ACCESS,
	)
	if err != nil {
		return 0, err
	}
	return uintptr(dup), nil
}

// notifyShutdownSignals registers for shutdown signals (SIGINT only on Windows).
func notifyShutdownSignals(ch chan<- os.Signal) {
	signal.Notify(ch, os.Interrupt)
}
