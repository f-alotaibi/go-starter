// For @error in views/components/form/error.templ

package middlewares

import (
	"context"

	"github.com/labstack/echo/v4"
)

func InjectFormErrorToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), "errors", map[string]string{})
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
