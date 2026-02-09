package services

import (
	"category-api/models"
	"category-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetSalesToday() (*models.ReportToday, error) {
	return s.repo.GetReportToday()
}
