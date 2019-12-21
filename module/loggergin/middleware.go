package loggergin

import (
	"fmt"
	"github.com/airbloc/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Middleware(loggerName string) gin.HandlerFunc {
	log := logger.New(loggerName)

	return func(c *gin.Context) {
		timer := log.Timer()
		c.Next()

		statusCode := c.Writer.Status()
		statusColor := logger.Green
		if statusCode < 200 || statusCode >= 300 {
			statusColor = logger.Red
		}
		info := logger.Attrs{
			"method": c.Request.Method,
			"url":    getRequestPath(c.Request),
			"status": statusCode,
			"client": c.ClientIP(),
		}
		timer.End("{method} {url} – {}HTTP {status}{} – {client}",
			statusColor,
			logger.Reset,
			info)
	}
}

func Recovery(loggerName string) gin.HandlerFunc {
	log := logger.New(loggerName)

	return func(c *gin.Context) {
		defer func() {
			info := logger.Attrs{
				"method": c.Request.Method,
				"url":    getRequestPath(c.Request),
				"client": c.ClientIP(),
			}
			if r := log.Recover(info); r != nil {
				c.Error(fmt.Errorf("panic: %v", r))
			}
		}()
		c.Next()
	}
}

func getRequestPath(r *http.Request) string {
	path := r.URL.Path
	raw := r.URL.RawQuery
	if raw != "" {
		return path + "?" + raw
	}
	return path
}
