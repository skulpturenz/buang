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

type FindProjectByRepositoryRequest struct {
	Repository string `schema:"repository" validate:"required"`
}

type FindProjectByRepositoryResponse struct {
	Id            int64     `json:"id"`
	Repository    string    `json:"repository"`
	RequiresAuthn bool      `json:"requiresAuthn"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ComposePath   string    `json:"composePath"`
}

// @summary	Find project by repository
// @tags		api.v1, projects
// @security	ApiKeyAuth
// @param		repository	query		string	required	"Repository"
// @success	200			{array}		FindProjectByRepositoryResponse
// @failure	400			{object}	any
// @failure	401
// @failure	500	{object}	string
// @router		/project [get]
func FindProjectByRepository(s app.ApplicationServices) http.HandlerFunc {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return func(w http.ResponseWriter, r *http.Request) {
		var req FindProjectByRepositoryRequest

		decoder, ok := s.GetGorillaSchemaDecoder()
		if !ok {
			slog.ErrorContext(r.Context(), "find project by repository", "err", "decoder not found")
			http.Error(w, errors.New("decoder not found").Error(), http.StatusInternalServerError)
			return
		}

		err := decoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "find project by repository", "err", err.Error())
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

		p := projects.FindProjectByRepositoryParams{
			Repository: req.Repository,
		}

		res, err := p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "find project by repository", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ret := FindProjectByRepositoryResponse{
			Id:            res.Project.GetId(),
			Repository:    res.Project.GetRepository(),
			RequiresAuthn: res.Project.GetRequiresAuthn(),
			CreatedAt:     res.Project.GetCreatedAt(),
			UpdatedAt:     res.Project.GetUpdatedAt(),
		}

		err = app.WriteJson(w, ret, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "find project by repository", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
