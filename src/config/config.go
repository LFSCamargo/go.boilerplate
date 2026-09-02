package config

import (
	"fmt"
	"os"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/zenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string `env:"PORT"`
	PostgresConnection string `env:"POSTGRES_CONNECTION"`
	BearerSecret       string `env:"BEARER_SECRET"`
	CORSOrigins        string `env:"CORS_ORIGINS"`
	SMTPHost           string `env:"SMTP_HOST"`
	SMTPPort           string `env:"SMTP_PORT"`
	SMTPSecurity       string `env:"SMTP_SECURITY"`
	SMTPUser           string `env:"SMTP_USER"`
	SMTPPassword       string `env:"SMTP_PASSWORD"`
	MailFrom           string `env:"MAIL_FROM"`
	EmailsDir          string `env:"EMAILS_DIR"`
	EmailAssetsBaseURL string `env:"EMAIL_ASSETS_BASE_URL"`
	AppPublicURL       string `env:"APP_PUBLIC_URL"`
	UploadsDir             string `env:"UPLOADS_DIR"`
	LogLevel               string `env:"LOG_LEVEL"`
	LogFormat              string `env:"LOG_FORMAT"`
	RateLimitMax           int    `env:"RATE_LIMIT_MAX"`
	RateLimitWindowMinutes int    `env:"RATE_LIMIT_WINDOW_MINUTES"`
}

var configSchema = z.Struct(z.Shape{
	"Port":               z.String().Default("5000"),
	"PostgresConnection": z.String().Required().Min(1),
	"BearerSecret":       z.String().Required().Min(32),
	"CORSOrigins":        z.String().Required().Min(1),
	"SMTPHost":           z.String().Optional(),
	"SMTPPort":           z.String().Optional(),
	"SMTPSecurity":       z.String().Optional(),
	"SMTPUser":           z.String().Optional(),
	"SMTPPassword":       z.String().Optional(),
	"MailFrom":           z.String().Optional(),
	"EmailsDir":          z.String().Default("emails"),
	"EmailAssetsBaseURL": z.String().Default("http://localhost:8080"),
	"AppPublicURL":       z.String().Default("http://localhost:8080"),
	"UploadsDir":         z.String().Default("uploads"),
	"LogLevel":               z.String().Default("info"),
	"LogFormat":              z.String().Default("pretty"),
	"RateLimitMax":           z.Int().Default(10),
	"RateLimitWindowMinutes": z.Int().Default(1),
})

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "No .env file found, using system environment variables")
	}

	var cfg Config
	if errs := configSchema.Parse(zenv.NewDataProvider(), &cfg); errs != nil {
		fmt.Fprintf(os.Stderr, "invalid configuration: %v\n", z.Issues.Flatten(errs))
		os.Exit(1)
	}

	return &cfg
}
