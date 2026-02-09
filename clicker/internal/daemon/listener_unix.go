//go:build !windows

package daemon

import "net"

// listen creates a Unix domain socket listener.
func listen(socketPath string) (net.Listener, error) {
	return net.Listen("unix", socketPath)
}
