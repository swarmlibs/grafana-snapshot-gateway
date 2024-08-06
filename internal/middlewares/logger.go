package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// StructuredLogger logs a gin HTTP request in JSON format. Allows to set the
// logger for testing purposes.
func StructuredLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldObserve(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now() // Start timer
		method := c.Request.Method
		path := stripRouteParams(c)

		// Call the next middleware
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.Path = path
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()

		// Log using the params
		var logEvent func(logger log.Logger) log.Logger
		if c.Writer.Status() >= 500 {
			logEvent = level.Error
		} else {
			logEvent = level.Info
		}

		msg := fmt.Sprintf("%s %s", method, c.Request.URL.Path)
		logEvent(*logger).Log("msg", msg, "method", param.Method, "path", param.Path, "status_code", param.StatusCode, "body_size", param.BodySize, "client_id", param.ClientIP, "latency", param.Latency.String(), "error", param.ErrorMessage)
	}
}
