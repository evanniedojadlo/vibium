//go:build windows

package daemon

import (
	"fmt"
	"net"
)

// listen creates a named pipe listener on Windows.
// TODO: Implement named pipe support for Windows.
func listen(socketPath string) (net.Listener, error) {
	return nil, fmt.Errorf("daemon mode is not yet supported on Windows")
}
