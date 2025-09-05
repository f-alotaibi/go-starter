package controllers

import (
	"net/http"

	"github.com/f-alotaibi/go-starter/models"
	"github.com/f-alotaibi/go-starter/repositories"
	"github.com/f-alotaibi/go-starter/utils"
	"github.com/f-alotaibi/go-starter/views"
	"github.com/go-pkgz/auth/v2"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SignupRequest struct {
	User                 string `form:"user" validate:"required,min=4"`
	Email                string `form:"email" validate:"required,email"`
	Password             string `form:"passwd" validate:"required,min=8,password"`
	PasswordConfirmation string `form:"passwd_confirm" validate:"required,eqfield=Password"`
}

type SignupController struct {
	authService *auth.Service
	db          *gorm.DB
}

func NewSignupController(database *gorm.DB, authService *auth.Service) *SignupController {
	return &SignupController{
		authService: authService,
		db:          database,
	}
}

func (c *SignupController) Show(ctx echo.Context) error {
	return views.Signup().Render(ctx.Request().Context(), ctx.Response())
}

func (c *SignupController) Post(ctx echo.Context) error {
	req := new(SignupRequest)
	if err := ctx.Bind(req); err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.Signup().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	if errors, ok := utils.ValidateStruct(req); !ok {
		return views.Signup().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	signedUser, err := repositories.FindUser(c.db, req.User, req.Email)
	if err == nil {
		errors := map[string]string{}
		if signedUser.Username == req.User {
			errors["user"] = "Username already exists"
		}
		if signedUser.Email == req.Email {
			errors["email"] = "Email already exists"
		}
		return views.Signup().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.Signup().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	user := &models.User{
		Username: req.User,
		Email:    req.Email,
		Password: hash,
	}

	err = repositories.CreateUser(c.db, user)
	if err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.Signup().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	// validate creds using our Direct provider
	provider, err := c.authService.Provider("users")
	if err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.Signup().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	crw := &captureWriter{
		headers: make(http.Header),
	}
	provider.LoginHandler(crw, ctx.Request())

	if crw.status != 0 {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.Signup().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	for _, cookie := range crw.headers["Set-Cookie"] {
		ctx.Response().Header().Add("Set-Cookie", cookie)
	}

	// redirect after signup
	return ctx.Redirect(http.StatusFound, "/")
}
