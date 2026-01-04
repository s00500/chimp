package main

import (
	"fmt"
	"os/exec"
)

type ToolDef struct {
	ToolName string
	Version  string // path of the file in the final location
}

var AllTools = map[string]ToolDef{
	"Tailwind": {
		ToolName: "github.com/hookenz/gotailwind/v4",
		Version:  "latest",
	},
	"Air": {
		ToolName: "github.com/air-verse/air",
		Version:  "latest",
	},
	"Templ": {
		ToolName: "github.com/a-h/templ/cmd/templ",
		Version:  "latest",
	},
	// "Mage": {
	// 	ToolName: "github.com/magefile/mage",
	// 	Version:  "latest",
	// },
}

func (t ToolDef) Install(basePath string) error {
	install := exec.Command("go", "get", "--tool", fmt.Sprintf("%s@%s", t.ToolName, t.Version))
	if basePath != "" {
		install.Dir = basePath
	}
	out, err := install.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed install of %s: %s, %w", t.ToolName, string(out), err)
	}
	return nil
}
