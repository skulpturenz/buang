package deployments

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type FindDeploymentByIdResponse struct {
	ID         int64      `json:"id"`
	ProjectID  int64      `json:"projectId"`
	Url        *string    `json:"url"`
	Status     int16      `json:"status"`
	Sha        string     `json:"sha"`
	DeployedAt *time.Time `json:"deployedAt"`
}

// @summary	Find deployment by ID
// @tags		api.v1, deployments
// @security	ApiKeyAuth
// @param		projectId		path		int	required	"Project ID"
// @param		deploymentId	path		int	required	"Deployment ID"
// @success	200				{object}	FindDeploymentByIdResponse
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment/{deploymentId} [get]
func FindDeploymentById(s app.ApplicationServices) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "find deployment by id", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		deploymentIdParam := chi.URLParam(r, "deploymentId")
		deploymentId, err := strconv.Atoi(deploymentIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "find deployment by id", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := deployments.FindDeploymentByIdParams{
			ID:        int64(deploymentId),
			ProjectId: int64(projectId),
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "find deployment by id", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ret := FindDeploymentByIdResponse{
			ID:         res.Deployment.GetId(),
			ProjectID:  res.Deployment.GetProjectId(),
			Url:        res.Deployment.GetUrl(),
			Status:     res.Deployment.GetStatus(),
			Sha:        res.Deployment.GetSha(),
			DeployedAt: res.Deployment.GetDeployedAt(),
		}

		err = app.WriteJson(w, ret, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "find deployment by id", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
