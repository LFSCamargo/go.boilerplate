package responses

import (
	"testing"
	"time"

	"go.boilerplate/src/modules/auth/services"

	"github.com/google/uuid"
)

func TestUserSchema_AcceptsUserFrom(t *testing.T) {
	name := "Ada"
	user := UserFrom(&services.AuthUser{
		ID:            uuid.MustParse("11111111-1111-4111-8111-111111111111"),
		Email:         "ada@example.com",
		DisplayName:   &name,
		EmailVerified: true,
	})
	if issues := UserSchema.Validate(&user); issues != nil {
		t.Fatalf("unexpected issues: %#v", issues)
	}
}

func TestUserSchema_RejectsInvalidEmail(t *testing.T) {
	user := User{ID: "11111111-1111-4111-8111-111111111111", Email: "nope"}
	if issues := UserSchema.Validate(&user); issues == nil {
		t.Fatal("expected validation issues")
	}
}

func TestTokenSchema_AcceptsTokenFrom(t *testing.T) {
	token := TokenFrom(&services.TokenPair{
		AccessToken: "jwt-token",
		ExpiresAt:   time.Now().Add(time.Hour),
	})
	if issues := TokenSchema.Validate(&token); issues != nil {
		t.Fatalf("unexpected issues: %#v", issues)
	}
}
