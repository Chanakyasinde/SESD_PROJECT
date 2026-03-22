package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"inventory-management-system/backend/internal/domain/usecases"
	"inventory-management-system/backend/pkg/response"
)

type AuthHandler struct {
	authUC usecases.AuthUsecase
}

func NewAuthHandler(authUC usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	result, err := h.authUC.Register(c.Request.Context(), usecases.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "user registered", result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	result, err := h.authUC.Login(c.Request.Context(), usecases.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login successful", result)
}
