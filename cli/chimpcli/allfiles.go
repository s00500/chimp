package main

import "path/filepath"

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
	"Layout template": {
		InPath:      "templates/layout_templ",
		OutPath:     "templates/layout.templ",
		UseTemplate: true,
	},
}

func (f FileDef) Render(basePath string, data TemplateData) error {
	if f.UseTemplate {
		return WriteTemplate(f.InPath, filepath.Join(basePath, f.OutPath), data)
	} else {
		return WriteEmbedded(f.InPath, f.OutPath)
	}
}
