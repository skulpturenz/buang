package projects

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/workflows"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.temporal.io/sdk/client"
)

type CreateDeploymentRequest struct {
	Branch            string         `json:"branch" validate:"required"`
	Sha               string         `json:"sha" validate:"required"`
	ServiceEntrypoint string         `json:"serviceEntrypoint" validate:"required,hostname_port"`
	Env               map[string]any `json:"env"`
}

// @summary	Spin up a preview deployment
// @tags		api.v1, project
// @security	ApiKeyAuth
// @param		projectId			path		int						true	"Project ID"
// @param		deploymentDetails	body		CreateDeploymentRequest	true	"Deployment details"
// @success	200					{object}	int
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/deployment [post]
func CreateDeployment(s app.ApplicationServices) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

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

		workflowId := fmt.Sprintf("create-project-%v-deployment-%v", projectId, res.Id)
		options := client.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueDeployment,
		}

		_, err = s.Temporal.ExecuteWorkflow(r.Context(), options, workflows.Deploy, int64(projectId), int64(res.Id))
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
