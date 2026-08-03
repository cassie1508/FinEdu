package service_test

import (
	"context"

	"finedu-backend/internal/models"
	"finedu-backend/internal/service"
)

// --- mockRepo implements repository.PortfolioRepo ---

type mockRepo struct {
	CreatePortfolioFn      func(ctx context.Context, userID, name, description string) (*models.Portfolio, error)
	GetPortfolioByIDFn     func(ctx context.Context, portfolioID string) (*models.Portfolio, error)
	ListPortfoliosByUserFn func(ctx context.Context, userID string) ([]models.PortfolioListItem, error)
	DeletePortfolioFn      func(ctx context.Context, portfolioID string) error
	AddHoldingFn           func(ctx context.Context, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error)
	GetHoldingByIDFn       func(ctx context.Context, holdingID string) (*models.PortfolioHolding, error)
	GetHoldingsByPortfolioFn func(ctx context.Context, portfolioID string) ([]models.PortfolioHolding, error)
	UpdateHoldingFn        func(ctx context.Context, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error)
	DeleteHoldingFn        func(ctx context.Context, holdingID string) error
}

func (m *mockRepo) CreatePortfolio(ctx context.Context, userID, name, description string) (*models.Portfolio, error) {
	return m.CreatePortfolioFn(ctx, userID, name, description)
}

func (m *mockRepo) GetPortfolioByID(ctx context.Context, portfolioID string) (*models.Portfolio, error) {
	return m.GetPortfolioByIDFn(ctx, portfolioID)
}

func (m *mockRepo) ListPortfoliosByUser(ctx context.Context, userID string) ([]models.PortfolioListItem, error) {
	return m.ListPortfoliosByUserFn(ctx, userID)
}

func (m *mockRepo) DeletePortfolio(ctx context.Context, portfolioID string) error {
	return m.DeletePortfolioFn(ctx, portfolioID)
}

func (m *mockRepo) AddHolding(ctx context.Context, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	return m.AddHoldingFn(ctx, portfolioID, symbol, shares, averageCost)
}

func (m *mockRepo) GetHoldingByID(ctx context.Context, holdingID string) (*models.PortfolioHolding, error) {
	return m.GetHoldingByIDFn(ctx, holdingID)
}

func (m *mockRepo) GetHoldingsByPortfolio(ctx context.Context, portfolioID string) ([]models.PortfolioHolding, error) {
	return m.GetHoldingsByPortfolioFn(ctx, portfolioID)
}

func (m *mockRepo) UpdateHolding(ctx context.Context, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	return m.UpdateHoldingFn(ctx, holdingID, shares, averageCost)
}

func (m *mockRepo) DeleteHolding(ctx context.Context, holdingID string) error {
	return m.DeleteHoldingFn(ctx, holdingID)
}

// --- mockMarketData implements service.MarketDataService ---

type mockMarketData struct {
	GetQuoteFn       func(ctx context.Context, symbol string) (*models.StockQuote, error)
	GetBatchQuotesFn func(ctx context.Context, symbols []string) (map[string]*models.StockQuote, error)
}

func (m *mockMarketData) GetQuote(ctx context.Context, symbol string) (*models.StockQuote, error) {
	return m.GetQuoteFn(ctx, symbol)
}

func (m *mockMarketData) GetBatchQuotes(ctx context.Context, symbols []string) (map[string]*models.StockQuote, error) {
	return m.GetBatchQuotesFn(ctx, symbols)
}

// --- mockAIRisk implements service.AIRiskAnalyzer ---

type mockAIRisk struct {
	AnalyzeRiskFn func(ctx context.Context, holdings []service.HoldingForAnalysis) (*models.PortfolioRisk, error)
}

func (m *mockAIRisk) AnalyzeRisk(ctx context.Context, holdings []service.HoldingForAnalysis) (*models.PortfolioRisk, error) {
	return m.AnalyzeRiskFn(ctx, holdings)
}
