package responses

import (
	"time"

	"go.boilerplate/src/modules/auth/services"

	z "github.com/Oudwins/zog"
)

type User struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	DisplayName   *string `json:"display_name"`
	AvatarURL     *string `json:"avatar_url"`
	EmailVerified bool    `json:"email_verified"`
}

type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

var UserSchema = z.Struct(z.Shape{
	"ID":            z.String().Required().UUID(),
	"Email":         z.String().Required().Email(),
	"DisplayName":   z.Ptr(z.String()),
	"AvatarURL":     z.Ptr(z.String()),
	"EmailVerified": z.Bool(),
})

var TokenSchema = z.Struct(z.Shape{
	"AccessToken": z.String().Required().Min(1),
	"ExpiresAt":   z.Time().Required(),
})

func UserFrom(u *services.AuthUser) User {
	return User{
		ID:            u.ID.String(),
		Email:         u.Email,
		DisplayName:   u.DisplayName,
		AvatarURL:     u.AvatarURL,
		EmailVerified: u.EmailVerified,
	}
}

func TokenFrom(t *services.TokenPair) Token {
	return Token{
		AccessToken: t.AccessToken,
		ExpiresAt:   t.ExpiresAt,
	}
}
