package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Bypass for go-pkgz/auth DisableXSRF
func AuthMiddlware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			xsrf, err := c.Cookie("XSRF-TOKEN")
			if err != nil {
				return c.String(http.StatusUnauthorized, "unauthorized")
			}
			c.Request().Header.Set("X-XSRF-TOKEN", xsrf.Value)
			return next(c)
		}
	}
}
