package diagnostics

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	constantsenvs "skulpture/buang/constants/envs"
)

type VersionRequest struct {
}

type VersionResult = string

// @summary	API Version
// @tags		api.v1, diagnostics
// @security	ApiKeyAuth
// @success	200	{object}	VersionResult
// @failure	401
// @failure	500	{object}	string
// @router		/diagnostics/version [get]
func Version(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.WriteText(w, constantsenvs.BUANG_VERSION, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "version", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
