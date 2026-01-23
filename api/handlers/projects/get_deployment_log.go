package projects

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	deploymentlogs "skulpture/buang/components/deployment_logs"
	workflowwrappers "skulpture/buang/workers/workflow_wrappers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetDeploymentLogsRequest struct {
	ProjectId         int64 `schema:"-"`
	DeploymentId      int64 `schema:"-"`
	WaitForDeployment bool  `schema:"waitForDeployment,default:false"`
}

type GetDeploymentLogsResult = string

// @summary	Get deployment logs
// @tags		api.v1, deployment
// @security	ApiKeyAuth
// @param		projectId			path		int		required	"Project ID"
// @param		deploymentId		path		int		required	"Deployment ID"
// @param		waitForDeployment	query		bool	false		"Wait for deployment completion"
// @success	200					{object}	string
// @success	204					{object}	nil
// @failure	400					{object}	any
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment/{deploymentId}/logs [get]
func GetDeploymentLogs(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GetDeploymentLogsRequest

		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.ProjectId = int64(projectId)

		deploymentIdParam := chi.URLParam(r, "deploymentId")
		deploymentId, err := strconv.Atoi(deploymentIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.DeploymentId = int64(deploymentId)

		err = s.SchemaDecoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.WaitForDeployment {
			wp := workflowwrappers.HasDeployedParams{
				ProjectId:    req.ProjectId,
				DeploymentId: req.DeploymentId,
			}

			err = wp.Exec(r.Context(), s)
			if err != nil {
				slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		p := deploymentlogs.GetDeploymentLogParams{
			ProjectId:    int64(projectId),
			DeploymentId: int64(deploymentId),
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res.Logs == nil {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		err = app.WriteText(w, *res.Logs, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "get deployment logs", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
