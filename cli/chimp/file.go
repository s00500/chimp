package main

import (
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/manifoldco/promptui"
	log "github.com/s00500/env_logger"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file [filetemplate]",
	Short: "Create a new project scaffold",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		basePath := filepath.Join(".") // we basically assume we are in the root of the project
		wd, err := os.Getwd()
		log.MustFatal(log.Wrap(err, "on getting current working directory"))

		data := TemplateData{ProjectName: filepath.Base(wd)}

		fileList := append([]string{"All Files"}, slices.Collect(maps.Keys(AllFiles))...)

		prompt := promptui.Select{
			Label: "Select a file template",
			Items: fileList,
		}

		result := ""

		if len(args) > 0 {
			result = args[0]
		} else {
			_, result, err = prompt.Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		switch result {
		case "All Files", "all":
			for _, f := range AllFiles {
				err := f.Render(basePath, data)
				log.MustFatal(log.Wrap(err, "on write file"))
			}
			log.Infof("All Files created ✅")
		default:
			f, ok := AllFiles[result]
			if !ok {
				log.Fatal("nonexisting file")
			}

			err := f.Render(basePath, data)
			log.MustFatal(log.Wrap(err, "on write file"))
		}

		log.Infof("File '%s' created ✅", result)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
