package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	buffer := bytes.NewBuffer([]byte{})

	url := "https://github.com/tailwindlabs/heroicons/archive/refs/heads/master.zip"
	fmt.Println("Downloading Heroicons...")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Unzipping...")
	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsing SVGs...")

	fmt.Fprintln(buffer, "package icon")
	fmt.Fprintln(buffer, "")
	fmt.Fprintln(buffer, "import \"github.com/s00500/chimp/icon/base\"")

	for _, f := range zipReader.File {
		if !strings.HasSuffix(f.Name, ".svg") {
			continue
		}

		// Target only `optimized/20/solid/` and `optimized/24/outline/` SVGs
		if !(strings.Contains(f.Name, "optimized/24/solid/") || strings.Contains(f.Name, "optimized/24/outline/")) {
			continue
		}

		parts := strings.Split(f.Name, "/")
		if len(parts) < 2 {
			continue
		}

		style := ""
		if strings.Contains(f.Name, "/solid/") {
			style = "Solid"
		} else if strings.Contains(f.Name, "/outline/") {
			style = "Outline"
		}

		// filename.svg → iconName
		filename := parts[len(parts)-1]
		iconName := strings.TrimSuffix(filename, ".svg")
		varName := toPascalCase(iconName) + style

		rc, err := f.Open()
		if err != nil {
			fmt.Println("Skipping", iconName, ":", err)
			continue
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			fmt.Println("Error reading", iconName, ":", err)
			continue
		}

		safeSVG := strings.ReplaceAll(string(content), "`", "` + \"`\" + `")
		safeSVG = strings.ReplaceAll(safeSVG, "\n", "")

		safeSVG = strings.TrimPrefix(safeSVG, `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true" data-slot="icon">`)
		safeSVG = strings.TrimSuffix(safeSVG, "</svg>")
		safeSVG = strings.TrimSpace(safeSVG)

		fmt.Fprintf(buffer, "const He%s base.IconBase = `%s`\n\n", varName, safeSVG)
	}

	// go save the buffer
	err = os.WriteFile("../../heroicon-icons.go", buffer.Bytes(), 0755)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

// Converts dashed/underscored names to PascalCase
func toPascalCase(name string) string {
	re := regexp.MustCompile(`[._\- ]+`)
	parts := re.Split(name, -1)
	for i, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}
