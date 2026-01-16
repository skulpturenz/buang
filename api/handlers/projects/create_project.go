package projects

import (
	"fmt"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/projects"
)

func CreateProject(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := projects.CreateProjectParams{}

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
