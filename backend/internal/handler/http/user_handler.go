package http

import (
	"net/http"
	"strconv"

	"inventory-management-system/backend/internal/domain/usecases"
	"inventory-management-system/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC usecases.UserUsecase
}

func NewUserHandler(userUC usecases.UserUsecase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

func (h *UserHandler) ListCustomers(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "50"), 10, 64)

	items, total, err := h.userUC.ListCustomers(c.Request.Context(), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "customers fetched", gin.H{
		"items":    items,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
