package projects

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	workflowwrappers "skulpture/buang/workers/workflow_wrappers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @summary	Spin down a preview deployment
// @tags		api.v1, project
// @security	ApiKeyAuth
// @param		projectId		path		int	true	"Project ID"
// @param		deploymentId	path		int	true	"Deployment ID"
// @success	204				{object}	nil
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment/{deploymentId} [delete]
func BuangDeployment(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		deploymentIdParam := chi.URLParam(r, "deploymentId")
		deploymentId, err := strconv.Atoi(deploymentIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		wp := workflowwrappers.BuangDeploymentParams{
			ProjectId:    int64(projectId),
			DeploymentId: int64(deploymentId),
		}

		err = wp.Exec(r.Context(), s)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
