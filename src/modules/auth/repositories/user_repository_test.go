package repositories_test

import (
	"context"
	"os"
	"testing"

	"go.boilerplate/src/db"
	"go.boilerplate/src/db/models"
	"go.boilerplate/src/modules/auth/repositories"

	"github.com/google/uuid"
)

func TestUserRepository_CreateAndFindByEmail(t *testing.T) {
	dsn := os.Getenv("POSTGRES_CONNECTION")
	if dsn == "" {
		dsn = "postgres://goboilerplate:goboilerplate@localhost:5432/goboilerplate?sslmode=disable"
	}
	database, err := db.Connect(dsn)
	if err != nil {
		t.Skipf("postgres unavailable: %v", err)
	}
	defer database.Close()
	if err := database.Ping(); err != nil {
		t.Skipf("postgres ping failed: %v", err)
	}
	if err := db.RunMigrations(dsn); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repo := repositories.NewUserRepository(database.DB, t.TempDir(), "http://localhost:8080")
	email := "repo-" + uuid.NewString() + "@example.com"
	user := &models.User{
		Email:        email,
		PasswordHash: "hash",
	}
	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatal(err)
	}

	found, err := repo.FindByEmail(context.Background(), email)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil || found.Email != email {
		t.Fatalf("expected user, got %+v", found)
	}
}
