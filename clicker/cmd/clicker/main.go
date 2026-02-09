package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vibium/clicker/internal/log"
	"github.com/vibium/clicker/internal/process"
)

var version = "dev"

// Global flags
var (
	headless   bool
	waitOpen   int
	waitClose  int
	verbose    bool
	jsonOutput bool
	oneshot    bool
)

func main() {
	// Setup signal handler to cleanup on Ctrl+C
	process.SetupSignalHandler()

	rootCmd := &cobra.Command{
		Use:   "clicker",
		Short: "Browser automation for AI agents and humans",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Enable logging only if --verbose is used
			if verbose {
				log.Setup(log.LevelVerbose)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add global flags for browser commands
	rootCmd.PersistentFlags().BoolVar(&headless, "headless", false, "Hide browser window (visible by default)")
	rootCmd.PersistentFlags().IntVar(&waitOpen, "wait-open", 0, "Seconds to wait after navigation for page to load")
	rootCmd.PersistentFlags().IntVar(&waitClose, "wait-close", 0, "Seconds to keep browser open before closing")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&oneshot, "oneshot", false, "One-shot mode (no daemon, launch+execute+teardown)")

	// VIBIUM_ONESHOT=1 env var forces oneshot mode (for tests)
	if os.Getenv("VIBIUM_ONESHOT") == "1" {
		oneshot = true
	}

	// Register all commands
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newPathsCmd())
	rootCmd.AddCommand(newInstallCmd())
	rootCmd.AddCommand(newLaunchTestCmd())
	rootCmd.AddCommand(newWSTestCmd())
	rootCmd.AddCommand(newBiDiTestCmd())
	rootCmd.AddCommand(newNavigateCmd())
	rootCmd.AddCommand(newScreenshotCmd())
	rootCmd.AddCommand(newEvalCmd())
	rootCmd.AddCommand(newFindCmd())
	rootCmd.AddCommand(newClickCmd())
	rootCmd.AddCommand(newTypeCmd())
	rootCmd.AddCommand(newCheckActionableCmd())
	rootCmd.AddCommand(newServeCmd())
	rootCmd.AddCommand(newMCPCmd())
	rootCmd.AddCommand(newDaemonCmd())

	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Clicker v{{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
