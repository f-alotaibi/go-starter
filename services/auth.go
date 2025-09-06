package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/f-alotaibi/go-starter/models/types"
	"github.com/f-alotaibi/go-starter/repositories"
	"github.com/f-alotaibi/go-starter/utils"
	"github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/avatar"
	"github.com/go-pkgz/auth/v2/logger"
	"github.com/go-pkgz/auth/v2/provider"
	"github.com/go-pkgz/auth/v2/token"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

func NewAuth(db *gorm.DB) (*auth.Service, error) {
	if db == nil {
		return nil, fmt.Errorf("auth: could not find database")
	}

	authService := auth.NewService(auth.Opts{
		SecretReader: token.SecretFunc(func(aud string) (string, error) {
			return os.Getenv("AUTH_SECRET"), nil
		}),
		TokenDuration:  time.Minute * 5,              // JWT
		CookieDuration: time.Hour * 24,               // Session
		Issuer:         os.Getenv("AUTH_JWT_ISSUER"), // TODO: env
		URL:            "",                           // TODO: env
		AvatarStore:    avatar.NewNoOp(),
		DisableXSRF:    false, // Checks for env(AUTH_XSRF_NAME) check middlewares/auth.go
		JWTCookieName:  os.Getenv("AUTH_JWT_COOKIE_NAME"),
		XSRFCookieName: os.Getenv("AUTH_XSRF_NAME"),
		XSRFHeaderKey:  fmt.Sprintf("X-%s", os.Getenv("AUTH_XSRF_NAME")),
		Logger:         logger.Std,
		ClaimsUpd: token.ClaimsUpdFunc(func(claims token.Claims) token.Claims {
			user, err := repositories.FindUserByUsername(db, claims.User.Name)
			if err != nil {
				return token.Claims{}
			}
			claims.User.SetStrAttr("pwdReset", user.LastPasswordReset.Time.Format(time.RFC3339Nano))
			claims.User.SetRole(string(user.Role))
			claims.User.SetAdmin(user.Role == types.AdminRole)
			return claims
		}),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			user, err := repositories.FindUserByUsername(db, claims.User.Name)
			if err != nil {
				return false
			}
			log.Println(claims.User.Attributes)
			claimLastPasswordResetTime, err := time.Parse(time.RFC3339Nano, claims.User.StrAttr("pwdReset"))
			if err != nil {
				log.Println("PARSE", err)
				return false
			}
			if !claimLastPasswordResetTime.Equal(user.LastPasswordReset.Time) {
				log.Println(claimLastPasswordResetTime, "NOT EQUAL", user.LastPasswordReset.Time)
				log.Println("NOT EQUAL")
				return false
			}
			return true
		}),
	})

	authService.AddDirectProvider("users", provider.CredCheckerFunc(func(username, password string) (ok bool, err error) {
		b4 := time.Now()
		user, err := repositories.FindUserByUsername(db, username)
		if err != nil {
			return false, err
		}
		if utils.VerifyPassword(password, string(user.Password)) {
			log.Println("DATABASE QUERY: ", time.Now().Sub(b4))
			return true, nil
		}
		return false, nil
	}))

	return authService, nil
}
