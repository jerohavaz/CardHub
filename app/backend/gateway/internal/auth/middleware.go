package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey struct{}

func UserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKey{}).(string)
	return id, ok && id != ""
}

func Middleware(verifier *Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if isPublic(r) {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Fields(
				r.Header.Get("Authorization"),
			)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") {
				http.Error(
					w, "unauthorized", http.StatusUnauthorized,
				)
				return
			}

			userID, err := verifier.Verify(parts[1])
			if err != nil {
				http.Error(
					w, "unauthorized", http.StatusUnauthorized,
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				contextKey{},
				userID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublic(r *http.Request) bool {
	path := r.URL.Path

	if r.Method == http.MethodGet {
		switch path {
		case "/healthz",
			"/docs",
			"/openapi.json",
			"/openapi.yaml":
			return true
		}

		// Huma's Swagger UI may request additional assets.
		if strings.HasPrefix(path, "/docs/") {
			return true
		}
	}

	if r.Method == http.MethodPost {
		switch path {
		case "/v1/auth/register",
			"/v1/auth/login",
			"/v1/auth/refresh",
			"/v1/auth/logout":
			return true
		}
	}

	return false
}
