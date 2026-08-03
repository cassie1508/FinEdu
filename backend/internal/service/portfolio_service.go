package service

import (
	"context"
	"errors"
	"strings"

	"finedu-backend/internal/models"
	"finedu-backend/internal/repository"
)

var (
	ErrForbidden = errors.New("access denied")
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
)

type PortfolioService struct {
	repo      repository.PortfolioRepo
	marketSvc QuoteService
	aiRiskSvc AIRiskAnalyzer
}

func NewPortfolioService(repo repository.PortfolioRepo, marketSvc QuoteService, aiRiskSvc AIRiskAnalyzer) *PortfolioService {
	return &PortfolioService{repo: repo, marketSvc: marketSvc, aiRiskSvc: aiRiskSvc}
}

// CreatePortfolio creates a new portfolio for the user.
func (s *PortfolioService) CreatePortfolio(ctx context.Context, userID, name, description string) (*models.Portfolio, error) {
	p, err := s.repo.CreatePortfolio(ctx, userID, name, description)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return p, nil
}

// ListPortfolios returns all portfolios for the authenticated user.
func (s *PortfolioService) ListPortfolios(ctx context.Context, userID string) ([]models.PortfolioListItem, error) {
	return s.repo.ListPortfoliosByUser(ctx, userID)
}

