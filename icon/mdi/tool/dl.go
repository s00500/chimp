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

	url := "https://github.com/Templarian/MaterialDesign/archive/refs/heads/master.zip"
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

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		panic(err)
	}

	mdiIcons := make(map[string]string)
	for _, f := range zipReader.File {
		if strings.HasPrefix(f.Name, "MaterialDesign-master/svg/") && strings.HasSuffix(f.Name, ".svg") {
			iconName := strings.TrimPrefix(f.Name, "MaterialDesign-master/svg/")
			iconName = strings.TrimSuffix(iconName, ".svg")

			rc, err := f.Open()
			if err != nil {
				continue
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				fmt.Println("Error reading", iconName, ":", err)
				continue
			}
			mdiIcons[iconName] = string(content)
		}
	}

	// Output map as Go source code
	fmt.Fprintln(buffer, "package icon")
	fmt.Fprintln(buffer, "")
	fmt.Fprintln(buffer, "import \"github.com/s00500/chimp/icon/base\"")

	fmt.Fprintf(buffer, "")
	for name, svg := range mdiIcons {
		safeSVG := strings.ReplaceAll(svg, "`", "` + \"`\" + `") // Escape backticks
		safeSVG = strings.TrimSuffix(safeSVG, "</svg>")
		safeSVG = strings.TrimPrefix(safeSVG, `<svg xmlns="http://www.w3.org/2000/svg" id="mdi-`+name+`" viewBox="0 0 24 24">`)
		safeSVG = strings.TrimSpace(safeSVG)

		varName := toPascalCase(name)
		fmt.Fprintf(buffer, "const Mdi%s base.IconBase = `%s`\n\n", varName, safeSVG)
	}

	// go save the buffer
	err = os.WriteFile("../../mdi-icons.go", buffer.Bytes(), 0755)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

// toPascalCase converts kebab/underscore/mixed-case names to PascalCase.
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
