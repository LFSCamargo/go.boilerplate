package main

import (
	"flag"
	"fmt"
	"os"

	"go.boilerplate/src/config"
	"go.boilerplate/src/db"
	applog "go.boilerplate/src/log"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up, down, version")
	flag.Parse()

	cfg := config.NewConfig()
	applog.Init(cfg.LogLevel, cfg.LogFormat)

	switch *direction {
	case "up":
		if err := db.RunMigrations(cfg.PostgresConnection); err != nil {
			applog.Fatal("migrate up failed", "err", err)
		}
		fmt.Println("migrations applied")
	case "down":
		if err := db.MigrateDown(cfg.PostgresConnection); err != nil {
			applog.Fatal("migrate down failed", "err", err)
		}
		fmt.Println("rolled back one migration")
	case "version":
		version, dirty, err := db.MigrationVersion(cfg.PostgresConnection)
		if err != nil {
			applog.Fatal("migration version failed", "err", err)
		}
		fmt.Printf("version=%d dirty=%v\n", version, dirty)
	default:
		fmt.Fprintf(os.Stderr, "unknown direction %q\n", *direction)
		os.Exit(1)
	}
}
