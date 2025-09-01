package middlewares

import (
	"context"
	"os"

	"github.com/labstack/echo/v4"
)

// Middleware that passes the CSRF token (from echo/middleware.CRSF()) and passes it to request's context value
func CSRFToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), "csrf", c.Get(os.Getenv("CSRF_KEY")))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
