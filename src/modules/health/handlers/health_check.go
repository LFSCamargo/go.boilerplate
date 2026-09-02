package HealthHandlers

import (
	"context"

	"go.boilerplate/src/modules/health/types/responses"
	"go.boilerplate/src/openapi"
)

type HealthOutput struct {
	Body responses.Health
}

func HealthCheckHandler(ctx context.Context, _ *struct{}) (*HealthOutput, error) {
	body := responses.OK()
	if err := openapi.ValidateResponse(responses.HealthSchema, body); err != nil {
		return nil, err
	}
	return &HealthOutput{Body: body}, nil
}
