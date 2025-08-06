package main

import (
	"github.com/manifoldco/promptui"
	log "github.com/s00500/env_logger"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file [filetemplate]",
	Short: "Create a new project scaffold",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// projectName := args[0]
		// data := map[string]string{"ProjectName": projectName}

		prompt := promptui.Select{
			Label: "Select a file template",
			Items: []string{"All Files", "Makefile", "css", "BaseCoatCSS", "main.go"},
		}
		result := ""
		var err error

		if len(args) > 0 {
			result = args[0]
		} else {
			_, result, err = prompt.Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		switch result {
		case "Makefile":
			err = WriteEmbedded("templates/Makefile", "Makefile")
			log.ShouldWrap(err, "on write file")
		case "css":
			err = WriteEmbedded("templates/input.css", "css/input.css")
			log.ShouldWrap(err, "on write file")
		case "BaseCoatCSS":
			err = WriteEmbedded("templates/basecoat.css", "css/basecoat.css")
			log.ShouldWrap(err, "on write file")
		case "main.go":
			err = WriteEmbedded("templates/main", "main.go")
			log.ShouldWrap(err, "on write file")
		case "All Files":

			err = WriteEmbedded("templates/Makefile", "Makefile")
			log.ShouldWrap(err, "on write file")
			err = WriteEmbedded("templates/input.css", "css/input.css")
			log.ShouldWrap(err, "on write file")
			err = WriteEmbedded("templates/basecoat.css", "css/basecoat.css")
			log.ShouldWrap(err, "on write file")
			err = WriteEmbedded("templates/main", "main.go")
			log.ShouldWrap(err, "on write file")
		default:
			log.Fatal("nonexisting file")
		}

		log.Infof("File '%s' created ✅", result)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
