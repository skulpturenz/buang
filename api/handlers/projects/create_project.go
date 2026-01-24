package projects

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/projects"

	"github.com/go-playground/validator/v10"
)

type CreateProjectRequest struct {
	Repository    string  `json:"repository" validate:"required"`
	RequiresAuthn bool    `json:"requiresAuthn"`
	Username      *string `json:"username,omitempty" validate:"required_with=Password,required_if=RequiresAuthn true"`
	Password      *string `json:"password,omitempty" validate:"required_with=Username,required_if=RequiresAuthn true"`
	ComposePath   string  `json:"composePath" validate:"required"`
}

// @summary	Create a project
// @tags		api.v1, projects
// @security	ApiKeyAuth
// @param		projectDetails	body		CreateProjectRequest	true	"Project details"
// @success	200				{object}	int
// @failure	400				{object}	any
// @failure	401
// @failure	500	{object}	string
// @router		/project [post]
func CreateProject(s app.ApplicationServices) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProjectRequest

		err := json.NewDecoder(r.Body).Decode(&req)
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

		p := projects.CreateProjectParams(req)

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "create project", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = fmt.Fprintf(w, "%v", res.Id)
		if err != nil {
			slog.ErrorContext(r.Context(), "create project", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
