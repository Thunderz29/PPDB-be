package middleware

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		var extra string
		if msg, exists := c.Get("res_msg"); exists {
			extra += fmt.Sprintf("[%s]", msg)
		}
		if errStr, exists := c.Get("res_err"); exists {
			cleanErr := strings.ReplaceAll(fmt.Sprintf("%v", errStr), "\n", " ")
			if extra != "" {
				extra += " "
			}
			extra += "Error: " + cleanErr
		}

		msg := fmt.Sprintf("%3d | %12v | %15s | %-7s %s", status, latency, clientIP, method, path)
		if query != "" {
			msg += "?" + query
		}

		if extra != "" {
			msg += " | " + extra
		}

		// Don't log health check spam
		if path == "/api/health" {
			return
		}

		if status >= 500 {
			slog.Error(msg)
		} else if status >= 400 {
			slog.Warn(msg)
		} else {
			slog.Info(msg)
		}
	}
}
