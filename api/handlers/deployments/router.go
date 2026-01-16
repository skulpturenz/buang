package deployments

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Route("/deployment/{deploymentId}", func(r chi.Router) {
		r.Put("/buang", BuangDeployment(s))
	})
}
