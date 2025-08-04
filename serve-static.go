package chimp

import (
	"embed"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ServeStatic is a helper to setup serving a static directory
func ServeStatic(r chi.Router, dir, prefix string) {
	p := strings.Trim(prefix, "/")
	r.Handle("/"+p+"/*", http.StripPrefix("/"+p+"/", http.FileServer(http.Dir(dir))))
}

// ServeStaticFromEmbedded is a helper to setup serving a static embedded fs
func ServeStaticFromEmbedded(r chi.Router, fs embed.FS, prefix string) {
	// use this then in the main UI
	// go:embed assets
	//var assetsFS embed.FS
	p := strings.Trim(prefix, "/")
	r.Handle("/"+p+"/*", http.StripPrefix("/", http.FileServer(http.FS(fs))))
}
