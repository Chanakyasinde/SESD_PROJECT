package http

import (
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/domain/usecases"
	"inventory-management-system/backend/internal/infrastructure/security"
	"inventory-management-system/backend/internal/middleware"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthUsecase    usecases.AuthUsecase
	UserUsecase    usecases.UserUsecase
	ProductUsecase usecases.ProductUsecase
	OrderUsecase   usecases.OrderUsecase
	JWTService     *security.JWTService
	Logger         *slog.Logger
	AllowedOrigins []string
}

func NewRouter(dep Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recover())
	r.Use(middleware.CORS(dep.AllowedOrigins))
	r.Use(middleware.Logging(dep.Logger))

	authHandler := NewAuthHandler(dep.AuthUsecase)
	userHandler := NewUserHandler(dep.UserUsecase)
	productHandler := NewProductHandler(dep.ProductUsecase)
	orderHandler := NewOrderHandler(dep.OrderUsecase)

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)

		secured := api.Group("/")
		secured.Use(middleware.Auth(dep.JWTService))
		{
			userRoutes := secured.Group("/users")
			userRoutes.GET("/customers", middleware.RequireRoles(entities.RoleAdmin, entities.RoleStaff), userHandler.ListCustomers)

			productRoutes := secured.Group("/products")
			productRoutes.GET("", productHandler.List)
			productRoutes.GET("/:id", productHandler.GetByID)
			productRoutes.POST("", middleware.RequireRoles(entities.RoleAdmin), productHandler.Create)
			productRoutes.PUT("/:id", middleware.RequireRoles(entities.RoleAdmin), productHandler.Update)
			productRoutes.DELETE("/:id", middleware.RequireRoles(entities.RoleAdmin), productHandler.Delete)

			orderRoutes := secured.Group("/orders")
			orderRoutes.GET("", orderHandler.List)
			orderRoutes.GET("/:id", orderHandler.GetByID)
			orderRoutes.POST("", middleware.RequireRoles(entities.RoleAdmin, entities.RoleStaff, entities.RoleCustomer), orderHandler.Create)
			orderRoutes.PATCH("/:id/status", middleware.RequireRoles(entities.RoleAdmin, entities.RoleStaff), orderHandler.UpdateStatus)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
