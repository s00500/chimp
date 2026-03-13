package main

import (
	"github.com/s00500/chimp/runner"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Touch the sync file to trigger reload in --sync consumers",
	Run: func(cmd *cobra.Command, args []string) {
		runner.TouchSyncFile()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
