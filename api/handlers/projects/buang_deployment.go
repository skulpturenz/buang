package projects

import (
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/workflows"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.temporal.io/sdk/client"
)

// @summary	Spin down a preview deployment
// @tags		api.v1, project
// @param		projectId		path		int	true	"Project ID"
// @param		deploymentId	path		int	true	"Deployment ID"
// @success	204				{object}	nil
// @failure	403
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment/{deploymentId} [delete]
func BuangDeployment(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		deploymentIdParam := chi.URLParam(r, "deploymentId")
		deploymentId, err := strconv.Atoi(deploymentIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		workflowId := fmt.Sprintf("buang-project-%v-deployment-%v", projectId, deploymentId)
		options := client.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueBuang,
		}

		_, err = s.Temporal.ExecuteWorkflow(r.Context(), options, workflows.Buang, projectId, deploymentId)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
