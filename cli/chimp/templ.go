package main

import (
	"context"
	"time"

	"github.com/s00500/chimp/runner"
	log "github.com/s00500/env_logger"
	"github.com/spf13/cobra"
)

var templCmd = &cobra.Command{
	Use:   "templ",
	Short: "Watch .templ files and run templ generate on changes",
	Run: func(cmd *cobra.Command, args []string) {
		postHooks, _ := cmd.Flags().GetStringSlice("post-hook")
		noInteractive, _ := cmd.Flags().GetBool("no-interactive")
		debounce, _ := cmd.Flags().GetDuration("debounce")
		master, _ := cmd.Flags().GetBool("master")
		sync, _ := cmd.Flags().GetBool("sync")

		cfg := runner.CommandConfig{
			WatchExts: []string{".templ"},
			Command: runner.HookConfig{
				Command: "go",
				Args:    []string{"tool", "templ", "generate"},
			},
			PostHooks:   parseHooks(postHooks),
			Master:      master,
			Sync:        sync,
			Interactive: !noInteractive,
			Debounce:    debounce,
		}

		cr := runner.NewCommandRunner(cfg)
		if err := cr.Run(context.Background()); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	templCmd.Flags().StringSlice("post-hook", nil, "Command to run after templ generate")
	templCmd.Flags().Bool("no-interactive", false, "Disable keypress reload")
	templCmd.Flags().Bool("master", true, "Touch sync file after generation")
	templCmd.Flags().Bool("sync", false, "Watch sync file for reload signals")
	templCmd.Flags().Duration("debounce", 200*time.Millisecond, "Debounce duration")
	rootCmd.AddCommand(templCmd)
}
