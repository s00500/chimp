package main

import (
	"maps"
	"slices"

	"github.com/manifoldco/promptui"
	log "github.com/s00500/env_logger"
	"github.com/spf13/cobra"
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Install tools",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolList := append([]string{"All Tools"}, slices.Collect(maps.Keys(AllTools))...)

		prompt := promptui.Select{
			Label: "Select a tool to install",
			Items: toolList,
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
		case "All Tools", "all":
			for _, t := range AllTools {
				err := t.Install("")
				log.MustFatal(log.Wrap(err, "on installing tool"))
			}
			log.Infof("All Tools '%s' installed ✅, use it with go tool <toolname>", result)
		default:
			t, ok := AllTools[result]
			if !ok {
				log.Fatal("nonexisting tool")
			}

			err := t.Install("")
			log.MustFatal(log.Wrap(err, "on installing tool"))
		}

		log.Infof("Tool '%s' installed ✅, use it with go tool <toolname>", result)
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)
}
