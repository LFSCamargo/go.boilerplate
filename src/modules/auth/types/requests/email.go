package requests

import z "github.com/Oudwins/zog"

type Email struct {
	Email string `json:"email" form:"email" doc:"Account email" example:"ada@example.com" format:"email"`
}

var EmailSchema = z.Struct(z.Shape{
	"Email": z.String().Required().Email(),
})
