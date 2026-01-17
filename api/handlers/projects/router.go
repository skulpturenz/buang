package projects

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", ListAllProjects(s))
	})

	r.Route("/project", func(r chi.Router) {
		r.Post("/", CreateProject(s))
		r.Get("/", FindProjectByRepository(s))

		r.Post("/{projectId}/deployment", CreateDeployment(s))

		r.Get("/{projectId}/deployment/{deploymentId}", FindDeploymentById(s))
		r.Delete("/{projectId}/deployment/{deploymentId}", BuangDeployment(s))

		r.Get("/{projectId}/deployments", ListAllDeployments(s))

		r.Delete("/{projectId}/branch", BuangBranch(s))
	})
}
