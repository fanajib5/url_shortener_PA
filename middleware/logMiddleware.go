package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func LogMiddleware(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} method=${method}, " +
			"uri=${uri}, status=${status}, latency_human=${latency_human}\n",
		// CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
}
