package projects

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/projects"
)

type CreateProjectRequest struct {
	Repository    string  `json:"repository"`
	RequiresAuthn bool    `json:"requiresAuthn"`
	Username      *string `json:"username,omitempty"`
	Password      *string `json:"password,omitempty"`
}

func CreateProject(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateProjectRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

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
