//go:build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// dupFd duplicates a file descriptor.
func dupFd(fd uintptr) (uintptr, error) {
	newFd, err := syscall.Dup(int(fd))
	return uintptr(newFd), err
}

// notifyShutdownSignals registers for shutdown signals (SIGINT + SIGTERM on Unix).
func notifyShutdownSignals(ch chan<- os.Signal) {
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
}
