package transport

import (
	"context"
	"net/http"
	"strings"
)

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := h.tokenManager.ParseToken(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)

		r = r.WithContext(ctx)
	})
}
