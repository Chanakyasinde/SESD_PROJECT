package http

import (
	"net/http"
	"strconv"

	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
	"inventory-management-system/backend/internal/domain/usecases"
	"inventory-management-system/backend/internal/middleware"
	"inventory-management-system/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	orderUC usecases.OrderUsecase
}

func NewOrderHandler(orderUC usecases.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUC: orderUC}
}

type createOrderRequest struct {
	CustomerID string                          `json:"customerId"`
	Items      []usecases.CreateOrderItemInput `json:"items" binding:"required,dive"`
}

type updateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed shipped delivered cancelled"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	actorIDHex, _ := c.Get(middleware.CtxUserIDKey)
	actorID, err := primitive.ObjectIDFromHex(actorIDHex.(string))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid actor id", nil)
		return
	}

	roleAny, _ := c.Get(middleware.CtxRoleKey)
	role, _ := roleAny.(string)
	switch role {
	case entities.RoleAdmin, entities.RoleStaff:
		if req.CustomerID == "" {
			response.Error(c, http.StatusBadRequest, "customerId is required", nil)
			return
		}
	case entities.RoleCustomer:
		req.CustomerID = actorID.Hex()
	default:
		response.Error(c, http.StatusForbidden, "insufficient role permissions", nil)
		return
	}

	result, err := h.orderUC.Create(c.Request.Context(), actorID, usecases.CreateOrderInput{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "order created", result)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id", nil)
		return
	}

	var req updateStatusRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	actorIDHex, _ := c.Get(middleware.CtxUserIDKey)
	actorID, err := primitive.ObjectIDFromHex(actorIDHex.(string))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid actor id", nil)
		return
	}

	if err = h.orderUC.UpdateStatus(c.Request.Context(), actorID, id, req.Status); err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "order status updated", nil)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order id", nil)
		return
	}
	result, err := h.orderUC.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	roleAny, _ := c.Get(middleware.CtxRoleKey)
	role, _ := roleAny.(string)
	if role == entities.RoleCustomer {
		actorIDHex, _ := c.Get(middleware.CtxUserIDKey)
		actorID, parseErr := primitive.ObjectIDFromHex(actorIDHex.(string))
		if parseErr != nil {
			response.Error(c, http.StatusUnauthorized, "invalid actor id", nil)
			return
		}
		if result.CustomerID != actorID {
			response.Error(c, http.StatusForbidden, "insufficient role permissions", nil)
			return
		}
	}
	response.Success(c, http.StatusOK, "order fetched", result)
}

func (h *OrderHandler) List(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 64)
	sortDir, _ := strconv.Atoi(c.DefaultQuery("sortDir", "-1"))

	actorIDHex, _ := c.Get(middleware.CtxUserIDKey)
	actorID, err := primitive.ObjectIDFromHex(actorIDHex.(string))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid actor id", nil)
		return
	}

	roleAny, _ := c.Get(middleware.CtxRoleKey)
	role, _ := roleAny.(string)

	var customerID *primitive.ObjectID
	switch role {
	case entities.RoleCustomer:
		customerID = &actorID
	case entities.RoleAdmin, entities.RoleStaff:
		if customerIDRaw := c.Query("customerId"); customerIDRaw != "" {
			if parsed, err := primitive.ObjectIDFromHex(customerIDRaw); err == nil {
				customerID = &parsed
			}
		}
	default:
		response.Error(c, http.StatusForbidden, "insufficient role permissions", nil)
		return
	}

	filter := repositories.OrderFilter{
		Status:     c.Query("status"),
		CustomerID: customerID,
		Pagination: entities.Pagination{
			Page:     page,
			PageSize: pageSize,
			SortBy:   c.DefaultQuery("sortBy", "createdAt"),
			SortDir:  sortDir,
		},
	}

	result, total, err := h.orderUC.List(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "orders fetched", gin.H{"items": result, "total": total, "page": page, "pageSize": pageSize})
}
