package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newTraceCmd() *cobra.Command {
	traceCmd := &cobra.Command{
		Use:   "trace",
		Short: "Record browser traces (screenshots and snapshots)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start recording a trace",
		Example: `  vibium trace start --screenshots
  # Start trace with periodic screenshots

  vibium trace start --screenshots --snapshots
  # Start trace with screenshots and HTML snapshots`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			screenshots, _ := cmd.Flags().GetBool("screenshots")
			snapshots, _ := cmd.Flags().GetBool("snapshots")
			bidi, _ := cmd.Flags().GetBool("bidi")
			name, _ := cmd.Flags().GetString("name")

			if !oneshot {
				callArgs := map[string]interface{}{}
				if name != "" {
					callArgs["name"] = name
				}
				if screenshots {
					callArgs["screenshots"] = true
				}
				if snapshots {
					callArgs["snapshots"] = true
				}
				if bidi {
					callArgs["bidi"] = true
				}
				result, err := daemonCall("browser_trace_start", callArgs)
				if err != nil {
					printError(err)
					return
				}
				printResult(result)
				return
			}

			fmt.Fprintf(os.Stderr, "Error: trace command requires daemon mode\n")
			os.Exit(1)
		},
	}
	startCmd.Flags().Bool("screenshots", false, "Capture screenshots periodically")
	startCmd.Flags().Bool("snapshots", false, "Capture HTML snapshots")
	startCmd.Flags().Bool("bidi", false, "Record raw BiDi commands in the trace")
	startCmd.Flags().String("name", "", "Name for the trace")

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop recording and save the trace",
		Example: `  vibium trace stop
  # Save trace to trace.zip

  vibium trace stop -o my-trace.zip
  # Save trace to custom path`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.Flags().GetString("output")

			if !oneshot {
				callArgs := map[string]interface{}{}
				if output != "" {
					callArgs["path"] = output
				}
				result, err := daemonCall("browser_trace_stop", callArgs)
				if err != nil {
					printError(err)
					return
				}
				printResult(result)
				return
			}

			fmt.Fprintf(os.Stderr, "Error: trace command requires daemon mode\n")
			os.Exit(1)
		},
	}
	stopCmd.Flags().StringP("output", "o", "", "Output file path (default: trace.zip)")

	traceCmd.AddCommand(startCmd)
	traceCmd.AddCommand(stopCmd)
	return traceCmd
}
