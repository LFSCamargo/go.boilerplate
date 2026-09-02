package main

import (
	"go.boilerplate/src/config"
	"go.boilerplate/src/db"
	applog "go.boilerplate/src/log"
	"go.boilerplate/src/mail"
	"go.boilerplate/src/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.NewConfig()
	applog.Init(cfg.LogLevel, cfg.LogFormat)

	database := db.MustConnect(cfg.PostgresConnection)
	defer func() {
		if err := database.Close(); err != nil {
			applog.Error("database close failed", "err", err)
		}
	}()

	if err := database.Ping(); err != nil {
		applog.Fatal("database ping failed", "err", err)
	}

	if err := db.RunMigrations(cfg.PostgresConnection); err != nil {
		applog.Fatal("migrations failed", "err", err)
	}

	var mailer *mail.Service
	if mail.IsConfigured(cfg) {
		svc, err := mail.NewFromConfig(cfg)
		if err != nil {
			applog.Fatal("mail service init failed", "err", err)
		}
		mailer = svc
		applog.Info("mail service ready", "mode", "react-email+smtp")
	} else {
		applog.Warn("mail service disabled", "reason", "SMTP_HOST or MAIL_FROM not set")
	}

	app := fiber.New()
	routes.Setup(app, routes.Deps{
		Config: cfg,
		DB:     database,
		Mail:   mailer,
	})

	applog.Info("server listening", "port", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		applog.Fatal("server failed", "err", err)
	}
}
