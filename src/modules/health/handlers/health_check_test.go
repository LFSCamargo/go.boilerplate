package HealthHandlers

import (
	"context"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	out, err := HealthCheckHandler(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if out == nil || out.Body.Status != "ok" {
		t.Fatalf("got %+v", out)
	}
}
