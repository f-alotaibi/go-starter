package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/f-alotaibi/go-starter/controllers"
	"github.com/f-alotaibi/go-starter/middlewares"
	"github.com/f-alotaibi/go-starter/services"
	"github.com/go-pkgz/auth/v2/token"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"

	_ "net/http/pprof"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if os.Getenv("APP_ENV") == "dev" {
		log.Println("Enabling pprof for profiling")
		go func() {
			log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
		}()
	}
	zapLogger := zap.Must(zap.NewDevelopment())
	defer zapLogger.Sync()
	slog.SetDefault(slog.New(zapslog.NewHandler(zapLogger.Core())))

	db, err := services.NewDB()
	if err != nil {
		panic(err)
	}

	authService, err := services.NewAuth(db)
	if err != nil {
		panic(err)
	}
	authMiddleware := authService.Middleware()
	//authHandler, _ := authService.Handlers()

	e := echo.New()

	e.Use(middleware.Gzip())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieName:  os.Getenv("CSRF_COOKIE_NAME"),
		ContextKey:  os.Getenv("CSRF_KEY"),
		TokenLookup: fmt.Sprintf("form:%s", os.Getenv("CSRF_KEY")),
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
			}

			switch {
			case v.Status >= 500:
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR", attrs...)
				break
			case v.Status >= 400:
				slog.LogAttrs(context.Background(), slog.LevelWarn, "REQUEST_WARN", attrs...)
				break
			default:
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST", attrs...)
				break
			}
			return nil
		},
	}))

	e.Use(middlewares.InjectCSRFToContext())
	e.Use(middlewares.InjectFormErrorToContext())

	e.Static("/*", "assets/public")

	indexController := controllers.NewIndexController()
	e.GET("/", indexController.Show)
	loginController := controllers.NewLoginController(authService)
	e.GET("/login", loginController.Show)
	e.POST("/login", loginController.Post)
	signupController := controllers.NewSignupController(db, authService)
	e.GET("/signup", signupController.Show)
	e.POST("/signup", signupController.Post)

	// TODO: Enable /auth if you going to use any of go-pkgz/auth features
	//e.Any("/auth/*", echo.WrapHandler(authHandler))

	e.GET("/private", func(c echo.Context) error {
		user, err := token.GetUserInfo(c.Request())
		if err == nil {
			log.Printf("%v", user)
		}
		log.Println(err)
		return c.String(http.StatusOK, "authed")
	}, middlewares.AuthMiddlware(), echo.WrapMiddleware(authMiddleware.Auth))

	e.Logger.Fatal(e.Start(fmt.Sprintf("127.0.0.1:%s", os.Getenv("HTTP_PORT"))))
}
