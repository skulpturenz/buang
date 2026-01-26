package projects

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"skulpture/buang/components/projects"
	workersinterfaces "skulpture/buang/workers/interfaces"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @summary	Delete a project
// @tags		api.v1, projects
// @security	ApiKeyAuth
// @param		projectId	path		int	true	"Project ID"
// @success	204			{object}	nil
// @failure	400			{object}	any
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId} [delete]
func DeleteProject(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "delete project", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d := deployments.FindActiveDeploymentsByProjectParams{
			ProjectID: int64(projectId),
		}

		activeDeployments, err := d.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "delete project", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, dply := range activeDeployments.Deployments {
			wp := workersinterfaces.BuangDeploymentParams{
				ProjectId:    int64(projectId),
				DeploymentId: dply.GetId(),
				Block:        true,
			}

			err = s.Workflows.BuangDeployment(r.Context(), wp)
			if err != nil {
				slog.ErrorContext(r.Context(), "delete project", "err", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		p := projects.DeleteProjectParams{
			Id: int64(projectId),
		}

		_, err = p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "delete project", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
