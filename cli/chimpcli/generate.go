package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"embed"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var templateFs embed.FS

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "Create a new project scaffold",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		basePath := filepath.Join(".", projectName)
		data := map[string]string{"ProjectName": projectName}

		dirs := []string{
			"handler",
			"templates",
			"static",
			"state",
			"css",
		}

		for _, dir := range dirs {
			err := os.MkdirAll(filepath.Join(basePath, dir), os.ModePerm)
			if err != nil {
				log.Fatalf("Error creating directory %s: %v", dir, err)
			}
		}

		writeTemplate("templates/gitignore.tmpl", filepath.Join(basePath, ".gitignore"), data)
		//writeTemplate("templates/Makefile.tmpl", filepath.Join(basePath, "Makefile"), data)
		//writeTemplate("templates/config.yaml.tmpl", filepath.Join(basePath, "config", "config.yaml"), data)

		fmt.Printf("✅ Project '%s' created successfully!\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func writeTemplate(templatePath, outPath string, data any) {
	tmpl, err := template.ParseFS(templateFs, templatePath)
	if err != nil {
		log.Fatalf("Failed to parse embedded template %s: %v", templatePath, err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", outPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		log.Fatalf("Error writing template to %s: %v", outPath, err)
	}
}

func WriteEmbedded(filename, dst string) error {
	in, err := templateFs.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()

	// Ensure the destination folder exists
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return os.Chmod(dst, 0755)
}
