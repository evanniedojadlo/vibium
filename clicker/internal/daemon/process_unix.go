//go:build !windows

package daemon

import "syscall"

// processExists checks if a process with the given PID exists.
func processExists(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || err == syscall.EPERM
}
