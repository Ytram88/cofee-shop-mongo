package service

import (
	"cofee-shop-mongo/models"
	"context"
)

type ReportRepository interface {
	GetPopularItems(ctx context.Context) ([]models.PopularItem, error)
	GetTotalSales(ctx context.Context) (float64, error)
}

type ReportService struct {
	repo ReportRepository
}

func NewReportService(repo ReportRepository) *ReportService {
	return &ReportService{repo}
}

func (s *ReportService) GetTotalSales(ctx context.Context) (float64, error) {
	return s.repo.GetTotalSales(ctx)
}

func (s *ReportService) GetPopularItems(ctx context.Context) ([]models.PopularItem, error) {
	return s.repo.GetPopularItems(ctx)
}
