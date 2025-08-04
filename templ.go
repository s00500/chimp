package chimp

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// Several Templ Helpers

func StyleWithCacheBust(url string, cacheBust bool) templ.Component {
	if cacheBust {
		return templ.Raw(`<link rel="stylesheet" href="` + fmt.Sprintf("%s?v=%d", templ.EscapeString(url), rand.Int()) + `"/>`)
	} else {
		return templ.Raw(`<link rel="stylesheet" href="` + templ.EscapeString(url) + `"/>`)
	}
}
