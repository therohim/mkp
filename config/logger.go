package config

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoggerConfig struct {
	// For skip middleware
	Next func(ctx *fiber.Ctx) bool

	// Default: Using Zap as logger library
	Logger *zap.Logger
}

func NewLogger(cfg ...LoggerConfig) fiber.Handler {
	var conf LoggerConfig
	if len(cfg) > 0 {
		conf = cfg[0]
	}

	return func(c *fiber.Ctx) error {
		// Don't execute the middleware if Next returns true
		if conf.Next != nil && conf.Next(c) {
			c.Next()
			return nil
		}

		loc, _ := time.LoadLocation(os.Getenv("APP_TIMEZONE"))
		start := time.Now().In(loc)
		headers := c.GetReqHeaders()

		// handle request
		c.Next()

		logTime := start.Format("02-Jan-2006 15:04:05")
		method := c.Method()
		path := c.Path()
		duration := time.Since(start).Milliseconds()
		msg := map[string]interface{}{
			"actor": headers["X-User"],
			"headers": headers,
			"time": logTime,
			"method": method,
			"path": path,
			"duration": fmt.Sprintf("%dms", duration),
		}

		code := c.Response().StatusCode()

		switch {
		case (code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError) || duration > 2000:
			conf.Logger.Warn("REQUEST_WARNING", zap.Any("data", msg))
		case code >= http.StatusInternalServerError:
			conf.Logger.Error("REQUEST_ERROR", zap.Any("data", msg))
		}

		return nil
	}
}
