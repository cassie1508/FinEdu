package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"finedu-backend/internal/models"
)

type CompanyDataService struct {
	apiKey     string
	httpClient *http.Client
}

func NewCompanyDataService(apiKey string) *CompanyDataService {
	return &CompanyDataService{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// --- JSON response shapes ---

type profileResp struct {
	Name                 string  `json:"name"`
	FinnhubIndustry      string  `json:"finnhubIndustry"`
	MarketCapitalization float64 `json:"marketCapitalization"` // đơn vị: triệu USD
}

type metricResp struct {
	Metric struct {
		PE            float64 `json:"peTTM"`
		EPS           float64 `json:"epsTTM"`
		DividendYield float64 `json:"dividendYieldIndicatedAnnual"`
		High52        float64 `json:"52WeekHigh"`
		Low52         float64 `json:"52WeekLow"`
		RevPerShare   float64 `json:"revenuePerShareTTM"`
	} `json:"metric"`
}

// FetchCompany gọi Finnhub, ghép data thành 1 models.Company.
func (s *CompanyDataService) FetchCompany(ctx context.Context, symbol string) (*models.Company, error) {
	var profile profileResp
	if err := s.get(ctx, "/stock/profile2", symbol, &profile); err != nil {
		return nil, err
	}
	if profile.Name == "" {
		return nil, fmt.Errorf("no data for symbol %s", symbol)
	}

	var metric metricResp
	if err := s.get(ctx, "/stock/metric", symbol, &metric); err != nil {
		return nil, err
	}

	// marketCap từ Finnhub là triệu USD → đổi ra USD
	marketCap := profile.MarketCapitalization * 1_000_000

	return &models.Company{
		Symbol:        symbol,
		Name:          profile.Name,
		Sector:        profile.FinnhubIndustry,
		Industry:      profile.FinnhubIndustry,
		MarketCap:     marketCap,
		Revenue:       metric.Metric.RevPerShare, // xem note bên dưới
		EPS:           metric.Metric.EPS,
		PERatio:       metric.Metric.PE,
		DividendYield: metric.Metric.DividendYield,
		WeekHigh52:    metric.Metric.High52,
		WeekLow52:     metric.Metric.Low52,
	}, nil
}

func (s *CompanyDataService) get(ctx context.Context, path, symbol string, out interface{}) error {
	url := fmt.Sprintf("https://finnhub.io/api/v1%s?symbol=%s&token=%s", path, symbol, s.apiKey)
	if path == "/stock/metric" {
		url += "&metric=all"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("finnhub %s status %d", path, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}