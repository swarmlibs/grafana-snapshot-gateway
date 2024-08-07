package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swarmlibs/grafana-snapshot-gateway/internal/metrics"
)

func MeasureRequestDuration(mc metrics.MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldObserve(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Measure the duration of the request
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		method := c.Request.Method
		path := stripRouteParams(c)
		statusCode := strconv.Itoa(c.Writer.Status())
		mc.RequestDurationSecondsHistogram.WithLabelValues(path, method, statusCode).Observe(duration.Seconds())
	}
}