// GetPortfolioDetail returns a portfolio with holdings enriched with live market data.
func (s *PortfolioService) GetPortfolioDetail(ctx context.Context, userID, portfolioID string) (*models.PortfolioDetail, error) {
	portfolio, err := s.verifyOwnership(ctx, userID, portfolioID)
	if err != nil {
		return nil, err
	}

	holdings, err := s.repo.GetHoldingsByPortfolio(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	// Collect unique symbols for batch quote fetch
	symbols := make([]string, 0, len(holdings))
	for _, h := range holdings {
		symbols = append(symbols, h.Symbol)
	}

	quotes, _ := s.marketSvc.GetBatchQuotes(ctx, symbols)

	var totalValue, totalCost float64
	details := make([]models.HoldingDetail, 0, len(holdings))

	for _, h := range holdings {
		cost := h.Shares * h.AverageCost
		totalCost += cost

		var currentPrice float64
		if q, ok := quotes[h.Symbol]; ok {
			currentPrice = q.CurrentPrice
		}

		currentValue := h.Shares * currentPrice
		totalValue += currentValue
		gainLoss := currentValue - cost
		var gainLossPercent float64
		if cost > 0 {
			gainLossPercent = (gainLoss / cost) * 100
		}

		details = append(details, models.HoldingDetail{
			ID:                        h.ID,
			PortfolioID:               h.PortfolioID,
			Symbol:                    h.Symbol,
			Shares:                    h.Shares,
			AverageCost:               h.AverageCost,
			CurrentPrice:              currentPrice,
			CurrentValue:              currentValue,
			UnrealizedGainLoss:        gainLoss,
			UnrealizedGainLossPercent: gainLossPercent,
		})
	}

	// Calculate allocation percentages
	if totalValue > 0 {
		for i := range details {
			details[i].AllocationPercent = (details[i].CurrentValue / totalValue) * 100
		}
	}

	var totalGainLossPercent float64
	if totalCost > 0 {
		totalGainLossPercent = ((totalValue - totalCost) / totalCost) * 100
	}

	return &models.PortfolioDetail{
		ID:                   portfolio.ID,
		Name:                 portfolio.Name,
		Description:          portfolio.Description,
		Holdings:             details,
		TotalValue:           totalValue,
		TotalCost:            totalCost,
		TotalGainLoss:        totalValue - totalCost,
		TotalGainLossPercent: totalGainLossPercent,
	}, nil
}

// DeletePortfolio removes a portfolio after verifying ownership.
func (s *PortfolioService) DeletePortfolio(ctx context.Context, userID, portfolioID string) error {
	if _, err := s.verifyOwnership(ctx, userID, portfolioID); err != nil {
		return err
	}
	err := s.repo.DeletePortfolio(ctx, portfolioID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// AddHolding adds a new holding to a portfolio.
func (s *PortfolioService) AddHolding(ctx context.Context, userID, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	if _, err := s.verifyOwnership(ctx, userID, portfolioID); err != nil {
		return nil, err
	}

	h, err := s.repo.AddHolding(ctx, portfolioID, strings.ToUpper(symbol), shares, averageCost)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return h, nil
}

// UpdateHolding modifies an existing holding after verifying ownership.
func (s *PortfolioService) UpdateHolding(ctx context.Context, userID, portfolioID, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	if _, err := s.verifyOwnership(ctx, userID, portfolioID); err != nil {
		return nil, err
	}

	holding, err := s.repo.GetHoldingByID(ctx, holdingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Verify the holding belongs to this portfolio
	if holding.PortfolioID != portfolioID {
		return nil, ErrNotFound
	}

	h, err := s.repo.UpdateHolding(ctx, holdingID, shares, averageCost)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return h, nil
}

// RemoveHolding deletes a holding after verifying ownership.
func (s *PortfolioService) RemoveHolding(ctx context.Context, userID, portfolioID, holdingID string) error {
	if _, err := s.verifyOwnership(ctx, userID, portfolioID); err != nil {
		return err
	}

	holding, err := s.repo.GetHoldingByID(ctx, holdingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if holding.PortfolioID != portfolioID {
		return ErrNotFound
	}

	return s.repo.DeleteHolding(ctx, holdingID)
}

// GetPortfolioSummary returns allocation and P&L data for a portfolio.
func (s *PortfolioService) GetPortfolioSummary(ctx context.Context, userID, portfolioID string) (*models.PortfolioSummary, error) {
	if _, err := s.verifyOwnership(ctx, userID, portfolioID); err != nil {
		return nil, err
	}

	holdings, err := s.repo.GetHoldingsByPortfolio(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	symbols := make([]string, 0, len(holdings))
	for _, h := range holdings {
		symbols = append(symbols, h.Symbol)
	}

	quotes, _ := s.marketSvc.GetBatchQuotes(ctx, symbols)

	var totalValue, totalCost float64
	values := make(map[string]float64, len(holdings))

	for _, h := range holdings {
		cost := h.Shares * h.AverageCost
		totalCost += cost

		var price float64
		if q, ok := quotes[h.Symbol]; ok {
			price = q.CurrentPrice
		}
		val := h.Shares * price
		totalValue += val
		values[h.Symbol] = val
	}

	allocations := make([]models.AllocationEntry, 0, len(holdings))
	for _, h := range holdings {
		var pct float64
		if totalValue > 0 {
			pct = (values[h.Symbol] / totalValue) * 100
		}
		allocations = append(allocations, models.AllocationEntry{
			Symbol:  h.Symbol,
			Percent: pct,
		})
	}

	var totalGainLossPercent float64
	if totalCost > 0 {
		totalGainLossPercent = ((totalValue - totalCost) / totalCost) * 100
	}

	return &models.PortfolioSummary{
		TotalValue:           totalValue,
		TotalCost:            totalCost,
		TotalGainLoss:        totalValue - totalCost,
		TotalGainLossPercent: totalGainLossPercent,
		Allocations:          allocations,
	}, nil
}

// GetPortfolioRisk returns AI-powered risk analysis for a portfolio.
func (s *PortfolioService) GetPortfolioRisk(ctx context.Context, userID, portfolioID string) (*models.PortfolioRisk, error) {
	if _, err := s.verifyOwnership(ctx, userID, portfolioID); err != nil {
		return nil, err
	}

	holdings, err := s.repo.GetHoldingsByPortfolio(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	if len(holdings) == 0 {
		return &models.PortfolioRisk{
			SectorConcentration: []models.SectorEntry{},
			Recommendations:     []string{},
			Message:             "Add holdings to receive risk analysis",
		}, nil
	}

	symbols := make([]string, len(holdings))
	for i, h := range holdings {
		symbols[i] = h.Symbol
	}
	quotes, _ := s.marketSvc.GetBatchQuotes(ctx, symbols)

	var totalValue float64
	values := make([]float64, len(holdings))
	for i, h := range holdings {
		price := 0.0
		if q, ok := quotes[h.Symbol]; ok {
			price = q.CurrentPrice
		}
		values[i] = h.Shares * price
		totalValue += values[i]
	}

	analysisInput := make([]HoldingForAnalysis, len(holdings))
	for i, h := range holdings {
		pct := 0.0
		if totalValue > 0 {
			pct = (values[i] / totalValue) * 100
		}
		price := 0.0
		if q, ok := quotes[h.Symbol]; ok {
			price = q.CurrentPrice
		}
		analysisInput[i] = HoldingForAnalysis{
			Symbol:            h.Symbol,
			Shares:            h.Shares,
			AllocationPercent: pct,
			CurrentPrice:      price,
		}
	}

	risk, err := s.aiRiskSvc.AnalyzeRisk(ctx, analysisInput)
	if err != nil {
		return &models.PortfolioRisk{
			SectorConcentration: []models.SectorEntry{},
			Recommendations:     []string{},
			Message:             "AI analysis temporarily unavailable",
		}, nil
	}

	return risk, nil
}

// verifyOwnership checks that the portfolio exists and belongs to the user.
func (s *PortfolioService) verifyOwnership(ctx context.Context, userID, portfolioID string) (*models.Portfolio, error) {
	portfolio, err := s.repo.GetPortfolioByID(ctx, portfolioID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if portfolio.UserID != userID {
		return nil, ErrForbidden
	}
	return portfolio, nil
}
