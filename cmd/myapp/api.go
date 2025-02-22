package main

import (
	"cofee-shop-mongo/internal/handlers"
	"cofee-shop-mongo/internal/repository"
	"cofee-shop-mongo/internal/service"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type APIServer struct {
	address string
	mux     *http.ServeMux
	db      *mongo.Database
	logger  *slog.Logger
}

func NewAPIServer(address string, mux *http.ServeMux, db *mongo.Database, logger *slog.Logger) *APIServer {
	return &APIServer{
		address: address,
		mux:     mux,
		db:      db,
		logger:  logger,
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
	orderService := service.NewOrderService(orderRepository, menuService, inventoryService)
	orderHandler := handlers.NewOrderHandler(orderService, as.logger)
	orderHandler.RegisterEndpoints(as.mux)

	mWChain := handlers.NewMiddleWareChain(handlers.Recovery, handlers.ContextMW)

	as.logger.Info("starting server", slog.String("address", as.address))
	http.ListenAndServe(as.address, mWChain(as.mux))

}
