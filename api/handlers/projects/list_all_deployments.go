package projects

import (
	"errors"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ListAllDeploymentsRequest struct {
	ProjectID int64  `schema:"-"`
	Limit     *int32 `schema:"limit,default:50" validate:"gte=1"`
	Page      *int32 `schema:"page,default:1" validate:"gte=1"`
}

type ListAllDeploymentsItem struct {
	ID         int64      `json:"id"`
	ProjectID  int64      `json:"projectId"`
	Url        *string    `json:"url"`
	Status     int16      `json:"status"`
	Sha        string     `json:"sha"`
	DeployedAt *time.Time `json:"deployedAt"`
}

// @summary	List all deployments
// @tags		api.v1, project
// @security	ApiKeyAuth
// @param		projectId	path		int	true	"Project ID"
// @param		limit		query		int	false	"Limit"
// @param		page		query		int	false	"Page"
// @success	200			{array}		ListAllDeploymentsItem
// @failure	400			{object}	any
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployments [get]
func ListAllDeployments(s app.ApplicationServices) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "list all deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req ListAllDeploymentsRequest

		err = s.SchemaDecoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "list all deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = validate.Struct(req)
		if err != nil {
			var validateErrs validator.ValidationErrors
			errs := map[string]string{}

			if errors.As(err, &validateErrs) {
				for _, e := range validateErrs {
					errs[e.Field()] = e.Error()
				}
			}

			app.WriteError(w, errs, http.StatusBadRequest)

			return
		}

		req.ProjectID = int64(projectId)

		p := deployments.ListDeploymentsParams{
			ProjectID: req.ProjectID,
			Limit:     *req.Limit,
			Page:      *req.Page,
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "list all deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ret := []ListAllDeploymentsItem{}
		for _, x := range res.Deployments {
			ret = append(ret, ListAllDeploymentsItem{
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
