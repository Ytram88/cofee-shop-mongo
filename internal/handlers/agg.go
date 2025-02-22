package handlers

import (
	"cofee-shop-mongo/internal/utils"
	"cofee-shop-mongo/models"
	"context"
	"fmt"
	"net/http"
)

type ReportService interface {
	GetPopularItems(ctx context.Context) ([]models.PopularItem, error)
	GetTotalSales(ctx context.Context) (float64, error)
}

type ReportHandler struct {
	Service ReportService
}

func NewReportHandler(rs ReportService) *ReportHandler {
	return &ReportHandler{rs}
}

func (h *ReportHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /reports/total-sales", h.GetTotalSales)
	mux.HandleFunc("GET /reports/total-sales/", h.GetTotalSales)

	mux.HandleFunc("GET /reports/popular-items", h.GetPopularItems)
	mux.HandleFunc("GET /reports/popular-items/", h.GetPopularItems)
}

func (h *ReportHandler) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	total, err := h.Service.GetTotalSales(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("couldn't get total sales: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]float64{"total_sales": total})
}

func (h *ReportHandler) GetPopularItems(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.Service.GetPopularItems(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("couldn't get popular items: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, popularItems)
}
