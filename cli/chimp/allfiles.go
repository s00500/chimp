package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

type TemplateData struct {
	ProjectName string
}

type FileDef struct {
	InPath      string // Path of the file in the embedded
	OutPath     string // path of the file in the final location
	UseTemplate bool
}

var AllFiles map[string]FileDef = map[string]FileDef{
	"Makefile": {
		InPath:      "templates/Makefile",
		OutPath:     "Makefile",
		UseTemplate: false,
	},
	"CSS": {
		InPath:      "templates/input.css",
		OutPath:     "css/input.css",
		UseTemplate: false,
	},
	"BaseCoatCSS": {
		InPath:      "templates/basecoat.css",
		OutPath:     "css/basecoat.css",
		UseTemplate: false,
	},
	"state.go": {
		InPath:      "templates/state",
		OutPath:     "state/state.go",
		UseTemplate: true,
	},
	"main.go": {
		InPath:      "templates/main",
		OutPath:     "main.go",
		UseTemplate: true,
	},
	"gitignore": {
		InPath:      "templates/gitignore",
		OutPath:     ".gitignore",
		UseTemplate: true,
	},
	"layout": {
		InPath:      "templates/layout_templ",
		OutPath:     "templates/layout.templ",
		UseTemplate: true,
	},
}

func (f FileDef) Render(basePath string, data TemplateData) error {
	if fileExists(filepath.Join(basePath, f.OutPath)) {
		prompt := promptui.Select{
			Label: fmt.Sprintf("File %s exists, should it be overwritten ?", f.OutPath),
			Items: []string{"No", "Yes"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "No" {
			return nil
		}
	}

	if f.UseTemplate {
		return WriteTemplate(f.InPath, filepath.Join(basePath, f.OutPath), data)
	} else {
		return WriteEmbedded(f.InPath, filepath.Join(basePath, f.OutPath))
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
