package projects

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", ListAllProjects(s))
	})
}
