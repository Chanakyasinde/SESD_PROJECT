package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"inventory-management-system/backend/pkg/response"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		roleAny, exists := c.Get(CtxRoleKey)
		if !exists {
			response.Error(c, http.StatusForbidden, "role not found in context", nil)
			c.Abort()
			return
		}

		role, ok := roleAny.(string)
		if !ok {
			response.Error(c, http.StatusForbidden, "invalid role context", nil)
			c.Abort()
			return
		}

		if _, ok = allowed[role]; !ok {
			response.Error(c, http.StatusForbidden, "insufficient role permissions", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
