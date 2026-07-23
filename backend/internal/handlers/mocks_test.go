package handlers_test

import (
	"context"
	"errors"

	"finedu-backend/internal/models"
	"finedu-backend/internal/service"
)

// ErrNotFound and ErrForbidden re-exported for tests in this package
var (
	errNotFound  = service.ErrNotFound
	errForbidden = service.ErrForbidden
	errDuplicate = service.ErrDuplicate
	errInternal  = errors.New("internal error")
)

// mockService implements service.PortfolioServicer
type mockService struct {
	CreatePortfolioFn    func(ctx context.Context, userID, name, description string) (*models.Portfolio, error)
	ListPortfoliosFn     func(ctx context.Context, userID string) ([]models.PortfolioListItem, error)
	GetPortfolioDetailFn func(ctx context.Context, userID, portfolioID string) (*models.PortfolioDetail, error)
	DeletePortfolioFn    func(ctx context.Context, userID, portfolioID string) error
	AddHoldingFn         func(ctx context.Context, userID, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error)
	UpdateHoldingFn      func(ctx context.Context, userID, portfolioID, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error)
	RemoveHoldingFn      func(ctx context.Context, userID, portfolioID, holdingID string) error
	GetPortfolioSummaryFn func(ctx context.Context, userID, portfolioID string) (*models.PortfolioSummary, error)
	GetPortfolioRiskFn   func(ctx context.Context, userID, portfolioID string) (*models.PortfolioRisk, error)
}

func (m *mockService) CreatePortfolio(ctx context.Context, userID, name, description string) (*models.Portfolio, error) {
	return m.CreatePortfolioFn(ctx, userID, name, description)
}

func (m *mockService) ListPortfolios(ctx context.Context, userID string) ([]models.PortfolioListItem, error) {
	return m.ListPortfoliosFn(ctx, userID)
}

func (m *mockService) GetPortfolioDetail(ctx context.Context, userID, portfolioID string) (*models.PortfolioDetail, error) {
	return m.GetPortfolioDetailFn(ctx, userID, portfolioID)
}

func (m *mockService) DeletePortfolio(ctx context.Context, userID, portfolioID string) error {
	return m.DeletePortfolioFn(ctx, userID, portfolioID)
}

func (m *mockService) AddHolding(ctx context.Context, userID, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	return m.AddHoldingFn(ctx, userID, portfolioID, symbol, shares, averageCost)
}

func (m *mockService) UpdateHolding(ctx context.Context, userID, portfolioID, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	return m.UpdateHoldingFn(ctx, userID, portfolioID, holdingID, shares, averageCost)
}

func (m *mockService) RemoveHolding(ctx context.Context, userID, portfolioID, holdingID string) error {
	return m.RemoveHoldingFn(ctx, userID, portfolioID, holdingID)
}

func (m *mockService) GetPortfolioSummary(ctx context.Context, userID, portfolioID string) (*models.PortfolioSummary, error) {
	return m.GetPortfolioSummaryFn(ctx, userID, portfolioID)
}

func (m *mockService) GetPortfolioRisk(ctx context.Context, userID, portfolioID string) (*models.PortfolioRisk, error) {
	return m.GetPortfolioRiskFn(ctx, userID, portfolioID)
}
