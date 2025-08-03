package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	buffer := bytes.NewBuffer([]byte{})

	url := "https://github.com/lucide-icons/lucide/archive/refs/heads/main.zip"
	fmt.Println("Downloading Lucide icon set...")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		panic(err)
	}

	fmt.Println("Unzipping...")
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(buffer, "package icon")
	fmt.Fprintln(buffer, "")
	fmt.Fprintln(buffer, "import \"github.com/s00500/chimp/icon/base\"")

	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "lucide-main/icons/") && strings.HasSuffix(f.Name, ".svg") {
			iconName := path.Base(f.Name)
			name := strings.TrimSuffix(iconName, ".svg")
			varName := toPascalCase(name)

			rc, err := f.Open()
			if err != nil {
				fmt.Println("Skipping", name, ":", err)
				continue
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				fmt.Println("Error reading", name, ":", err)
				continue
			}

			safeSVG := strings.ReplaceAll(string(content), "`", "` + \"`\" + `")
			safeSVG = strings.ReplaceAll(safeSVG, "\n", "")

			safeSVG = strings.TrimPrefix(safeSVG, `<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round">`)
			safeSVG = strings.TrimSuffix(safeSVG, "</svg>")
			safeSVG = strings.TrimSpace(safeSVG)

			fmt.Fprintf(buffer, "const Lu%s base.IconBase = `%s`\n\n", varName, safeSVG)
		}
	}

	// go save the buffer
	err = os.WriteFile("../../lucide-icons.go", buffer.Bytes(), 0755)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func toPascalCase(name string) string {
	re := regexp.MustCompile(`[._\- ]+`)
	parts := re.Split(name, -1)
	for i, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}
