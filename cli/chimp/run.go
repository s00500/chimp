package main

import (
	"context"
	"strings"
	"time"

	"github.com/s00500/chimp/runner"
	log "github.com/s00500/env_logger"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [flags] [-- binary-args...]",
	Short: "Build and run the project with live reload on file changes",
	Long: `Builds and runs the Go project in the current directory with live reload.
File changes trigger automatic rebuilds. Press Enter for manual reload.

Examples:
  chimp run
  chimp run -- -port 8080
  chimp run --master
  chimp run --sync
  chimp run --pre-hook "echo reloading" --ext .templ`,
	Run: func(cmd *cobra.Command, args []string) {
		exts, _ := cmd.Flags().GetStringSlice("ext")
		ignoreDirs, _ := cmd.Flags().GetStringSlice("ignore")
		master, _ := cmd.Flags().GetBool("master")
		sync, _ := cmd.Flags().GetBool("sync")
		preHooks, _ := cmd.Flags().GetStringSlice("pre-hook")
		postHooks, _ := cmd.Flags().GetStringSlice("post-hook")
		noInteractive, _ := cmd.Flags().GetBool("no-interactive")
		debounce, _ := cmd.Flags().GetDuration("debounce")

		// Everything after "--" goes to the binary
		binaryArgs := cmd.ArgsLenAtDash()
		var passArgs []string
		if binaryArgs >= 0 {
			passArgs = args[binaryArgs:]
		}

		cfg := runner.Config{
			Args:        passArgs,
			WatchExts:   exts,
			IgnoreDirs:  ignoreDirs,
			Master:      master,
			Sync:        sync,
			PreHooks:    parseHooks(preHooks),
			PostHooks:   parseHooks(postHooks),
			Interactive: !noInteractive,
			Debounce:    debounce,
		}

		r := runner.New(cfg)
		if err := r.Run(context.Background()); err != nil {
			log.Fatal(err)
		}
	},
}

func parseHooks(cmds []string) []runner.HookConfig {
	hooks := make([]runner.HookConfig, 0, len(cmds))
	for _, c := range cmds {
		parts := strings.Fields(c)
		if len(parts) == 0 {
			continue
		}
		hooks = append(hooks, runner.HookConfig{
			Command: parts[0],
			Args:    parts[1:],
		})
	}
	return hooks
}

func init() {
	runCmd.Flags().StringSliceP("ext", "e", nil, "Additional extensions to watch (e.g. .templ,.html)")
	runCmd.Flags().StringSliceP("ignore", "i", nil, "Additional directories to ignore")
	runCmd.Flags().Bool("master", false, "Touch sync file on reload (producer)")
	runCmd.Flags().Bool("sync", false, "Watch sync file for reload signals (consumer)")
	runCmd.Flags().StringSlice("pre-hook", nil, "Command to run before restart")
	runCmd.Flags().StringSlice("post-hook", nil, "Command to run after restart")
	runCmd.Flags().Bool("no-interactive", false, "Disable keypress reload, pass stdin to child")
	runCmd.Flags().Duration("debounce", 200*time.Millisecond, "Debounce duration for file changes")

	rootCmd.AddCommand(runCmd)
}
