package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	cc "github.com/s00500/chimp/components"

	"github.com/s00500/chimp"
	"github.com/starfederation/datastar-go/datastar"
)

func main() {
	r := chi.NewRouter()
	r.Use(chimp.BaseMiddleware(""))

	r.Get("/", templ.Handler(page()).ServeHTTP)

	r.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)

		files, err := cc.ReadUploadedFiles(r, "file")
		if err != nil {
			cc.SendError(sse, "Upload failed: "+err.Error())
			return
		}
		if len(files) == 0 {
			cc.SendError(sse, "No files received")
			return
		}

		for _, f := range files {
			fmt.Printf("Received: %s (%d bytes)\n", f.Name, f.Size)
		}
		cc.SendSuccess(sse, fmt.Sprintf("Uploaded %d file(s)", len(files)))
	})

	chimp.ServeIncludedAssets(r, "")
	chimp.ServeHotReload(r, true)

	fmt.Println("Dropzone demo at http://localhost:8090")
	http.ListenAndServe(":8090", r)
}
