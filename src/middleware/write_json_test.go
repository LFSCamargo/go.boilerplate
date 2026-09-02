package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	z "github.com/Oudwins/zog"
	"github.com/gofiber/fiber/v3"
)

func TestWriteJSON_WritesValidPayload(t *testing.T) {
	schema := z.Struct(z.Shape{"Status": z.String().Required().OneOf([]string{"ok"})})
	type out struct {
		Status string `json:"status"`
	}

	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return WriteJSON(c, schema, out{Status: "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Fatalf("body %s", body)
	}
}
