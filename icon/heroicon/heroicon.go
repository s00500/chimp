package heroicon

import (
	"strings"

	"codeberg.org/jmansfield/heroiconsgo"
	"github.com/a-h/templ"
)

func Icon(id, classes string) templ.Component {
	return templ.Raw(strings.Replace(heroiconsgo.Get("outline", id), "<svg", `<svg class="`+classes+`" `, 1))
}

func SolidIcon(id, classes string) templ.Component {
	return templ.Raw(strings.Replace(heroiconsgo.Get("solid", id), "<svg", `<svg class="`+classes+`" `, 1))
}

func MiniIcon(id, classes string) templ.Component {
	return templ.Raw(strings.Replace(heroiconsgo.Get("mini", id), "<svg", `<svg class="`+classes+`" `, 1))
}

func MicroIcon(id, classes string) templ.Component {
	return templ.Raw(strings.Replace(heroiconsgo.Get("micro", id), "<svg", `<svg class="`+classes+`" `, 1))
}

// TODO: this is not ellegant, make it nicer, maybe make stroke width a variadic arg ? WithStroke ?
func Icon2(id, classes string) templ.Component {
	return templ.Raw(strings.Replace(strings.Replace(heroiconsgo.Get("outline", id), "<svg", `<svg class="`+classes+`" `, 1), `stroke-width="1.5"`, `stroke-width="2"`, 1))
}
