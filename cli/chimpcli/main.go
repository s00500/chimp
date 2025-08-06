package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CLI Helper to generate boilerplate from templates

type TemplateData struct {
	ProjectName string
}

var rootCmd = &cobra.Command{
	Use:   "chimp",
	Short: "chimp cli can help you scaffold your project and add parts",
	Long:  `chimp cli can help you scaffold your project and add parts`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
