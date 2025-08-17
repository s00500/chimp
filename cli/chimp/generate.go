package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"embed"

	log "github.com/s00500/env_logger"
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
		data := TemplateData{ProjectName: projectName}

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

		for _, f := range AllFiles {
			err := f.Render(basePath, data)
			log.MustFatal(log.Wrap(err, "on write file"))
		}

		init := exec.Command("go", "mod", "init", projectName)
		init.Dir = basePath
		init.CombinedOutput()

		for _, t := range AllTools {
			err := t.Install(basePath)
			log.MustFatal(log.Wrap(err, "on installing tool"))
		}

		tidy := exec.Command("go", "mod", "tidy")
		tidy.Dir = basePath
		tidy.CombinedOutput()

		fmt.Printf("✅ Project '%s' created successfully!\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func WriteTemplate(templatePath, outPath string, data TemplateData) error {
	tmpl, err := template.ParseFS(templateFs, templatePath)
	if err != nil {
		return fmt.Errorf("Failed to parse embedded template %s: %v", templatePath, err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("Error creating file %s: %v", outPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("Error writing template to %s: %v", outPath, err)
	}
	return nil
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
