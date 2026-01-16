package projects

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/projects"
	"time"
)

type ListProjectRequest struct {
	Limit int32 `schema:"limit,default:50"`
	Page  int32 `schema:"page,default:1"`
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

func ListAllProjects(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ListProjectRequest

		err := s.SchemaDecoder.Decode(&req, r.URL.Query())
		if err != nil {
			slog.ErrorContext(r.Context(), "list all projects", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := projects.ListProjectParams(req)

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
