package projects

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/workflows"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.temporal.io/sdk/client"
)

type CreateDeploymentRequest struct {
	Sha *string `json:"sha"`
}

// @summary	Spin up a preview deployment
// @tags		api.v1, project
// @param		projectId			path		int						true	"Project ID"
// @param		deploymentDetails	body		CreateDeploymentRequest	true	"Deployment details"
// @success	200					{object}	int
// @failure	403
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment [post]
func CreateDeployment(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req CreateDeploymentRequest

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		p := deployments.CreateDeploymentParams{
			ProjectID: int64(projectId),
			Sha:       req.Sha,
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		workflowId := fmt.Sprintf("create-project-%v-deployment-%v", projectId, res.Id)
		options := client.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueDeployment,
		}

		_, err = s.Temporal.ExecuteWorkflow(r.Context(), options, workflows.Deploy, projectId, res.Id)
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = fmt.Fprintf(w, "%v", res.Id)
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
