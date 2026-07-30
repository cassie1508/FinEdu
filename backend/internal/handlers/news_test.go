package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"finedu-backend/internal/models"
	"finedu-backend/internal/service"
)

func TestRangeToParams(t *testing.T) {
	tests := []struct {
		rangeParam     string
		wantResolution string
	}{
		{"1W", "D"},
		{"1M", "D"},
		{"6M", "W"},
		{"1Y", "W"},
		{"5Y", "W"},
	}

	for _, tt := range tests {
		resolution, from, to, ok := rangeToParams(tt.rangeParam)
		if !ok {
			t.Errorf("rangeToParams(%q) returned ok=false, want true", tt.rangeParam)
			continue
		}
		if resolution != tt.wantResolution {
			t.Errorf("rangeToParams(%q) resolution = %q, want %q", tt.rangeParam, resolution, tt.wantResolution)
		}
		if !from.Before(to) {
			t.Errorf("rangeToParams(%q) from (%v) is not before to (%v)", tt.rangeParam, from, to)
		}
		if to.After(time.Now().Add(time.Minute)) {
			t.Errorf("rangeToParams(%q) to (%v) is unexpectedly in the future", tt.rangeParam, to)
		}
	}
}

func TestRangeToParams_InvalidRange(t *testing.T) {
	// 1D was intentionally dropped: with only daily-granularity data
	// available, a true 1-day window is a single point and a widened one
	// just duplicates 1W.
	if _, _, _, ok := rangeToParams("1D"); ok {
		t.Error(`rangeToParams("1D") returned ok=true, want false (1D was removed)`)
	}
	if _, _, _, ok := rangeToParams("bogus"); ok {
		t.Error(`rangeToParams("bogus") returned ok=true, want false`)
	}
}

// stubMarketData returns a canned response for every call, letting handler
// tests exercise each response-mapping branch without a real provider.
type stubMarketData struct {
	candles []models.CandlePoint
	err     error
}

func (s *stubMarketData) GetCandles(ctx context.Context, symbol, resolution string, from, to time.Time) ([]models.CandlePoint, error) {
	return s.candles, s.err
}

func newTestRouter(marketData service.MarketDataService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewChartsHandler(marketData)
	r.GET("/companies/:symbol/prices", h.GetPriceHistory)
	return r
}

func TestGetPriceHistory_Success(t *testing.T) {
	stub := &stubMarketData{candles: []models.CandlePoint{{Time: 1, Close: 100}}}
	r := newTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/companies/aapl/prices?range=1M", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetPriceHistory_InvalidRange(t *testing.T) {
	r := newTestRouter(&stubMarketData{})

	req := httptest.NewRequest(http.MethodGet, "/companies/AAPL/prices?range=bogus", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetPriceHistory_DefaultsRangeTo1M(t *testing.T) {
	stub := &stubMarketData{candles: []models.CandlePoint{{Time: 1, Close: 100}}}
	r := newTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/companies/AAPL/prices", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetPriceHistory_UpstreamDeclinedMapsTo503(t *testing.T) {
	stub := &stubMarketData{err: service.ErrUpstreamUnavailable}
	r := newTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/companies/AAPL/prices?range=1M", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", w.Code, http.StatusServiceUnavailable)
	}
}

