package middlewares

import "github.com/bmatcuk/doublestar/v4"

var (
	observeRoutes = []string{
		"/api/**",
	}
)

func shouldObserve(route string) bool {
	for _, r := range observeRoutes {
		if ok, _ := doublestar.Match(r, route); ok {
			return true
		}
	}
	return false
}
