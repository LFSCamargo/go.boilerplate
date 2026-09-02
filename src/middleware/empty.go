package middleware

import z "github.com/Oudwins/zog"

// Empty is a request body with no fields (GET /health, GET /me, POST /logout).
type Empty struct{}

var EmptySchema = z.Struct(z.Shape{})
