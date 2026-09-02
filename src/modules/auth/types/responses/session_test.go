package responses

import (
	"testing"
	"time"

	"go.boilerplate/src/modules/auth/services"

	"github.com/google/uuid"
)

func TestSessionSchema_AcceptsValidSession(t *testing.T) {
	session := Session{
		User: UserFrom(&services.AuthUser{
			ID:    uuid.MustParse("11111111-1111-4111-8111-111111111111"),
			Email: "ada@example.com",
		}),
		Token: TokenFrom(&services.TokenPair{
			AccessToken: "jwt-token",
			ExpiresAt:   time.Now().Add(time.Hour),
		}),
	}
	if issues := SessionSchema.Validate(&session); issues != nil {
		t.Fatalf("unexpected issues: %#v", issues)
	}
}

func TestOKSchema_AcceptsTrue(t *testing.T) {
	ok := OK{OK: true}
	if issues := OKSchema.Validate(&ok); issues != nil {
		t.Fatalf("unexpected issues: %#v", issues)
	}
}
