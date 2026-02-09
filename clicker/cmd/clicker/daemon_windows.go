//go:build windows

package main

import "os/exec"

// setSysProcAttr configures the child process for Windows.
// TODO: Use CREATE_NEW_PROCESS_GROUP when Windows daemon support is added.
func setSysProcAttr(cmd *exec.Cmd) {
	// No-op on Windows for now
}
