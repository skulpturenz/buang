package deployments

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Route("/project/{projectId}/deployment", func(r chi.Router) {

		r.Post("/", CreateDeployment(s))
		r.Get("/{deploymentId}", FindDeploymentById(s))
		r.Delete("/{deploymentId}", BuangDeployment(s))
	})
}
