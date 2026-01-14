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

func JsWithCacheBust(url string, cacheBust bool) templ.Component {
	if cacheBust {
		return templ.Raw(`<script src="` + fmt.Sprintf("%s?v=%d", templ.EscapeString(url), rand.Int()) + `"></script>`)
	} else {
		return templ.Raw(`<script src="` + templ.EscapeString(url) + `"></script>`)
	}
}

func JsModuleWithCacheBust(url string, cacheBust bool) templ.Component {
	if cacheBust {
		return templ.Raw(`<script type="module" src="` + fmt.Sprintf("%s?v=%d", templ.EscapeString(url), rand.Int()) + `"></script>`)
	} else {
		return templ.Raw(`<script type="module" src="` + templ.EscapeString(url) + `"></script>`)
	}
}
