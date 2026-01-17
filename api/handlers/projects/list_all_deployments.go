package projects

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ListDeploymentsRequest struct {
	ProjectID int64 `schema:"-"`
	Limit     int32 `schema:"limit,default:50"`
	Page      int32 `schema:"page,default:1"`
}

type Deployment struct {
	ID         int64      `json:"id"`
	ProjectID  int64      `json:"projectId"`
	Url        *string    `json:"url"`
	Status     int16      `json:"status"`
	Sha        *string    `json:"sha"`
	DeployedAt *time.Time `json:"deployedAt"`
}

// @summary	List all deployments
// @tags		api.v1, project
// @security	ApiKeyAuth
// @param		projectId	path	int	true	"Project ID"
// @param		limit		query	int	false	"Limit"
// @param		page		query	int	false	"Page"
// @success	200			{array}	Deployment
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployments [get]
func ListAllDeployments(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "list all deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req ListDeploymentsRequest

		err = s.SchemaDecoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "list all deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.ProjectID = int64(projectId)

		p := deployments.ListDeploymentsParams(req)

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "list all deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ret := []Deployment{}
		for _, x := range res.Deployments {
			ret = append(ret, Deployment{
				ID:         x.GetId(),
				ProjectID:  x.GetProjectId(),
				Url:        x.GetUrl(),
				Status:     x.GetStatus(),
				Sha:        x.GetSha(),
				DeployedAt: x.GetDeployedAt(),
			})
		}

		err = app.WriteJson(w, ret, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "list all deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
