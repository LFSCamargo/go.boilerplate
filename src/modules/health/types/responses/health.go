package responses

import z "github.com/Oudwins/zog"

type Health struct {
	Status string `json:"status"`
}

var HealthSchema = z.Struct(z.Shape{
	"Status": z.String().Required().OneOf([]string{"ok"}),
})

func OK() Health {
	return Health{Status: "ok"}
}
