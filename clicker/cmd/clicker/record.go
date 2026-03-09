package main

import (
	"github.com/spf13/cobra"
)

func newRecordCmd() *cobra.Command {
	recordCmd := &cobra.Command{
		Use:   "record",
		Short: "Record browser sessions (screenshots and snapshots)",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a recording",
		Example: `  vibium record start --screenshots
  # Start recording with periodic screenshots (JPEG, quality 0.5)

  vibium record start --screenshots --snapshots
  # Start recording with screenshots and HTML snapshots

  vibium record start --snapshots --format png
  # Use PNG format instead of JPEG (larger files, lossless)

  vibium record start --snapshots --quality 0.5
  # Lower JPEG quality for smaller recording files`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			screenshots, _ := cmd.Flags().GetBool("screenshots")
			snapshots, _ := cmd.Flags().GetBool("snapshots")
			bidi, _ := cmd.Flags().GetBool("bidi")
			name, _ := cmd.Flags().GetString("name")
			format, _ := cmd.Flags().GetString("format")
			quality, _ := cmd.Flags().GetFloat64("quality")

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
			if format != "jpeg" {
				callArgs["format"] = format
			}
			if quality != 0.5 {
				callArgs["quality"] = quality
			}
			result, err := daemonCall("browser_record_start", callArgs)
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	startCmd.Flags().Bool("screenshots", false, "Capture screenshots periodically")
	startCmd.Flags().Bool("snapshots", false, "Capture HTML snapshots")
	startCmd.Flags().Bool("bidi", false, "Record raw BiDi commands in the recording")
	startCmd.Flags().String("name", "", "Name for the recording")
	startCmd.Flags().String("format", "jpeg", "Screenshot format: jpeg or png")
	startCmd.Flags().Float64("quality", 0.5, "JPEG quality 0.0-1.0 (ignored for png)")

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop recording and save",
		Example: `  vibium record stop
  # Save recording to record.zip

  vibium record stop -o my-recording.zip
  # Save recording to custom path`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			output, _ := cmd.Flags().GetString("output")

			callArgs := map[string]interface{}{}
			if output != "" {
				callArgs["path"] = output
			}
			result, err := daemonCall("browser_record_stop", callArgs)
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	stopCmd.Flags().StringP("output", "o", "", "Output file path (default: record.zip)")

	recordCmd.AddCommand(startCmd)
	recordCmd.AddCommand(stopCmd)
	return recordCmd
}
