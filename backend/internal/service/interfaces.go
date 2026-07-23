package service

import (
	"context"

	"finedu-backend/internal/models"
)

// PortfolioServicer defines the business logic contract used by HTTP handlers.
type PortfolioServicer interface {
	CreatePortfolio(ctx context.Context, userID, name, description string) (*models.Portfolio, error)
	ListPortfolios(ctx context.Context, userID string) ([]models.PortfolioListItem, error)
	GetPortfolioDetail(ctx context.Context, userID, portfolioID string) (*models.PortfolioDetail, error)
	DeletePortfolio(ctx context.Context, userID, portfolioID string) error
	AddHolding(ctx context.Context, userID, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error)
	UpdateHolding(ctx context.Context, userID, portfolioID, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error)
	RemoveHolding(ctx context.Context, userID, portfolioID, holdingID string) error
	GetPortfolioSummary(ctx context.Context, userID, portfolioID string) (*models.PortfolioSummary, error)
	GetPortfolioRisk(ctx context.Context, userID, portfolioID string) (*models.PortfolioRisk, error)
}
