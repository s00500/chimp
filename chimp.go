package chimp

import (
	"bytes"
	"context"
	"net/http"

	"github.com/a-h/templ"
	log "github.com/s00500/env_logger"
)

func RenderString(c context.Context, component templ.Component) string {
	buf := bytes.NewBuffer([]byte{})
	log.Should(component.Render(c, buf))
	return buf.String()
}

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	component.Render(r.Context(), w)
}
