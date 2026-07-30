package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"finedu-backend/internal/models"
	"finedu-backend/internal/repository"
)

// mockCompanyRepo là repo giả dùng cho test, không cần database thật.
type mockCompanyRepo struct {
	company   *models.Company
	companies []models.Company
	err       error
}

func (m *mockCompanyRepo) GetBySymbol(_ context.Context, _ string) (*models.Company, error) {
	return m.company, m.err
}

func (m *mockCompanyRepo) Search(_ context.Context, _ string) ([]models.Company, error) {
	return m.companies, m.err
}

var sampleCompany = &models.Company{
	Symbol: "AAPL", Name: "Apple Inc.", Sector: "Technology",
	Industry: "Consumer Electronics", MarketCap: 3e12, Revenue: 3.83e11,
	EPS: 6.1, PERatio: 28.5, DividendYield: 0.5, WeekHigh52: 199.62, WeekLow52: 164.08,
}

func setupRouter(repo CompanyRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCompanyHandler(repo)
	r.GET("/api/v1/companies", h.SearchCompanies)
	r.GET("/api/v1/companies/:symbol", h.GetCompany)
	return r
}

func TestGetCompany(t *testing.T) {
	tests := []struct {
		name       string
		repo       *mockCompanyRepo
		url        string
		wantStatus int
		wantSymbol string
	}{
		{
			name:       "found returns 200",
			repo:       &mockCompanyRepo{company: sampleCompany},
			url:        "/api/v1/companies/AAPL",
			wantStatus: http.StatusOK,
			wantSymbol: "AAPL",
		},
		{
			name:       "not found returns 404",
			repo:       &mockCompanyRepo{err: repository.ErrNotFound},
			url:        "/api/v1/companies/ZZZZ",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "repo error returns 500",
			repo:       &mockCompanyRepo{err: errorsNew("db down")},
			url:        "/api/v1/companies/AAPL",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.repo)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantSymbol != "" {
				var got models.Company
				if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
					t.Fatalf("cannot parse body: %v", err)
				}
				if got.Symbol != tt.wantSymbol {
					t.Errorf("symbol = %q, want %q", got.Symbol, tt.wantSymbol)
				}
			}
		})
	}
}

func TestSearchCompanies(t *testing.T) {
	tests := []struct {
		name       string
		repo       *mockCompanyRepo
		url        string
		wantStatus int
		wantCount  int
	}{
		{
			name:       "empty query returns empty list",
			repo:       &mockCompanyRepo{},
			url:        "/api/v1/companies",
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:       "with results",
			repo:       &mockCompanyRepo{companies: []models.Company{*sampleCompany}},
			url:        "/api/v1/companies?q=apple",
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:       "nil results normalized to empty array",
			repo:       &mockCompanyRepo{companies: nil},
			url:        "/api/v1/companies?q=xyz",
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:       "repo error returns 500",
			repo:       &mockCompanyRepo{err: errorsNew("db down")},
			url:        "/api/v1/companies?q=apple",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.repo)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusOK {
				var body struct {
					Companies []models.Company `json:"companies"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatalf("cannot parse body: %v", err)
				}
				if len(body.Companies) != tt.wantCount {
					t.Errorf("count = %d, want %d", len(body.Companies), tt.wantCount)
				}
			}
		})
	}
}

// errorsNew tạo error đơn giản cho test.
func errorsNew(msg string) error {
	return &simpleError{msg}
}

type simpleError struct{ s string }

func (e *simpleError) Error() string { return e.s }