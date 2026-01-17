package projects

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	constantstaskqueues "skulpture/buang/constants/task_queues"
	"skulpture/buang/workers/workflows"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.temporal.io/sdk/client"
)

type BuangBranchRequest struct {
	Branch string `json:"repository" validate:"required"`
}

// @summary	Spin down a preview branch
// @tags		api.v1, project
// @security	ApiKeyAuth
// @param		projectId		path		int					true	"Project ID"
// @param		branchDetails	body		BuangBranchRequest	true	"Branch details"
// @success	204				{object}	nil
// @failure	401
// @failure	500	{object}	string
// @router		/project/{projectId}/branch [delete]
func BuangBranch(s app.ApplicationServices) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		projectIdParam := chi.URLParam(r, "projectId")
		projectId, err := strconv.Atoi(projectIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang branch", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var req BuangBranchRequest

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

		workflowId := fmt.Sprintf("buang-project-%v-branch-%v", projectId, req.Branch)
		options := client.StartWorkflowOptions{
			ID:        workflowId,
			TaskQueue: constantstaskqueues.QueueBuang,
		}

		_, err = s.Temporal.ExecuteWorkflow(r.Context(), options, workflows.BuangBranch, int64(projectId), req.Branch)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang branch", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
