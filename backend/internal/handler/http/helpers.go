package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/pkg/response"
)

func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, entities.ErrInvalidInput):
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
	case errors.Is(err, entities.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
	case errors.Is(err, entities.ErrForbidden):
		response.Error(c, http.StatusForbidden, err.Error(), nil)
	case errors.Is(err, entities.ErrNotFound):
		response.Error(c, http.StatusNotFound, err.Error(), nil)
	case errors.Is(err, entities.ErrConflict):
		response.Error(c, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, entities.ErrInvalidTransition):
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
	case mongo.IsDuplicateKeyError(err):
		response.Error(c, http.StatusConflict, "resource already exists", nil)
	case strings.Contains(strings.ToLower(err.Error()), "insufficient stock"):
		response.Error(c, http.StatusConflict, err.Error(), nil)
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error", nil)
	}
}
