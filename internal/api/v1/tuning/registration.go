package tuning

import (
	"net/http"

	"github.com/bigstack-oss/cube-cos-api/internal/api"
	definition "github.com/bigstack-oss/cube-cos-api/internal/definition/v1"
)

var (
	handlers = []api.Handler{
		{
			Version: api.V1,
			Method:  http.MethodGet,
			Path:    "/tunings",
			Func:    getTunings,
		},
		{
			Version: api.V1,
			Method:  http.MethodGet,
			Path:    "/tunings/specs",
			Func:    getTuningSpecs,
		},
		{
			Version: api.V1,
			Method:  http.MethodPut,
			Path:    "/tunings/:ParameterName",
			Func:    applyTuning,
		},
		{
			Version: api.V1,
			Method:  http.MethodPut,
			Path:    "/tunings",
			Func:    applyTunings,
		},
		{
			Version: api.V1,
			Method:  http.MethodPut,
			Path:    "/tunings/:ParameterName/status",
			Func:    updateTuningStatus,
		},
		{
			Version: api.V1,
			Method:  http.MethodDelete,
			Path:    "/tuning/:ParameterName",
			Func:    deleteTuning,
		},
		{
			Version: api.V1,
			Method:  http.MethodDelete,
			Path:    "/tunings",
			Func:    deleteTunings,
		},
	}
)

func init() {
	api.RegisterHandlersToRoles(
		tunings,
		handlers,
		definition.RoleControl,
		definition.RoleCompute,
	)
}
