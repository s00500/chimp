package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	fmt.Println("Scanning .templ files for Tailwind classes...")

	classes := make(map[string]struct{})
	// Match static class="..." attributes
	classRegex := regexp.MustCompile(`class="([^"]*)"`)
	// Match templ.KV("classes", ...) calls
	kvRegex := regexp.MustCompile(`templ\.KV\("([^"]*)"`)
	// Match templ class={ ... } expressions (the class attribute, not
	// data-attr:class / data-bind:class — hence the leading whitespace).
	// (?s) lets the body span multiple lines; .*? stops at the first closing brace.
	classExprRegex := regexp.MustCompile(`(?s)\sclass=\{(.*?)\}`)
	// Match double-quoted string literals, used to pull class strings out of
	// class={ "literal", config.Class } expression bodies.
	literalRegex := regexp.MustCompile(`"([^"]*)"`)

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".templ") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		for _, regex := range []*regexp.Regexp{classRegex, kvRegex} {
			matches := regex.FindAllSubmatch(content, -1)
			for _, match := range matches {
				if len(match) > 1 {
					classStr := string(match[1])
					for _, class := range strings.Fields(classStr) {
						classes[class] = struct{}{}
					}
				}
			}
		}

		// Pull class strings out of templ class={ "literal", ... } expressions,
		// including any literals inside templ.KV(...) within the same body.
		for _, expr := range classExprRegex.FindAllSubmatch(content, -1) {
			if len(expr) < 2 {
				continue
			}
			for _, lit := range literalRegex.FindAllSubmatch(expr[1], -1) {
				if len(lit) > 1 {
					for _, class := range strings.Fields(string(lit[1])) {
						classes[class] = struct{}{}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to walk directory: %v\n", err)
		os.Exit(1)
	}

	// Sort classes
	sorted := make([]string, 0, len(classes))
	for class := range classes {
		sorted = append(sorted, class)
	}
	sort.Strings(sorted)

	// Write to cli/chimp/templates for embedding in CLI
	output := strings.Join(sorted, "\n") + "\n"
	outputPath := filepath.Join("cli", "chimp", "templates", "classes.txt")
	if err := os.WriteFile(outputPath, []byte(output), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Extracted %d unique classes to %s\n", len(sorted), outputPath)
}
