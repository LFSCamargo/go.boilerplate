package db_test

import (
	"testing"

	"go.boilerplate/src/db"
)

func TestConnect_InvalidDSN(t *testing.T) {
	_, err := db.Connect("postgres://invalid:invalid@127.0.0.1:1/nope?sslmode=disable")
	if err == nil {
		t.Fatal("expected error for invalid dsn")
	}
}

func TestRunMigrations_InvalidDSN(t *testing.T) {
	err := db.RunMigrations("postgres://invalid:invalid@127.0.0.1:1/nope?sslmode=disable")
	if err == nil {
		t.Fatal("expected migration error for invalid dsn")
	}
}
