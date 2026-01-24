package diagnostics

import (
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	constantsenvs "skulpture/buang/constants/envs"
	"time"
)

type InfoRequest struct {
}

type InfoResult struct {
	Version  string `json:"version"`
	Timezone string `json:"timezone"`
	Uptime   string `json:"uptime"`
}

// @summary	Service info
// @tags		api.v1, diagnostics
// @security	ApiKeyAuth
// @success	200	{object}	InfoResult
// @failure	401
// @failure	500	{object}	string
// @router		/diagnostics/info [get]
func Info(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zoneName, _ := time.Now().Zone()

		about := InfoResult{
			Version:  constantsenvs.BUANG_VERSION,
			Timezone: zoneName,
			Uptime:   constantsenvs.Uptime().String(),
		}

		err := app.WriteJson(w, about, http.StatusOK)
		if err != nil {
			slog.ErrorContext(r.Context(), "version", "err", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
