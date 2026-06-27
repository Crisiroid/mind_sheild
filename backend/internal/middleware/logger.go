package middleware

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		log.Printf(
			"[%s] %s %s - %d - %v",
			c.Request().Method,
			c.Request().RequestURI,
			c.RealIP(),
			c.Response().Status,
			time.Since(start),
		)

		return err
	}
}
