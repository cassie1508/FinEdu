package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"finedu-backend/internal/handlers"
	"finedu-backend/internal/models"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter(svc *mockService) *gin.Engine {
	r := gin.New()
	// Inject a test userID the same way the auth middleware would
	r.Use(func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	})

	h := handlers.NewPortfolioHandler(svc)
	r.POST("/portfolios", h.CreatePortfolio)
	r.GET("/portfolios", h.ListPortfolios)
	r.GET("/portfolios/:portfolioId", h.GetPortfolioDetail)
	r.DELETE("/portfolios/:portfolioId", h.DeletePortfolio)
	r.POST("/portfolios/:portfolioId/holdings", h.AddHolding)
	r.PUT("/portfolios/:portfolioId/holdings/:holdingId", h.UpdateHolding)
	r.DELETE("/portfolios/:portfolioId/holdings/:holdingId", h.RemoveHolding)
	r.GET("/portfolios/:portfolioId/summary", h.GetPortfolioSummary)
	r.GET("/portfolios/:portfolioId/risk", h.GetPortfolioRisk)
	return r
}

func doRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func testPortfolioModel() *models.Portfolio {
	return &models.Portfolio{
		ID:          "p1",
		UserID:      "test-user-id",
		Name:        "Test Portfolio",
		Description: "A test portfolio",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// --- CreatePortfolio ---

func TestCreatePortfolio_Handler_Success(t *testing.T) {
	svc := &mockService{
		CreatePortfolioFn: func(_ context.Context, _, _, _ string) (*models.Portfolio, error) {
			return testPortfolioModel(), nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodPost, "/portfolios", map[string]string{
		"name": "Test Portfolio", "description": "A test portfolio",
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp models.Portfolio
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "p1", resp.ID)
}

func TestCreatePortfolio_Handler_BadJSON(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/portfolios", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePortfolio_Handler_MissingName(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodPost, "/portfolios", map[string]string{"description": "no name"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePortfolio_Handler_Conflict(t *testing.T) {
	svc := &mockService{
		CreatePortfolioFn: func(_ context.Context, _, _, _ string) (*models.Portfolio, error) {
			return nil, errDuplicate
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodPost, "/portfolios", map[string]string{"name": "Dup"})
	assert.Equal(t, http.StatusConflict, w.Code)
}

// --- ListPortfolios ---

func TestListPortfolios_Handler_Empty(t *testing.T) {
	svc := &mockService{
		ListPortfoliosFn: func(_ context.Context, _ string) ([]models.PortfolioListItem, error) {
			return nil, nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]models.PortfolioListItem
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Empty(t, resp["portfolios"])
}

func TestListPortfolios_Handler_InternalError(t *testing.T) {
	svc := &mockService{
		ListPortfoliosFn: func(_ context.Context, _ string) ([]models.PortfolioListItem, error) {
			return nil, errInternal
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios", nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListPortfolios_Handler_WithItems(t *testing.T) {
	svc := &mockService{
		ListPortfoliosFn: func(_ context.Context, _ string) ([]models.PortfolioListItem, error) {
			return []models.PortfolioListItem{{ID: "p1", Name: "P1"}}, nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]models.PortfolioListItem
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Len(t, resp["portfolios"], 1)
}

// --- GetPortfolioDetail ---

func TestGetPortfolioDetail_Handler_NotFound(t *testing.T) {
	svc := &mockService{
		GetPortfolioDetailFn: func(_ context.Context, _, _ string) (*models.PortfolioDetail, error) {
			return nil, errNotFound
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios/bad-id", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPortfolioDetail_Handler_Forbidden(t *testing.T) {
	svc := &mockService{
		GetPortfolioDetailFn: func(_ context.Context, _, _ string) (*models.PortfolioDetail, error) {
			return nil, errForbidden
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios/p1", nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetPortfolioDetail_Handler_Success(t *testing.T) {
	svc := &mockService{
		GetPortfolioDetailFn: func(_ context.Context, _, _ string) (*models.PortfolioDetail, error) {
			return &models.PortfolioDetail{ID: "p1", Name: "P1", Holdings: []models.HoldingDetail{}}, nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios/p1", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- DeletePortfolio ---

func TestDeletePortfolio_Handler_Success(t *testing.T) {
	svc := &mockService{
		DeletePortfolioFn: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodDelete, "/portfolios/p1", nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeletePortfolio_Handler_Forbidden(t *testing.T) {
	svc := &mockService{
		DeletePortfolioFn: func(_ context.Context, _, _ string) error {
			return errForbidden
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodDelete, "/portfolios/p1", nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// --- AddHolding ---

func TestAddHolding_Handler_Success(t *testing.T) {
	svc := &mockService{
		AddHoldingFn: func(_ context.Context, _, _, _ string, _, _ float64) (*models.PortfolioHolding, error) {
			return &models.PortfolioHolding{ID: "h1", Symbol: "AAPL", Shares: 5, AverageCost: 100}, nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodPost, "/portfolios/p1/holdings", map[string]interface{}{
		"symbol": "AAPL", "shares": 5, "averageCost": 100,
	})
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddHolding_Handler_ValidationError(t *testing.T) {
	svc := &mockService{}
	r := setupRouter(svc)
	// shares is required and must be > 0
	w := doRequest(r, http.MethodPost, "/portfolios/p1/holdings", map[string]interface{}{
		"symbol": "AAPL", "shares": -1, "averageCost": 100,
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- UpdateHolding ---

func TestUpdateHolding_Handler_Success(t *testing.T) {
	svc := &mockService{
		UpdateHoldingFn: func(_ context.Context, _, _, _ string, _, _ float64) (*models.PortfolioHolding, error) {
			return &models.PortfolioHolding{ID: "h1", Symbol: "AAPL", Shares: 10, AverageCost: 120}, nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodPut, "/portfolios/p1/holdings/h1", map[string]interface{}{
		"shares": 10, "averageCost": 120,
	})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateHolding_Handler_NotFound(t *testing.T) {
	svc := &mockService{
		UpdateHoldingFn: func(_ context.Context, _, _, _ string, _, _ float64) (*models.PortfolioHolding, error) {
			return nil, errNotFound
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodPut, "/portfolios/p1/holdings/bad", map[string]interface{}{
		"shares": 10, "averageCost": 100,
	})
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// --- RemoveHolding ---

func TestRemoveHolding_Handler_Success(t *testing.T) {
	svc := &mockService{
		RemoveHoldingFn: func(_ context.Context, _, _, _ string) error {
			return nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodDelete, "/portfolios/p1/holdings/h1", nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

// --- GetPortfolioRisk ---

func TestGetPortfolioRisk_Handler_Success(t *testing.T) {
	h := 75.0
	d := 60.0
	rl := "Moderate"
	svc := &mockService{
		GetPortfolioRiskFn: func(_ context.Context, _, _ string) (*models.PortfolioRisk, error) {
			return &models.PortfolioRisk{
				HealthScore:          &h,
				RiskLevel:            &rl,
				DiversificationScore: &d,
				SectorConcentration:  []models.SectorEntry{{Sector: "Technology", Percent: 100}},
				Recommendations:      []string{"Diversify"},
			}, nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios/p1/risk", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var risk models.PortfolioRisk
	require.NoError(t, json.NewDecoder(w.Body).Decode(&risk))
	require.NotNil(t, risk.HealthScore)
	assert.Equal(t, 75.0, *risk.HealthScore)
}

func TestGetPortfolioRisk_Handler_InternalError(t *testing.T) {
	svc := &mockService{
		GetPortfolioRiskFn: func(_ context.Context, _, _ string) (*models.PortfolioRisk, error) {
			return nil, errInternal
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios/p1/risk", nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// --- GetPortfolioSummary ---

func TestGetPortfolioSummary_Handler_Success(t *testing.T) {
	svc := &mockService{
		GetPortfolioSummaryFn: func(_ context.Context, _, _ string) (*models.PortfolioSummary, error) {
			return &models.PortfolioSummary{TotalValue: 5000, Allocations: []models.AllocationEntry{}}, nil
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios/p1/summary", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPortfolioSummary_Handler_NotFound(t *testing.T) {
	svc := &mockService{
		GetPortfolioSummaryFn: func(_ context.Context, _, _ string) (*models.PortfolioSummary, error) {
			return nil, errNotFound
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodGet, "/portfolios/p1/summary", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRemoveHolding_Handler_NotFound(t *testing.T) {
	svc := &mockService{
		RemoveHoldingFn: func(_ context.Context, _, _, _ string) error {
			return errNotFound
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodDelete, "/portfolios/p1/holdings/bad", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddHolding_Handler_Conflict(t *testing.T) {
	svc := &mockService{
		AddHoldingFn: func(_ context.Context, _, _, _ string, _, _ float64) (*models.PortfolioHolding, error) {
			return nil, errDuplicate
		},
	}
	r := setupRouter(svc)
	w := doRequest(r, http.MethodPost, "/portfolios/p1/holdings", map[string]interface{}{
		"symbol": "AAPL", "shares": 5, "averageCost": 100,
	})
	assert.Equal(t, http.StatusConflict, w.Code)
}
