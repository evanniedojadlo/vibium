//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// setSysProcAttr configures the child process to run in a new session (detached).
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
}
