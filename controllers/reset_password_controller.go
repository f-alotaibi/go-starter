package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/f-alotaibi/go-starter/models"
	"github.com/f-alotaibi/go-starter/repositories"
	"github.com/f-alotaibi/go-starter/utils"
	"github.com/f-alotaibi/go-starter/views"
	"github.com/labstack/echo/v4"
	"github.com/wneessen/go-mail"
	"gorm.io/gorm"
)

type ResetPasswordRequest struct {
	Email string `form:"email" validate:"required,email"`
}

type ResetPasswordController struct {
	db         *gorm.DB
	mailClient *mail.Client
}

func NewResetPasswordController(database *gorm.DB, mailClient *mail.Client) *ResetPasswordController {
	return &ResetPasswordController{
		db:         database,
		mailClient: mailClient,
	}
}

func (c *ResetPasswordController) Show(ctx echo.Context) error {
	return views.ResetPassword(false).Render(ctx.Request().Context(), ctx.Response())
}

func (c *ResetPasswordController) Post(ctx echo.Context) error {
	req := new(ResetPasswordRequest)
	if err := ctx.Bind(req); err != nil {
		errors := map[string]string{
			"result": "bad request",
		}
		return views.ResetPassword(false).Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	if errors, ok := utils.ValidateStruct(req); !ok {
		return views.ResetPassword(false).Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
	}

	user, err := repositories.FindUserByEmail(c.db, req.Email)
	if err == nil {
		// Send an email
		rawToken, hashedToken, err := utils.GenerateSecureToken()
		if err != nil {
			errors := map[string]string{
				"result": "bad request",
			}
			return views.ResetPassword(false).Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
		}

		resetToken := &models.PasswordResetToken{
			UserID:     user.ID,
			Token:      hashedToken,
			Expiration: time.Now().Add(time.Minute * 5),
		}

		err = repositories.CreatePasswordResetToken(c.db, resetToken)
		if err == gorm.ErrDuplicatedKey {
			// Retry 3 times if not just exit
			for i := 0; i < 3; i++ {
				rawToken, hashedToken, err = utils.GenerateSecureToken()
				if err != nil {
					errors := map[string]string{
						"result": "bad request",
					}
					return views.ResetPassword(false).Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
				}

				resetToken = &models.PasswordResetToken{
					UserID:     user.ID,
					Token:      hashedToken,
					Expiration: time.Now().Add(time.Minute * 5),
				}
				err = repositories.CreatePasswordResetToken(c.db, resetToken)
				if err == gorm.ErrDuplicatedKey {
					continue
				} else {
					break
				}
			}
		}

		if err != nil {
			errors := map[string]string{
				"result": "bad request",
			}
			return views.ResetPassword(false).Render(utils.WithErrors(ctx.Request().Context(), errors), ctx.Response())
		}

		// TODO: Move to services/
		go func(mailClient *mail.Client, email string, token string) {
			message := mail.NewMsg()
			if err := message.From(os.Getenv("MAIL_RESET_PASSWORD_SENDER")); err != nil {
				log.Fatalf("failed to set From address: %s", err)
			}
			if err := message.To(email); err != nil {
				log.Fatalf("failed to set To address: %s", err)
			}
			message.Subject("Password reset")
			message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf("Expires in 5 minutes: localhost:8080/change_password?token=%s", token))
			if err := mailClient.DialAndSend(message); err != nil {
				log.Fatalf("failed to send mail: %s", err)
			}
		}(c.mailClient, user.Email, rawToken)

		log.Println(rawToken, hashedToken)
	}

	return views.ResetPassword(true).Render(context.WithValue(ctx.Request().Context(), "pwdResetEmail", req.Email), ctx.Response())
}
