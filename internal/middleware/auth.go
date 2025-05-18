package middleware

import (
	"log"
	"net/http"
	"strings"
)

const UserIDHeader = "x-user-id"

type AuthMiddleware struct {
	// authClient        authProto.AuthServiceClient
	protectedPrefixes []string
}

func NewAuthMiddleware(
	/*client authProto.AuthServiceClient, */
	protectedPrefixes []string,
) *AuthMiddleware {
	return &AuthMiddleware{
		// authClient: client,
		protectedPrefixes: protectedPrefixes,
	}
}

func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isProtected := false
		path := r.URL.Path

		for _, prefix := range am.protectedPrefixes {
			if strings.HasPrefix(path, prefix) {
				isProtected = true
				break
			}
		}

		if !isProtected {
			log.Printf("TEMP: not protected, skip auth at %s", path)
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, "")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}

		// TODO: auth service: validate token
		_ = parts[1] // token
		userID := "dde009e4-aad0-4570-b40a-cb0caee2a1c1"

		r.Header.Set(UserIDHeader, userID)
		next.ServeHTTP(w, r)
	})
}
