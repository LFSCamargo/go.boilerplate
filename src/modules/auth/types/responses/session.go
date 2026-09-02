package responses

import z "github.com/Oudwins/zog"

type Session struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type Me struct {
	User User `json:"user"`
}

type RegisterCreated struct {
	User                 User `json:"user"`
	RequiresVerification bool `json:"requires_verification"`
}

type OK struct {
	OK bool `json:"ok"`
}

type Verified struct {
	OK       bool `json:"ok"`
	Verified bool `json:"verified"`
}

var SessionSchema = z.Struct(z.Shape{
	"User":  UserSchema,
	"Token": TokenSchema,
})

var MeSchema = z.Struct(z.Shape{
	"User": UserSchema,
})

var RegisterCreatedSchema = z.Struct(z.Shape{
	"User":                 UserSchema,
	"RequiresVerification": z.Bool(),
})

var OKSchema = z.Struct(z.Shape{
	"OK": z.Bool(),
})

var VerifiedSchema = z.Struct(z.Shape{
	"OK":       z.Bool(),
	"Verified": z.Bool(),
})
