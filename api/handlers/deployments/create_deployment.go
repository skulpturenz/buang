package deployments

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	workersinterfaces "skulpture/buang/workers/interfaces"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CreateDeploymentRequest struct {
	Branch            string         `json:"branch" validate:"required" schema:"-"`
	Sha               string         `json:"sha" validate:"required" schema:"-"`
	ServiceEntrypoint string         `json:"serviceEntrypoint" validate:"required,hostname_port" schema:"-"`
	Env               map[string]any `json:"env" schema:"-"`
	WaitForDeployment bool           `schema:"waitForDeployment,default:false" swaggerignore:"true"`
}

// @summary	Spin up a preview deployment
// @tags		api.v1, deployments
// @security	ApiKeyAuth
// @param		projectId			path		int						true	"Project ID"
// @param		deploymentDetails	body		CreateDeploymentRequest	true	"Deployment details"
// @param		waitForDeployment	query		bool					false	"Wait for deployment completion"
// @success	200					{object}	int
// @failure	400					{object}	any
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment [post]
func CreateDeployment(s app.ApplicationServices) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
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

		decoder, ok := s.GetGorillaSchemaDecoder()
		if !ok {
			slog.ErrorContext(r.Context(), "create deployment", "err", "decoder not found")
			http.Error(w, errors.New("decoder not found").Error(), http.StatusInternalServerError)
			return
		}

		err = decoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := deployments.CreateDeploymentParams{
			ProjectID:         int64(projectId),
			Branch:            req.Branch,
			Sha:               req.Sha,
			ServiceEntrypoint: req.ServiceEntrypoint,
			Env:               req.Env,
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wp := workersinterfaces.CreateDeploymentParams{
			ProjectId:    p.ProjectID,
			DeploymentId: res.Id,
			Block:        req.WaitForDeployment,
		}

		workflows, ok := s.GetWorkflows()
		if !ok {
			slog.ErrorContext(r.Context(), "create deployment", "err", "workflows not found")
			http.Error(w, errors.New("workflows not found").Error(), http.StatusInternalServerError)
			return
		}

		err = workflows.CreateDeployment(r.Context(), wp)
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprintf(w, "%v", res.Id)
		if err != nil {
			slog.ErrorContext(r.Context(), "create deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
