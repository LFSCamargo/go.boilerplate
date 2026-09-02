package requests

import z "github.com/Oudwins/zog"

type VerifyEmailCode struct {
	Code string `json:"code" form:"code" doc:"6-digit verification code" minLength:"6" maxLength:"6" example:"123456"`
}

var VerifyEmailCodeSchema = z.Struct(z.Shape{
	"Code": z.String().Required().Min(6).Max(6),
})

type VerifyEmailToken struct {
	Token string `json:"token" query:"token"`
}

var VerifyEmailTokenSchema = z.Struct(z.Shape{
	"Token": z.String().Required().Min(1),
})
