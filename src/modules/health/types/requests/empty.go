package requests

import z "github.com/Oudwins/zog"

type Empty struct{}

var EmptySchema = z.Struct(z.Shape{})
