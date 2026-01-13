package chimp

import "github.com/a-h/templ"

// Base inserts a base tag, removing the need from many other places to use BaseURL
func Base(baseUrl string) templ.Component {
	return templ.Raw(`<base href=` + templ.URL(baseUrl) + `></div>`)
}
