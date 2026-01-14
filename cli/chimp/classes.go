package main

import (
	"os"
	"path/filepath"

	log "github.com/s00500/env_logger"
	"github.com/spf13/cobra"
)

var classesCmd = &cobra.Command{
	Use:   "classes",
	Short: "Copy chimp's Tailwind classes to css/classes.txt",
	Run: func(cmd *cobra.Command, args []string) {
		content, err := templateFs.ReadFile("templates/classes.txt")
		log.MustFatal(log.Wrap(err, "failed to read embedded classes.txt"))

		// Ensure css directory exists
		err = os.MkdirAll("css", 0755)
		log.MustFatal(log.Wrap(err, "failed to create css directory"))

		// Write to file
		outputPath := filepath.Join("css", "classes.txt")
		err = os.WriteFile(outputPath, content, 0644)
		log.MustFatal(log.Wrap(err, "failed to write "+outputPath))

		log.Infof("Wrote %s", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(classesCmd)
}
