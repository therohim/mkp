package config

import "github.com/gofiber/fiber/v2/middleware/cors"

func NewCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins: "http://localhost:4000, http://localhost:3300, http://localhost:3100, https://app.vcgamers.io, https://auth.vcgamers.io, https://arena.vcgamers.io, https://hub.vcgamers.io, https://app.vcgamers.com, https://auth.vcgamers.com, https://arena.vcgamers.com, https://hub.vcgamers.com",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}
}
