package healthRouter

import (
	"net/http"

	HealthHandlers "go.boilerplate/src/modules/health/handlers"

	"github.com/danielgtaylor/huma/v2"
)

func Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Liveness",
		Description: "Returns `{status: ok}` when the process is up.",
		Tags:        []string{"Health"},
	}, HealthHandlers.HealthCheckHandler)
}
