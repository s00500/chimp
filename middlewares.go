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
			if isSSE(r) {

				_, rest := splitAfterThirdSlash(r.Referer())
				ctx := context.WithValue(r.Context(), "urlpath", strings.TrimPrefix("/"+rest, strings.TrimSuffix(baseUrl, "/")))
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				ctx := context.WithValue(r.Context(), "urlpath", strings.TrimPrefix(r.URL.Path, baseUrl))
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

// Function to use in your templates for get URL
func URL(ctx context.Context) string {
	s, ok := ctx.Value("urlpath").(string)
	if !ok {
		//log.Debug("no url")
		return ""
	}
	return s
}

func isSSE(r *http.Request) bool {
	return r.Header.Get("Datastar-Request") == "true"
}

func splitAfterThirdSlash(s string) (string, string) {
	slashCount := 0

	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			slashCount++
			if slashCount == 3 {
				// Return the split parts
				return s[:i+1], s[i+1:]
			}
		}
	}
	return s, ""
}