func TestGetPriceHistory_TransportErrorMapsTo502(t *testing.T) {
	stub := &stubMarketData{err: errors.New("connection refused")}
	r := newTestRouter(stub)

	req := httptest.NewRequest(http.MethodGet, "/companies/AAPL/prices?range=1M", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
}

// --- NewsHandler ---

// stubNewsProvider returns canned articles/error for every call, letting
// handler tests exercise each response-mapping branch without a real
// provider.
type stubNewsProvider struct {
	articles []models.NewsArticle
	err      error
}

func (s *stubNewsProvider) GetTickerNews(ctx context.Context, ticker string) ([]models.NewsArticle, error) {
	return s.articles, s.err
}

func (s *stubNewsProvider) GetGeneralNews(ctx context.Context, category string) ([]models.NewsArticle, error) {
	return s.articles, s.err
}

// stubSummaryService returns a canned summary/error for every call.
type stubSummaryService struct {
	summary models.NewsDailySummary
	err     error
}

func (s *stubSummaryService) GetSummary(ctx context.Context, ticker string) (models.NewsDailySummary, error) {
	return s.summary, s.err
}

func newTestNewsRouter(articles service.NewsProvider, summaries service.SummaryService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewNewsHandler(articles, summaries)
	r.GET("/news/general", h.GetGeneralNews)
	r.GET("/news/:ticker", h.GetTickerNews)
	r.GET("/news/:ticker/summary", h.GetTickerNewsSummary)
	return r
}

func TestGetTickerNews_Success(t *testing.T) {
	articles := &stubNewsProvider{articles: []models.NewsArticle{
		{ID: "1", Headline: "AAPL beats earnings", PublishedAt: time.Now()},
	}}
	r := newTestNewsRouter(articles, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/aapl", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}

	var got []models.NewsArticle
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d articles, want 1", len(got))
	}
}

func TestGetTickerNews_InvalidLimit(t *testing.T) {
	r := newTestNewsRouter(&stubNewsProvider{}, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/AAPL?limit=0", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetTickerNews_InvalidSince(t *testing.T) {
	r := newTestNewsRouter(&stubNewsProvider{}, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/AAPL?since=not-a-date", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetTickerNews_SinceFiltersAndRespectsLimit(t *testing.T) {
	now := time.Now()
	articles := &stubNewsProvider{articles: []models.NewsArticle{
		{ID: "new", PublishedAt: now},
		{ID: "old", PublishedAt: now.AddDate(0, 0, -10)},
	}}
	r := newTestNewsRouter(articles, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/AAPL?since="+now.AddDate(0, 0, -1).Format("2006-01-02"), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var got []models.NewsArticle
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) != 1 || got[0].ID != "new" {
		t.Errorf("got %+v, want only the article newer than `since`", got)
	}
}

func TestGetTickerNews_UpstreamDeclinedMapsTo503(t *testing.T) {
	r := newTestNewsRouter(&stubNewsProvider{err: service.ErrUpstreamUnavailable}, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/AAPL", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", w.Code, http.StatusServiceUnavailable)
	}
}

func TestGetTickerNewsSummary_Success(t *testing.T) {
	summary := &stubSummaryService{summary: models.NewsDailySummary{
		Ticker: "AAPL", Sentiment: "bullish", SourceArticleIDs: []string{"1", "2"},
	}}
	r := newTestNewsRouter(&stubNewsProvider{}, summary)

	req := httptest.NewRequest(http.MethodGet, "/news/AAPL/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetTickerNewsSummary_NoArticlesMapsTo404(t *testing.T) {
	r := newTestNewsRouter(&stubNewsProvider{}, &stubSummaryService{err: service.ErrNoArticles})

	req := httptest.NewRequest(http.MethodGet, "/news/AAPL/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetTickerNewsSummary_UpstreamDeclinedMapsTo503(t *testing.T) {
	r := newTestNewsRouter(&stubNewsProvider{}, &stubSummaryService{err: service.ErrUpstreamUnavailable})

	req := httptest.NewRequest(http.MethodGet, "/news/AAPL/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", w.Code, http.StatusServiceUnavailable)
	}
}

func TestGetGeneralNews_Success(t *testing.T) {
	articles := &stubNewsProvider{articles: []models.NewsArticle{{ID: "1"}}}
	r := newTestNewsRouter(articles, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/general?category=earnings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetGeneralNews_DefaultsCategoryToMarket(t *testing.T) {
	articles := &stubNewsProvider{articles: []models.NewsArticle{{ID: "1"}}}
	r := newTestNewsRouter(articles, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/general", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetGeneralNews_InvalidCategory(t *testing.T) {
	r := newTestNewsRouter(&stubNewsProvider{}, &stubSummaryService{})

	req := httptest.NewRequest(http.MethodGet, "/news/general?category=bogus", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
