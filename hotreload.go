package chimp

import (
	"net/http"
	"sync"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

// ServeHotReload adds a datastar based hotreload handler. Use the enable flag to easily disable in production
func ServeHotReload(r *chi.Mux, baseURL string, enable bool) {
	if !enable {
		return
	}
	var hotReloadOnlyOnce sync.Once
	r.Get(baseURL+"/hotreload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		hotReloadOnlyOnce.Do(func() {
			sse.ExecuteScript("window.location.reload()")
		})
		<-r.Context().Done()
	})
}

// HotReload adds a datastar based hotreload element to the frontend. Use the enable flag to easily disable in production
func HotReload(baseUrl string, enable bool) templ.Component {
	if !enable {
		return templ.Raw(``)
	}
	return templ.Raw(`<div id="hotreload" data-init="@get('` + baseUrl + `/hotreload', {retryMaxCount: 1000,retryInterval:20, retryMaxWaitMs:200})"></div>`)
}
