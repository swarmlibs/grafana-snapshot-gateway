package middlewares

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swarmlibs/grafana-snapshot-gateway/internal/metrics"
)

func MeasureResponseDuration(mc metrics.MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldObserve(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		duration := time.Since(start)
		route := c.Request.URL.Path

		// Remove parameters in path
		for _, p := range c.Params {
			route = strings.Replace(route, p.Value, "", -1)
		}
		// Remote trailing slash
		route = strings.TrimRight(route, "/")

		method := c.Request.Method
		statusCode := strconv.Itoa(c.Writer.Status())
		mc.RequestDurationSecondsHistogram.WithLabelValues(route, method, statusCode).Observe(duration.Seconds())
	}
}
