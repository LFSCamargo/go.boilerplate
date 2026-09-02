package requests

import z "github.com/Oudwins/zog"

type Register struct {
	Email       string `json:"email" form:"email" doc:"Account email" example:"ada@example.com" format:"email"`
	Password    string `json:"password" form:"password" doc:"At least 8 characters" minLength:"8" example:"password12"`
	DisplayName string `json:"display_name,omitempty" form:"display_name" doc:"Optional display name" example:"Ada"`
	AvatarURL   string `json:"-"`
}

var RegisterSchema = z.Struct(z.Shape{
	"Email":       z.String().Required().Email(),
	"Password":    z.String().Required().Min(8),
	"DisplayName": z.String().Optional(),
})
