package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/daemon"
	"github.com/vibium/clicker/internal/paths"
)

func newConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect <url>",
		Short: "Connect to a remote browser",
		Long: `Stop any running daemon and start a new one connected to a remote
BiDi WebSocket endpoint. Subsequent commands will use the remote browser.

Set VIBIUM_CONNECT_API_KEY to send an Authorization: Bearer header.`,
		Example: `  vibium connect ws://remote:9515/session
  vibium go https://example.com
  vibium disconnect

  # With authentication
  export VIBIUM_CONNECT_API_KEY=my-api-key
  vibium connect wss://cloud.example.com/session`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			connectURL := args[0]

			// Stop existing daemon if running
			if daemon.IsRunning() {
				if err := daemon.Shutdown(); err != nil {
					fmt.Fprintf(os.Stderr, "Error stopping existing daemon: %v\n", err)
					os.Exit(1)
				}
				// Wait briefly for socket cleanup
				time.Sleep(200 * time.Millisecond)
			}

			// Clean stale files
			daemon.CleanStale()

			// Start a new daemon with --connect
			exe, err := os.Executable()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding executable: %v\n", err)
				os.Exit(1)
			}

			daemonArgs := []string{"daemon", "start", "--_internal", "--idle-timeout=30m",
				fmt.Sprintf("--connect=%s", connectURL)}
			if headless {
				daemonArgs = append(daemonArgs, "--headless")
			}

			// Forward API key from env
			_, envHeaders := connectFromEnv()
			for key, vals := range envHeaders {
				for _, v := range vals {
					daemonArgs = append(daemonArgs, fmt.Sprintf("--connect-header=%s: %s", key, v))
				}
			}

			child := exec.Command(exe, daemonArgs...)
			child.Stdout = nil
			child.Stderr = nil
			child.Stdin = nil
			setSysProcAttr(child)

			if err := child.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Error starting daemon: %v\n", err)
				os.Exit(1)
			}

			// Wait for socket
			socketPath, _ := paths.GetSocketPath()
			if err := waitForSocket(socketPath, 5*time.Second); err != nil {
				fmt.Fprintf(os.Stderr, "Daemon failed to start: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Connected to %s (daemon pid %d)\n", connectURL, child.Process.Pid)
		},
	}
}

func newDisconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from a remote browser",
		Long: `Stop the daemon. The next command will auto-start a new daemon
with a local browser (unless VIBIUM_CONNECT_URL is set).`,
		Run: func(cmd *cobra.Command, args []string) {
			if !daemon.IsRunning() {
				fmt.Println("No daemon running.")
				return
			}

			if err := daemon.Shutdown(); err != nil {
				fmt.Fprintf(os.Stderr, "Error stopping daemon: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Disconnected.")
		},
	}
}
