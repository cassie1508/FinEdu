package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"finedu-backend/internal/models"
	"finedu-backend/internal/repository"
	"finedu-backend/internal/service"
)

// helpers

func testPortfolio(id, userID string) *models.Portfolio {
	return &models.Portfolio{
		ID:          id,
		UserID:      userID,
		Name:        "Test Portfolio",
		Description: "desc",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func testHolding(id, portfolioID, symbol string) models.PortfolioHolding {
	return models.PortfolioHolding{
		ID:          id,
		PortfolioID: portfolioID,
		Symbol:      symbol,
		Shares:      10,
		AverageCost: 100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func noQuotes(_ context.Context, _ []string) (map[string]*models.StockQuote, error) {
	return map[string]*models.StockQuote{}, nil
}

func quoteMap(symbol string, price float64) func(context.Context, []string) (map[string]*models.StockQuote, error) {
	return func(_ context.Context, _ []string) (map[string]*models.StockQuote, error) {
		return map[string]*models.StockQuote{
			symbol: {Symbol: symbol, CurrentPrice: price},
		}, nil
	}
}

func noRisk(_ context.Context, _ []service.HoldingForAnalysis) (*models.PortfolioRisk, error) {
	h := 75.0
	d := 60.0
	rl := "Moderate"
	return &models.PortfolioRisk{
		HealthScore:          &h,
		RiskLevel:            &rl,
		SectorConcentration:  []models.SectorEntry{{Sector: "Technology", Percent: 100}},
		DiversificationScore: &d,
		Recommendations:      []string{"Diversify"},
	}, nil
}

func buildSvc(r *mockRepo, m *mockMarketData, ai *mockAIRisk) *service.PortfolioService {
	return service.NewPortfolioService(r, m, ai)
}

// --- CreatePortfolio ---

func TestCreatePortfolio_Success(t *testing.T) {
	p := testPortfolio("p1", "u1")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			CreatePortfolioFn: func(_ context.Context, _, _, _ string) (*models.Portfolio, error) {
				return p, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	got, err := svc.CreatePortfolio(context.Background(), "u1", "Test Portfolio", "desc")
	require.NoError(t, err)
	assert.Equal(t, p.ID, got.ID)
}

func TestCreatePortfolio_Duplicate(t *testing.T) {
	svc := buildSvc(
		&mockRepo{
			CreatePortfolioFn: func(_ context.Context, _, _, _ string) (*models.Portfolio, error) {
				return nil, repository.ErrDuplicate
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.CreatePortfolio(context.Background(), "u1", "Dup", "")
	assert.ErrorIs(t, err, service.ErrDuplicate)
}

// --- ListPortfolios ---

func TestListPortfolios(t *testing.T) {
	items := []models.PortfolioListItem{{ID: "p1", Name: "P1"}}
	svc := buildSvc(
		&mockRepo{
			ListPortfoliosByUserFn: func(_ context.Context, _ string) ([]models.PortfolioListItem, error) {
				return items, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	got, err := svc.ListPortfolios(context.Background(), "u1")
	require.NoError(t, err)
	assert.Len(t, got, 1)
}

// --- GetPortfolioDetail ---

func TestGetPortfolioDetail_Success(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "p1", "AAPL")

	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingsByPortfolioFn: func(_ context.Context, _ string) ([]models.PortfolioHolding, error) {
				return []models.PortfolioHolding{h}, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: quoteMap("AAPL", 150.0)},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	detail, err := svc.GetPortfolioDetail(context.Background(), "u1", "p1")
	require.NoError(t, err)
	assert.Equal(t, "p1", detail.ID)
	require.Len(t, detail.Holdings, 1)
	assert.Equal(t, 150.0, detail.Holdings[0].CurrentPrice)
	assert.Equal(t, 1500.0, detail.Holdings[0].CurrentValue)        // 10 shares * $150
	assert.Equal(t, 500.0, detail.Holdings[0].UnrealizedGainLoss)   // $1500 - $1000 cost
}

func TestGetPortfolioDetail_GainLossCalculation(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := models.PortfolioHolding{
		ID: "h1", PortfolioID: "p1", Symbol: "AAPL",
		Shares: 10, AverageCost: 100, // cost = $1000
	}

	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingsByPortfolioFn: func(_ context.Context, _ string) ([]models.PortfolioHolding, error) {
				return []models.PortfolioHolding{h}, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: quoteMap("AAPL", 150.0)},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	detail, err := svc.GetPortfolioDetail(context.Background(), "u1", "p1")
	require.NoError(t, err)
	// 10 shares * $150 = $1500 value, cost = $1000, gain = $500
	assert.Equal(t, 500.0, detail.Holdings[0].UnrealizedGainLoss)
	assert.InDelta(t, 50.0, detail.Holdings[0].UnrealizedGainLossPercent, 0.001)
	assert.Equal(t, 100.0, detail.Holdings[0].AllocationPercent) // only one holding
}

func TestGetPortfolioDetail_NotFound(t *testing.T) {
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return nil, repository.ErrNotFound
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.GetPortfolioDetail(context.Background(), "u1", "bad-id")
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestGetPortfolioDetail_Forbidden(t *testing.T) {
	p := testPortfolio("p1", "other-user")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.GetPortfolioDetail(context.Background(), "u1", "p1")
	assert.ErrorIs(t, err, service.ErrForbidden)
}

// --- DeletePortfolio ---

func TestDeletePortfolio_Success(t *testing.T) {
	p := testPortfolio("p1", "u1")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			DeletePortfolioFn: func(_ context.Context, _ string) error {
				return nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	err := svc.DeletePortfolio(context.Background(), "u1", "p1")
	require.NoError(t, err)
}

func TestDeletePortfolio_Forbidden(t *testing.T) {
	p := testPortfolio("p1", "other-user")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	err := svc.DeletePortfolio(context.Background(), "u1", "p1")
	assert.ErrorIs(t, err, service.ErrForbidden)
}

// --- AddHolding ---

func TestAddHolding_Success(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "p1", "AAPL")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			AddHoldingFn: func(_ context.Context, _, _ string, _, _ float64) (*models.PortfolioHolding, error) {
				return &h, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	got, err := svc.AddHolding(context.Background(), "u1", "p1", "aapl", 10, 100)
	require.NoError(t, err)
	assert.Equal(t, "AAPL", got.Symbol)
}

func TestAddHolding_Duplicate(t *testing.T) {
	p := testPortfolio("p1", "u1")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			AddHoldingFn: func(_ context.Context, _, _ string, _, _ float64) (*models.PortfolioHolding, error) {
				return nil, repository.ErrDuplicate
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.AddHolding(context.Background(), "u1", "p1", "AAPL", 10, 100)
	assert.ErrorIs(t, err, service.ErrDuplicate)
}

// --- UpdateHolding ---

func TestUpdateHolding_Success(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "p1", "AAPL")
	updated := h
	updated.Shares = 20
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingByIDFn: func(_ context.Context, _ string) (*models.PortfolioHolding, error) {
				return &h, nil
			},
			UpdateHoldingFn: func(_ context.Context, _ string, _, _ float64) (*models.PortfolioHolding, error) {
				return &updated, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	got, err := svc.UpdateHolding(context.Background(), "u1", "p1", "h1", 20, 100)
	require.NoError(t, err)
	assert.Equal(t, 20.0, got.Shares)
}

func TestUpdateHolding_HoldingNotInPortfolio(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "other-portfolio", "AAPL")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingByIDFn: func(_ context.Context, _ string) (*models.PortfolioHolding, error) {
				return &h, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.UpdateHolding(context.Background(), "u1", "p1", "h1", 5, 100)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestUpdateHolding_HoldingNotFound(t *testing.T) {
	p := testPortfolio("p1", "u1")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingByIDFn: func(_ context.Context, _ string) (*models.PortfolioHolding, error) {
				return nil, repository.ErrNotFound
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.UpdateHolding(context.Background(), "u1", "p1", "bad-holding", 10, 100)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestUpdateHolding_UpdateRepoError(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "p1", "AAPL")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingByIDFn: func(_ context.Context, _ string) (*models.PortfolioHolding, error) {
				return &h, nil
			},
			UpdateHoldingFn: func(_ context.Context, _ string, _, _ float64) (*models.PortfolioHolding, error) {
				return nil, repository.ErrNotFound
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.UpdateHolding(context.Background(), "u1", "p1", "h1", 5, 50)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

// --- RemoveHolding ---

func TestRemoveHolding_Success(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "p1", "AAPL")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingByIDFn: func(_ context.Context, _ string) (*models.PortfolioHolding, error) {
				return &h, nil
			},
			DeleteHoldingFn: func(_ context.Context, _ string) error {
				return nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	err := svc.RemoveHolding(context.Background(), "u1", "p1", "h1")
	require.NoError(t, err)
}

func TestRemoveHolding_NotFound(t *testing.T) {
	p := testPortfolio("p1", "u1")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingByIDFn: func(_ context.Context, _ string) (*models.PortfolioHolding, error) {
				return nil, repository.ErrNotFound
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	err := svc.RemoveHolding(context.Background(), "u1", "p1", "bad-holding")
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestRemoveHolding_HoldingNotInPortfolio(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "other-portfolio", "AAPL")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingByIDFn: func(_ context.Context, _ string) (*models.PortfolioHolding, error) {
				return &h, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	err := svc.RemoveHolding(context.Background(), "u1", "p1", "h1")
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestCreatePortfolio_RepoError(t *testing.T) {
	svc := buildSvc(
		&mockRepo{
			CreatePortfolioFn: func(_ context.Context, _, _, _ string) (*models.Portfolio, error) {
				return nil, errors.New("db error")
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.CreatePortfolio(context.Background(), "u1", "P1", "")
	assert.Error(t, err)
}

// --- GetPortfolioSummary ---

func TestGetPortfolioSummary_Success(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := models.PortfolioHolding{ID: "h1", PortfolioID: "p1", Symbol: "AAPL", Shares: 10, AverageCost: 100}
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingsByPortfolioFn: func(_ context.Context, _ string) ([]models.PortfolioHolding, error) {
				return []models.PortfolioHolding{h}, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: quoteMap("AAPL", 200.0)},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	summary, err := svc.GetPortfolioSummary(context.Background(), "u1", "p1")
	require.NoError(t, err)
	assert.Equal(t, 2000.0, summary.TotalValue)
	assert.Equal(t, 1000.0, summary.TotalCost)
	assert.Equal(t, 1000.0, summary.TotalGainLoss)
	require.Len(t, summary.Allocations, 1)
	assert.Equal(t, 100.0, summary.Allocations[0].Percent)
}

// --- GetPortfolioRisk ---

func TestGetPortfolioRisk_EmptyHoldings(t *testing.T) {
	p := testPortfolio("p1", "u1")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingsByPortfolioFn: func(_ context.Context, _ string) ([]models.PortfolioHolding, error) {
				return []models.PortfolioHolding{}, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	risk, err := svc.GetPortfolioRisk(context.Background(), "u1", "p1")
	require.NoError(t, err)
	assert.Nil(t, risk.HealthScore)
	assert.Equal(t, "Add holdings to receive risk analysis", risk.Message)
}

func TestGetPortfolioRisk_AISuccess(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "p1", "AAPL")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingsByPortfolioFn: func(_ context.Context, _ string) ([]models.PortfolioHolding, error) {
				return []models.PortfolioHolding{h}, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: quoteMap("AAPL", 150.0)},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	risk, err := svc.GetPortfolioRisk(context.Background(), "u1", "p1")
	require.NoError(t, err)
	require.NotNil(t, risk.HealthScore)
	assert.Equal(t, 75.0, *risk.HealthScore)
	require.NotNil(t, risk.RiskLevel)
	assert.Equal(t, "Moderate", *risk.RiskLevel)
}

func TestGetPortfolioRisk_AIFailure_GracefulDegradation(t *testing.T) {
	p := testPortfolio("p1", "u1")
	h := testHolding("h1", "p1", "AAPL")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
			GetHoldingsByPortfolioFn: func(_ context.Context, _ string) ([]models.PortfolioHolding, error) {
				return []models.PortfolioHolding{h}, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: quoteMap("AAPL", 150.0)},
		&mockAIRisk{
			AnalyzeRiskFn: func(_ context.Context, _ []service.HoldingForAnalysis) (*models.PortfolioRisk, error) {
				return nil, errors.New("openai down")
			},
		},
	)

	risk, err := svc.GetPortfolioRisk(context.Background(), "u1", "p1")
	require.NoError(t, err) // graceful — no error returned to caller
	assert.Nil(t, risk.HealthScore)
	assert.Equal(t, "AI analysis temporarily unavailable", risk.Message)
}

func TestGetPortfolioRisk_NotOwner(t *testing.T) {
	p := testPortfolio("p1", "other-user")
	svc := buildSvc(
		&mockRepo{
			GetPortfolioByIDFn: func(_ context.Context, _ string) (*models.Portfolio, error) {
				return p, nil
			},
		},
		&mockMarketData{GetBatchQuotesFn: noQuotes},
		&mockAIRisk{AnalyzeRiskFn: noRisk},
	)

	_, err := svc.GetPortfolioRisk(context.Background(), "u1", "p1")
	assert.ErrorIs(t, err, service.ErrForbidden)
}
