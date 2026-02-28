//go:build windows

package browser

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

// platformChromeArgs returns Windows-specific Chrome launch arguments.
func platformChromeArgs() []string {
	// Chrome for Testing sandbox cannot access its own executable in AppData
	// due to Windows filesystem permission restrictions.
	return []string{"--no-sandbox"}
}

// setProcGroup is a no-op on Windows.
func setProcGroup(cmd *exec.Cmd) {
	// Windows doesn't use process groups the same way
}

// killByPid kills a process tree by PID on Windows.
func killByPid(pid int) {
	// /T kills the entire process tree, /F forces termination
	exec.Command("taskkill", "/T", "/F", "/PID", fmt.Sprintf("%d", pid)).Run()
}

// isProcessAlive checks whether a process with the given PID is still running.
func isProcessAlive(pid int) bool {
	// tasklist /FI filters by PID and /FO CSV gives parseable output.
	// Exit code 0 + output containing the PID means the process exists.
	out, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH").Output()
	if err != nil {
		return false
	}
	return len(out) > 0 && bytes.Contains(out, []byte(fmt.Sprintf("%d", pid)))
}

// waitForProcessesDead polls until all PIDs have exited or timeout is reached.
func waitForProcessesDead(pids []int, timeout time.Duration) {
	// Brief initial sleep to let the OS reap process table entries
	time.Sleep(50 * time.Millisecond)

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		allDead := true
		for _, pid := range pids {
			if isProcessAlive(pid) {
				allDead = false
				break
			}
		}
		if allDead {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
