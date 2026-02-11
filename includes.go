package chimp

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	_ "embed"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

// I need to setup static hosting and THEN I can add template functions for them...

var modTime = time.Date(2025, 8, 4, 0, 0, 0, 0, time.UTC)

//go:embed static/datastar.min.js
var datastarBytes []byte

func IncludedDatastar() templ.Component {
	return templ.Raw(`<script type="module" src="static/datastar.min.js?v=rc6"></script>`)
}

func datastarHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewReader(datastarBytes)
	http.ServeContent(w, r, "datastar.min.js", modTime, buf)
}

//go:embed static/basecoat.cdn.min.css
var basecoatCss []byte

//go:embed static/all.min.js
var basecoatJs []byte

func IncludedBaseCoatCSS() templ.Component {
	return templ.Raw(`<link rel="stylesheet" href="static/basecoat.min.css"/>`)
}

func IncludedBaseCoatJS() templ.Component {
	return templ.Raw(`<script type="module" src="static/basecoat.min.js"></script>`)
}

func baseCoatCSSHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewReader(basecoatCss)
	http.ServeContent(w, r, "basecoat.min.css", modTime, buf)
}

func baseCoatJSHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewReader(basecoatJs)
	http.ServeContent(w, r, "basecoat.min.js", modTime, buf)
}

//go:embed static/neo-sans-std-m.otf
var neoSansFontBytes []byte

func neoSansFontHandler(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewReader(neoSansFontBytes)
	http.ServeContent(w, r, "neo-sans-std-m.otf", modTime, buf)
}

func IncludedFont() templ.Component {
	return templ.Raw(`<style>
@font-face {
    font-family: 'Neo Sans';
    src: url('static/neo-sans-std-m.otf') format('opentype');
    font-weight: 500;
    font-style: normal;
    font-display: swap;
}
body {
    font-family: 'Neo Sans', sans-serif;
}
</style>`)
}

// ServeIncludedAssets serves included datastar, basecoat, and chimp assets
func ServeIncludedAssets(r chi.Router, baseURL string) {
	p := strings.Trim(baseURL, "/")
	if p != "/" && p != "" {
		p = "/" + p
	}
	r.Get(p+"/static/datastar.min.js", datastarHandler)
	r.Get(p+"/static/basecoat.min.css", baseCoatCSSHandler)
	r.Get(p+"/static/basecoat.min.js", baseCoatJSHandler)
	r.Get(p+"/static/neo-sans-std-m.otf", neoSansFontHandler)
}
