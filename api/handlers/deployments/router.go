package deployments

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Route("/project/{projectId}/deployment", func(r chi.Router) {

		r.Post("/{projectId}/deployment", CreateDeployment(s))
		r.Get("/{projectId}/deployment/{deploymentId}", FindDeploymentById(s))
		r.Delete("/{projectId}/deployment/{deploymentId}", BuangDeployment(s))
	})
}
