//go:build !windows

package browser

import (
	"os/exec"
	"syscall"
	"time"
)

// setProcGroup sets the process group for the command (Unix only).
func setProcGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killByPid sends SIGKILL to a process by PID.
func killByPid(pid int) {
	syscall.Kill(pid, syscall.SIGKILL)
}

// waitForProcessesDead polls until all PIDs have exited or timeout is reached.
func waitForProcessesDead(pids []int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		allDead := true
		for _, pid := range pids {
			if syscall.Kill(pid, 0) == nil {
				allDead = false
				break
			}
		}
		if allDead {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
}
