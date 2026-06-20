package mcp

import (
	"skulpture/buang/app"

	"github.com/go-chi/chi/v5"
)

func Router(s app.ApplicationServices, r chi.Router) {
	r.Post("/", Handler(s))
}
