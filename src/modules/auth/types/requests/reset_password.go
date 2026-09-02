package requests

import z "github.com/Oudwins/zog"

type ResetPassword struct {
	Email       string `json:"email,omitempty" form:"email" doc:"Required when using an OTP code" example:"ada@example.com"`
	Code        string `json:"code,omitempty" form:"code" doc:"6-digit OTP from the reset email" example:"123456"`
	Token       string `json:"token,omitempty" form:"token" doc:"Magic-link token (alternative to email+code)"`
	NewPassword string `json:"new_password" form:"new_password" doc:"New password" minLength:"8" example:"newpassword1"`
}

var ResetPasswordSchema = z.Struct(z.Shape{
	"Email":       z.String().Email().Optional(),
	"Code":        z.String().Optional(),
	"Token":       z.String().Optional(),
	"NewPassword": z.String().Required().Min(8),
})

func (r ResetPassword) HasToken() bool {
	return r.Token != ""
}

func (r ResetPassword) HasCode() bool {
	return r.Email != "" && r.Code != ""
}
