package main

import (
	"cofee-shop-mongo/internal/config"
	"cofee-shop-mongo/internal/handlers"
	"cofee-shop-mongo/internal/handlers/middleware"
	"cofee-shop-mongo/internal/repository"
	"cofee-shop-mongo/internal/service"
	"fmt"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type APIServer struct {
	config *config.Config
	mux    *http.ServeMux
	db     *mongo.Database
	logger *slog.Logger
}

func NewAPIServer(config *config.Config, mux *http.ServeMux, db *mongo.Database, logger *slog.Logger) *APIServer {
	return &APIServer{
		config: config,
		mux:    mux,
		db:     db,
		logger: logger,
	}
}

func (as *APIServer) Run() {
	inventoryRepository := repository.NewInventoryRepository(as.db)
	inventoryService := service.NewInventoryService(inventoryRepository)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, as.logger)
	inventoryHandler.RegisterEndpoints(as.mux)

	menuRepository := repository.NewMenuRepository(as.db)
	menuService := service.NewMenuService(menuRepository)
	menuHandler := handlers.NewMenuHandler(menuService, as.logger)
	menuHandler.RegisterEndpoints(as.mux)

	orderRepository := repository.NewOrderRepository(as.db)
	orderService := service.NewOrderService(orderRepository, menuService, inventoryService) //order needs access to menu and inventory so you need to pass repo or service
	orderHandler := handlers.NewOrderHandler(orderService, as.logger)
	orderHandler.RegisterEndpoints(as.mux)

	userRepository := repository.NewUserRepository(as.db)
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService, as.logger)
	userHandler.RegisterEndpoints(as.mux)

	reportRepository := repository.NewReportRepository(as.db)
	reportService := service.NewReportService(reportRepository)
	reportHandler := handlers.NewReportHandler(reportService)
	reportHandler.RegisterEndpoints(as.mux)

	authService := service.NewAuthService(userRepository, as.config.JWTConfig)
	authHandler := handlers.NewAuthHandler(authService)
	authHandler.RegisterEndpoints(as.mux)

	mWChain := middleware.NewMiddleWareChain(middleware.Recovery, middleware.ContextMW)

	address := fmt.Sprintf("0.0.0.0:%s", as.config.Port)
	as.logger.Info("starting server", slog.String("address", address))
	http.ListenAndServe(address, mWChain(as.mux))

}
