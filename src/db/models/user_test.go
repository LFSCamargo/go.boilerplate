package models

import "testing"

func TestUserTableName(t *testing.T) {
	if got := (User{}).TableName(); got != "users" {
		t.Fatalf("got %s", got)
	}
}
