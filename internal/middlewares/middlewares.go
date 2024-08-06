package middlewares

import (
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/gin-gonic/gin"
)

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

func stripRouteParams(c *gin.Context) string {
	path := c.Request.URL.Path

	// Remove parameters in path
	for _, p := range c.Params {
		path = strings.Replace(path, p.Value, "", -1)
	}
	// Remote trailing slash
	path = strings.TrimRight(path, "/")

	return path
}
