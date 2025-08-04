package chimp

import (
	"bytes"
	"net/http"
	"time"

	_ "embed"

	"github.com/a-h/templ"
)

// I need to setup static hosting and THEN I can add template functions for them...

var modTime = time.Date(2025, 8, 4, 0, 0, 0, 0, time.UTC)

//go:embed static/datastar.min.js
var datastar []byte

func DatastarTemplate(baseUrl string) templ.Component {
	return templ.Raw(`<script type="module" src="` + baseUrl + `/static/datastar.min.js"></script>`)
}

func datastarHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewReader(datastar)
	http.ServeContent(w, r, "datastar.min.js", modTime, buf)
}

//go:embed static/basecoat.cdn.min.css
var basecoatCss []byte

//go:embed static/all.min.js
var basecoatJs []byte

func BaseCoatTemplate(baseUrl string) templ.Component {
	return templ.Raw(`<script type="module" src="` + baseUrl + `/static/datastar.min.js"></script>
  <link rel="stylesheet" href="` + baseUrl + `/static/basecoat.min.css"/>`)
}

func baseCoatCSSHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewReader(basecoatCss)
	http.ServeContent(w, r, "basecoat.min.css", modTime, buf)
}

func baseCoatJSHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewReader(basecoatJs)
	http.ServeContent(w, r, "basecoat.min.js", modTime, buf)
}

func StaticMiddleware(baseURL string) func(next http.Handler) http.Handler {
	datastarURL := baseURL + "/static/datastar.min.js"
	baseCoatCSSURL := baseURL + "/static/basecoat.min.css"
	baseCoatJSURL := baseURL + "/static/basecoat.min.js"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case datastarURL:
				datastarHandler(w, r)
			case baseCoatCSSURL:
				baseCoatCSSHandler(w, r)
			case baseCoatJSURL:
				baseCoatJSHandler(w, r)
			default:
				next.ServeHTTP(w, r)
			}
		})
	}
}
