package routes

import (
	"path/filepath"

	"go.boilerplate/src/config"
	"go.boilerplate/src/db"
	applog "go.boilerplate/src/log"
	"go.boilerplate/src/mail"
	"go.boilerplate/src/middleware"
	authRouter "go.boilerplate/src/modules/auth"
	authHandlers "go.boilerplate/src/modules/auth/handlers"
	"go.boilerplate/src/modules/auth/repositories"
	"go.boilerplate/src/modules/auth/services"
	healthRouter "go.boilerplate/src/modules/health"
	"go.boilerplate/src/openapi"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

type Deps struct {
	Config *config.Config
	DB     *db.DB
	Mail   *mail.Service
}

func Setup(app *fiber.App, deps Deps) {
	app.Use(applog.HTTP())
	if deps.Config != nil {
		app.Use(middleware.AuthRateLimit(deps.Config))
	}

	api := openapi.Attach(app)
	healthRouter.Register(api)

	if deps.Config != nil {
		staticDir := filepath.Join(deps.Config.EmailsDir, "static")
		app.Use("/static", static.New(staticDir))
		app.Use("/uploads", static.New(deps.Config.UploadsDir))
	}

	if deps.DB == nil || deps.Config == nil {
		return
	}

	var mailer services.Mailer = services.NoopMailer{}
	if deps.Mail != nil {
		mailer = deps.Mail
	}

	uploads := "uploads"
	publicURL := "http://localhost:8080"
	if deps.Config.UploadsDir != "" {
		uploads = deps.Config.UploadsDir
	}
	if deps.Config.AppPublicURL != "" {
		publicURL = deps.Config.AppPublicURL
	}

	authSvc := services.NewAuthService(
		repositories.NewUserRepository(deps.DB.DB, uploads, publicURL),
		repositories.NewOTPRepository(deps.DB.DB),
		repositories.NewRevokedTokenRepository(deps.DB.DB),
		mailer,
		services.Config{
			BearerSecret: deps.Config.BearerSecret,
			AppPublicURL: deps.Config.AppPublicURL,
		},
	)

	authRouter.Register(api, authHandlers.New(authSvc), authSvc)
}
