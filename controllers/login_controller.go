package controllers

import (
	"net/http"

	"github.com/f-alotaibi/go-starter/views"
	"github.com/go-pkgz/auth/v2"
	"github.com/labstack/echo/v4"
)

type LoginController struct {
	authService *auth.Service
}

func NewLoginController(authService *auth.Service) *LoginController {
	return &LoginController{
		authService: authService,
	}
}

func (c *LoginController) Show(ctx echo.Context) error {
	return views.Login().Render(ctx.Request().Context(), ctx.Response())
}

func (c *LoginController) Post(ctx echo.Context) error {
	// validate creds using our Direct provider
	provider, err := c.authService.Provider("users")
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "invalid credentials")
	}

	crw := &captureWriter{
		headers: make(http.Header),
	}
	provider.LoginHandler(crw, ctx.Request())

	if crw.status != 0 {
		return ctx.String(crw.status, "Could not capture")
	}

	for _, cookie := range crw.headers["Set-Cookie"] {
		ctx.Response().Header().Add("Set-Cookie", cookie)
	}

	// redirect after login
	return ctx.Redirect(http.StatusFound, "/")
}

type captureWriter struct {
	headers http.Header
	status  int
}

func (cw *captureWriter) Header() http.Header {
	return cw.headers
}

func (cw *captureWriter) WriteHeader(status int) {
	cw.status = status
}

func (cw *captureWriter) Write(b []byte) (int, error) {
	return len(b), nil // ignore body
}
