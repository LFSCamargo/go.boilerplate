package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	z "github.com/Oudwins/zog"
	"github.com/gofiber/fiber/v3"
)

type loginPayload struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

var loginSchema = z.Struct(z.Shape{
	"Email":    z.String().Required().Email(),
	"Password": z.String().Required().Min(1),
})

func TestValidateBody_InjectsParsedJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/login", ValidateBody[loginPayload](loginSchema), func(c fiber.Ctx) error {
		body := Payload[loginPayload](c)
		return c.JSON(fiber.Map{"email": body.Email})
	})

	raw, _ := json.Marshal(map[string]string{"email": "ada@example.com", "password": "secret"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status %d %s", resp.StatusCode, body)
	}

	var out struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.Email != "ada@example.com" {
		t.Fatalf("got %+v", out)
	}
}

func TestValidateBody_RejectsInvalidJSON(t *testing.T) {
	app := fiber.New()
	app.Post("/login", ValidateBody[loginPayload](loginSchema), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":"nope"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestValidateBody_AllowsEmptyBody(t *testing.T) {
	app := fiber.New()
	app.Get("/health", ValidateBody[Empty](EmptySchema), func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestWriteJSON_RejectsInvalidResponse(t *testing.T) {
	schema := z.Struct(z.Shape{"Email": z.String().Required().Email()})
	type out struct {
		Email string `json:"email"`
	}

	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return WriteJSON(c, schema, out{Email: "nope"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestValidateQuery_InjectsParsedQuery(t *testing.T) {
	type tokenPayload struct {
		Token string `json:"token" query:"token"`
	}
	schema := z.Struct(z.Shape{
		"Token": z.String().Required().Min(1),
	})

	app := fiber.New()
	app.Get("/verify", ValidateQuery[tokenPayload](schema), func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"token": Payload[tokenPayload](c).Token})
	})

	req := httptest.NewRequest(http.MethodGet, "/verify?token=abc", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
}
