package middleware

import (
	"log"
	"net/http"
	"strings"

	authProto "github.com/SeiFlow-3P2/auth_service/pkg/proto/v1"
	"github.com/gin-gonic/gin"
)

const UserIDHeader = "x-user-id"

type AuthMiddleware struct {
	authClient        authProto.AuthServiceClient
	protectedPrefixes []string
}

func NewAuthMiddleware(
	client authProto.AuthServiceClient,
	protectedPrefixes []string,
) *AuthMiddleware {
	return &AuthMiddleware{
		authClient:        client,
		protectedPrefixes: protectedPrefixes,
	}
}

func (am *AuthMiddleware) Handler(c *gin.Context) {
	isProtected := false
	path := c.Request.URL.Path

	for _, prefix := range am.protectedPrefixes {
		if strings.HasPrefix(path, prefix) {
			isProtected = true
			break
		}
	}

	if !isProtected {
		log.Printf("TEMP: not protected, skip auth at %s", path)
		c.Next()
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    16,
			"message": "Unauthorized: missing token",
			"details": []any{},
		})
		c.Abort()
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    16,
			"message": "Unauthorized: invalid token format",
			"details": []any{},
		})
		c.Abort()
		return
	}

	// TODO: auth service: validate token
	_ = parts[1] // token
	userID := "dde009e4-aad0-4570-b40a-cb0caee2a1c1"

	c.Set(UserIDHeader, userID)
	c.Next()
}
