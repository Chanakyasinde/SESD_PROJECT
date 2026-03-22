package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"inventory-management-system/backend/internal/infrastructure/security"
	"inventory-management-system/backend/pkg/response"
)

const (
	CtxUserIDKey = "userId"
	CtxRoleKey   = "role"
)

func Auth(jwtSvc *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization header", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header", nil)
			c.Abort()
			return
		}

		claims, err := jwtSvc.ParseToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxRoleKey, claims.Role)
		c.Next()
	}
}
