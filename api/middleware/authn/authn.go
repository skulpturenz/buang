package authn

import (
	"context"
	"net/http"
	enumsauthnmethod "skulpture/buang/enums/authn_method"
)

type AuthnMiddlewareConfig struct {
	ApiKey string
}

func (config *AuthnMiddlewareConfig) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(("X-API-Key"))

		if apiKey != config.ApiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		ctx := context.WithValue(r.Context(), "authn_method", enumsauthnmethod.ApiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
