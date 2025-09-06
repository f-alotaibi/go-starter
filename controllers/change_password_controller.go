package controllers

import (
	"net/http"
	"time"

	"github.com/f-alotaibi/go-starter/repositories"
	"github.com/f-alotaibi/go-starter/utils"
	"github.com/f-alotaibi/go-starter/views"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ChangePasswordRequest struct {
	Password             string `form:"passwd" validate:"required,min=8,password"`
	PasswordConfirmation string `form:"passwd_confirm" validate:"required,eqfield=Password"`
}

type ChangePasswordController struct {
	db *gorm.DB
}

func NewChangePasswordController(database *gorm.DB) *ChangePasswordController {
	return &ChangePasswordController{
		db: database,
	}
}

func (c *ChangePasswordController) Show(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	encodedToken := utils.GetEncodedHashedToken(token)
	resetToken, err := repositories.FindPasswordResetToken(c.db, encodedToken)
	if err != nil || time.Now().After(resetToken.Expiration) || resetToken.Used {
		return ctx.String(http.StatusNotFound, "invalid token")
	}

	return views.ChangePasswordForm().Render(ctx.Request().Context(), ctx.Response())
}

func (c *ChangePasswordController) Post(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	encodedToken := utils.GetEncodedHashedToken(token)
	resetToken, err := repositories.FindPasswordResetToken(c.db, encodedToken)
	if err != nil || time.Now().After(resetToken.Expiration) || resetToken.Used {
		return ctx.String(http.StatusNotFound, "invalid token")
	}

	req := new(ChangePasswordRequest)
	if err := ctx.Bind(req); err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.ChangePasswordForm().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	if errors, ok := utils.ValidateStruct(req); !ok {
		return views.ChangePasswordForm().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	err = repositories.SetPasswordResetTokenAsUsed(c.db, encodedToken)
	if err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.ChangePasswordForm().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.ChangePasswordForm().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	err = repositories.UpdateUserPassword(c.db, resetToken.UserID, hash)
	if err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.ChangePasswordForm().Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	return views.ChangePasswordSuccess().Render(ctx.Request().Context(), ctx.Response())
}
