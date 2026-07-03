package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ZerologLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		dummyLogger := log.Info()
		if len(c.Errors) > 0 {
			dummyLogger = log.Error().Err(c.Errors.Last())
		} else if c.Writer.Status() >= 400 && c.Writer.Status() < 500 {
			dummyLogger = log.Warn()
		} else if c.Writer.Status() >= 500 {
			dummyLogger = log.Error()
		}

		dummyLogger.
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Msg(c.Request.Method + " " + path)
	}
}
