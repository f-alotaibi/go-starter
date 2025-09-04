package controllers

import (
	"net/http"

	"github.com/f-alotaibi/go-starter/utils"
	"github.com/f-alotaibi/go-starter/views"
	"github.com/go-pkgz/auth/v2"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	User     string `form:"user" validate:"required"`
	Password string `form:"passwd" validate:"required"`
}

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
	req := new(LoginRequest)
	if err := ctx.Bind(req); err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.Login().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	if errors, ok := utils.ValidateStruct(req); !ok {
		return views.Login().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	// validate creds using our Direct provider
	provider, err := c.authService.Provider("users")
	if err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.Login().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	crw := &captureWriter{
		headers: make(http.Header),
	}
	provider.LoginHandler(crw, ctx.Request())

	if crw.status != 0 {
		errors := map[string]string{
			"result": "Invalid credentials",
		}
		return views.Login().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
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
