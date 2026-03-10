package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	httpHandler "inventory-management-system/backend/internal/handler/http"
	"inventory-management-system/backend/internal/infrastructure/config"
	"inventory-management-system/backend/internal/infrastructure/db"
	"inventory-management-system/backend/internal/infrastructure/security"
	mongoRepo "inventory-management-system/backend/internal/repository/mongo"
	"inventory-management-system/backend/internal/usecase"
	"inventory-management-system/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	appLogger := logger.New()

	ctx := context.Background()
	mongoClient, collections, err := db.Connect(ctx, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = mongoClient.Disconnect(shutdownCtx)
	}()

	jwtSvc := security.NewJWTService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTExpiresInHour)

	userRepo := mongoRepo.NewUserRepository(collections.Users)
	productRepo := mongoRepo.NewProductRepository(collections.Products)
	orderRepo := mongoRepo.NewOrderRepository(collections.Orders)
	stockLogRepo := mongoRepo.NewStockLogRepository(collections.StockLogs)

	authUC := usecase.NewAuthUsecase(userRepo, jwtSvc)
	userUC := usecase.NewUserUsecase(userRepo)
	productUC := usecase.NewProductUsecase(productRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, productRepo, stockLogRepo)

	router := httpHandler.NewRouter(httpHandler.Dependencies{
		AuthUsecase:    authUC,
		UserUsecase:    userUC,
		ProductUsecase: productUC,
		OrderUsecase:   orderUC,
		JWTService:     jwtSvc,
		Logger:         appLogger,
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		appLogger.Info("server started", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quitCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-quitCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("graceful shutdown failed", "error", err)
	}
	appLogger.Info("server stopped")
}
