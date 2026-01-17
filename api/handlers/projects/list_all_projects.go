package projects

import (
	"errors"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/projects"
	"time"

	"github.com/go-playground/validator/v10"
)

type ListProjectRequest struct {
	Limit *int32 `schema:"limit,default:50" validate:"gte=1"`
	Page  *int32 `schema:"page,default:1" validate:"gte=1"`
}

type Project struct {
	Id            int64     `json:"id"`
	Repository    string    `json:"repository"`
	RequiresAuthn bool      `json:"requiresAuthn"`
	Username      *string   `json:"username,omitempty"`
	Password      *string   `json:"password,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// @summary	List all projects
// @tags		api.v1, projects
// @security	ApiKeyAuth
// @param		limit	query	int	false	"Limit"
// @param		page	query	int	false	"Page"
// @success	200		{array}	Project
// @failure	401
// @failure	500	{object}	string
// @router		/projects [get]
func ListAllProjects(s app.ApplicationServices) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		var req ListProjectRequest

		err := s.SchemaDecoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "list all projects", "err", err.Error())
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

		p := projects.ListProjectParams{
			Limit: *req.Limit,
			Page:  *req.Page,
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "list all projects", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ret := []Project{}
		for _, x := range res.Projects {
			ret = append(ret, Project{
				Id:            x.GetId(),
				Repository:    x.GetRepository(),
				RequiresAuthn: x.GetRequiresAuthn(),
				Username:      x.GetUsername(),
				Password:      x.GetPassword(),
				CreatedAt:     x.GetCreatedAt(),
				UpdatedAt:     x.GetUpdatedAt(),
			})
		}

		err = app.WriteJson(w, ret, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "list all projects", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
