//go:build windows

package daemon

import "os"

// processExists checks if a process with the given PID exists.
func processExists(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// On Windows, FindProcess always succeeds. We rely on the PID file being present.
	_ = p
	return true
}
