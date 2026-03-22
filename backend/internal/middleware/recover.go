package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"inventory-management-system/backend/pkg/response"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				response.Error(c, http.StatusInternalServerError, "internal server error", fmt.Sprintf("%v", rec))
				c.Abort()
			}
		}()
		c.Next()
	}
}
