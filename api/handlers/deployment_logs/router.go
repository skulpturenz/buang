package deploymentlogs

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Route("/project/{projectId}/deployment/{deploymentId}/logs", func(r chi.Router) {
		r.Get("/", GetDeploymentLog(s))
	})
}
