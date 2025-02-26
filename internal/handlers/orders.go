package handlers

import (
	"cofee-shop-mongo/internal/auth"
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"log/slog"
	"net/http"
)

type OrderService interface {
	CreateOrder(ctx context.Context, item models.Order) (string, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	GetOrderById(ctx context.Context, OrderId string) (models.Order, error)
	UpdateOrderById(ctx context.Context, OrderId string, item models.Order) error
	DeleteOrderById(ctx context.Context, OrderId string) error
	CloseOrderById(ctx context.Context, OrderId string) error
}

type OrderHandler struct {
	Service OrderService
	Logger  *slog.Logger
}

func NewOrderHandler(orderService OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{orderService, logger}
}

func (h *OrderHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", h.CreateOrder)
	mux.HandleFunc("POST /orders/", h.CreateOrder)

	mux.HandleFunc("GET /orders", h.GetAllOrders)
	mux.HandleFunc("GET /orders/", h.GetAllOrders)

	mux.HandleFunc("GET /orders/{id}", h.GetOrderById)
	mux.HandleFunc("GET /orders/{id}/", h.GetOrderById)

	mux.HandleFunc("PUT /orders/{id}", h.UpdateOrderById)
	mux.HandleFunc("PUT /orders/{id}/", h.UpdateOrderById)

	mux.HandleFunc("DELETE /orders/{id}", h.DeleteOrderById)
	mux.HandleFunc("DELETE /orders/{id}/", h.DeleteOrderById)

	mux.HandleFunc("POST /orders/{id}/close", auth.WithJWTAuth(models.StaffAccess, h.CloseOrderById))
	mux.HandleFunc("POST /orders/{id}/close/", auth.WithJWTAuth(models.StaffAccess, h.CloseOrderById))
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := utils.ParseJSON(r, &order)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	_, err = h.Service.CreateOrder(r.Context(), order)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write([]byte("new item was created successfully"))
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.GetAllOrders(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, items)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := h.Service.GetOrderById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedOrder models.Order

	err := utils.ParseJSON(r, &updatedOrder)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.Service.UpdateOrderById(r.Context(), id, updatedOrder)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedOrder)
}

func (h *OrderHandler) DeleteOrderById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.Service.DeleteOrderById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (h *OrderHandler) CloseOrderById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.Service.CloseOrderById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
