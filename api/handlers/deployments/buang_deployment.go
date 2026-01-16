package deployments

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	"skulpture/buang/components/deployments"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func BuangDeployment(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deploymentIdParam := chi.URLParam(r, "deploymentId")
		deploymentId, err := strconv.Atoi(deploymentIdParam)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployments", "err", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p := deployments.BuangDeploymentParams{
			ID: int64(deploymentId),
		}

		_, err = p.Exec(r.Context(), &s)
		if err != nil {
			slog.ErrorContext(r.Context(), "buang deployment", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
