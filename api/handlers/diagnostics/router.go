package diagnostics

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Route("/diagnostics", func(r chi.Router) {
		r.Get("/stats", DockerStats(s))
		r.Get("/info", Info(s))
	})
}
