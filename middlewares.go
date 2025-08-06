package chimp

import (
	"context"
	"net/http"
	"strings"
)

// Uses the refer header of th erequest to determin the current page in the templates
func UrlPathMiddleware(baseUrl string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.SplitAfterN(r.Referer(), "/", 2)
			if len(parts) > 1 {
				ctx := context.WithValue(r.Context(), "urlpath", strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, parts[1]), baseUrl))
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				//ctx := context.WithValue(r.Context(), "urlpath", strings.TrimPrefix(r.URL.Path, baseURL))
				next.ServeHTTP(w, r)
			}
		})
	}
}
