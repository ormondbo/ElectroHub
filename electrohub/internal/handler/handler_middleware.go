package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	authpb "contracts/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userIDKey struct{}
type accessTokenKey struct{}

func (api *API) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		response, err := api.auth.ValidateToken(r.Context(), &authpb.ValidateTokenRequest{AccessToken: parts[1]})
		if err != nil {
			log.Printf("validate access token: %v", err)
			if status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.DeadlineExceeded {
				writeError(w, http.StatusServiceUnavailable, "authentication service unavailable")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to validate access token")
			return
		}
		if response == nil || !response.GetValid() {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey{}, response.GetUserId())
		ctx = context.WithValue(ctx, accessTokenKey{}, parts[1])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
