package loggergin

import (
	"github.com/airbloc/logger"
	"github.com/gin-gonic/gin"
)

func Middleware(loggerName string) gin.HandlerFunc {
	log := logger.New(loggerName)

	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			url = url + "?" + c.Request.URL.RawQuery
		}

		timer := log.Timer()
		c.Next()

		statusCode := c.Writer.Status()
		statusColor := logger.Green
		if statusCode < 200 || statusCode >= 300 {
			statusColor = logger.Red
		}
		info := logger.Attrs{
			"method": c.Request.Method,
			"url":    url,
			"status": statusCode,
			"client": c.ClientIP(),
		}
		timer.End("{method} {url} – HTTP {}{status}{} – {client}",
			statusColor,
			logger.Reset,
			info)
	}
}
