package icon

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/s00500/chimp/icon/base"
)

// Props defines the properties that can be set for an icon.
type props struct {
	Size        int
	Color       string
	Fill        string
	Stroke      string
	StrokeWidth string // Stroke Width of Icon, Usage: "2.5"
	Class       string
}

type iconProperty func(*props)

func Icon(icon base.IconBase, classes string, properties ...iconProperty) templ.Component {
	p := props{
		Size:        24,
		Fill:        "none",
		Stroke:      "currentColor",
		StrokeWidth: "1.5",
	}

	for _, prop := range properties {
		prop(&p)
	}

	// Hm... the defaults for fill and stroke are funkey on mdi...

	return templ.Raw(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" class=%q fill=%q stroke=%q stroke-width=%q viewBox="0 0 24 24">%s</svg>`, p.Size, p.Size, classes, p.Fill, p.Stroke, p.StrokeWidth, string(icon)))
}

func WithSize(size int) iconProperty {
	return func(p *props) {
		p.Size = size
	}
}
func WithFill(value string) iconProperty {
	return func(p *props) {
		p.Fill = value
	}
}
func WithStroke(value string) iconProperty {
	return func(p *props) {
		p.Stroke = value
	}
}
func WithStrokeWidth(value string) iconProperty {
	return func(p *props) {
		p.StrokeWidth = value
	}
}
