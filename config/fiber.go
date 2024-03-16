package config

import (
	"test-mkp/exception"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewFiberConfig() fiber.Config {
	return fiber.Config{
		ErrorHandler: exception.ErrorHandler,
		BodyLimit:    31457280, // (around 30MB-31MB) value in bytes (1024B = 1KB, 1024KB = 1MB)
	}
}

func NewFiberLoggerConfig() logger.Config {
	return logger.Config{
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}
}
