package requests

import z "github.com/Oudwins/zog"

type Login struct {
	Email    string `json:"email" form:"email" doc:"Account email" example:"ada@example.com" format:"email"`
	Password string `json:"password" form:"password" doc:"Account password" minLength:"1" example:"password12"`
}

var LoginSchema = z.Struct(z.Shape{
	"Email":    z.String().Required().Email(),
	"Password": z.String().Required().Min(1),
})
