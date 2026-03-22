package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/repositories"
	"inventory-management-system/backend/internal/domain/usecases"
	"inventory-management-system/backend/pkg/response"
)

type ProductHandler struct {
	productUC usecases.ProductUsecase
}

func NewProductHandler(productUC usecases.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUC: productUC}
}

type productRequest struct {
	Name              string  `json:"name" binding:"required,min=2,max=120"`
	Description       string  `json:"description" binding:"max=500"`
	Category          string  `json:"category" binding:"required,min=2,max=80"`
	Price             float64 `json:"price" binding:"required,gt=0"`
	StockQuantity     int64   `json:"stockQuantity" binding:"required,gte=0"`
	LowStockThreshold int64   `json:"lowStockThreshold" binding:"required,gte=0"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req productRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	result, err := h.productUC.Create(c.Request.Context(), entities.Product{
		Name:              req.Name,
		Description:       req.Description,
		Category:          req.Category,
		Price:             req.Price,
		StockQuantity:     req.StockQuantity,
		LowStockThreshold: req.LowStockThreshold,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "product created", result)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid product id", nil)
		return
	}

	var req productRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	result, err := h.productUC.Update(c.Request.Context(), id, entities.Product{
		Name:              req.Name,
		Description:       req.Description,
		Category:          req.Category,
		Price:             req.Price,
		StockQuantity:     req.StockQuantity,
		LowStockThreshold: req.LowStockThreshold,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "product updated", result)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid product id", nil)
		return
	}
	if err = h.productUC.Delete(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "product deleted", nil)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid product id", nil)
		return
	}
	result, err := h.productUC.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "product fetched", result)
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("pageSize", "10"), 10, 64)
	sortDir, _ := strconv.Atoi(c.DefaultQuery("sortDir", "1"))

	var minPricePtr *float64
	if minPriceRaw := c.Query("minPrice"); minPriceRaw != "" {
		minPrice, err := strconv.ParseFloat(minPriceRaw, 64)
		if err == nil {
			minPricePtr = &minPrice
		}
	}
	var maxPricePtr *float64
	if maxPriceRaw := c.Query("maxPrice"); maxPriceRaw != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceRaw, 64)
		if err == nil {
			maxPricePtr = &maxPrice
		}
	}
	var lowStockPtr *bool
	if lowStockRaw := c.Query("lowStock"); lowStockRaw != "" {
		value := lowStockRaw == "true"
		lowStockPtr = &value
	}

	filter := repositories.ProductFilter{
		Category: c.Query("category"),
		Search:   c.Query("search"),
		MinPrice: minPricePtr,
		MaxPrice: maxPricePtr,
		LowStock: lowStockPtr,
		Pagination: entities.Pagination{
			Page:     page,
			PageSize: pageSize,
			SortBy:   c.DefaultQuery("sortBy", "createdAt"),
			SortDir:  sortDir,
		},
	}

	result, total, err := h.productUC.List(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "products fetched", gin.H{"items": result, "total": total, "page": page, "pageSize": pageSize})
}
