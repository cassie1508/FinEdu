package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"finedu-backend/internal/models"
	"finedu-backend/internal/service"
)

type finnhubQuoteResponse struct {
	CurrentPrice float64 `json:"c"`
	Change       float64 `json:"d"`
	ChangePercent float64 `json:"dp"`
	High         float64 `json:"h"`
	Low          float64 `json:"l"`
	Open         float64 `json:"o"`
	PrevClose    float64 `json:"pc"`
}

func newFinnhubServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestFinnhubGetQuote_Success(t *testing.T) {
	srv := newFinnhubServer(func(w http.ResponseWriter, r *http.Request) {
		resp := finnhubQuoteResponse{CurrentPrice: 150.0, Change: 2.0, ChangePercent: 1.35, PrevClose: 148.0}
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	md := service.NewFinnhubMarketDataWithBaseURL("test-key", srv.URL)
	quote, err := md.GetQuote(context.Background(), "AAPL")
	require.NoError(t, err)
	assert.Equal(t, "AAPL", quote.Symbol)
	assert.Equal(t, 150.0, quote.CurrentPrice)
	assert.Equal(t, 2.0, quote.Change)
	assert.InDelta(t, 1.35, quote.ChangePercent, 0.001)
}

func TestFinnhubGetQuote_InvalidSymbol(t *testing.T) {
	srv := newFinnhubServer(func(w http.ResponseWriter, r *http.Request) {
		// Finnhub returns zeros for invalid symbols
		json.NewEncoder(w).Encode(finnhubQuoteResponse{})
	})
	defer srv.Close()

	md := service.NewFinnhubMarketDataWithBaseURL("test-key", srv.URL)
	_, err := md.GetQuote(context.Background(), "INVALID")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no quote data available")
}

func TestFinnhubGetQuote_HTTPError(t *testing.T) {
	srv := newFinnhubServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	})
	defer srv.Close()

	md := service.NewFinnhubMarketDataWithBaseURL("test-key", srv.URL)
	_, err := md.GetQuote(context.Background(), "AAPL")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "finnhub API error")
}

func TestFinnhubGetQuote_BadJSON(t *testing.T) {
	srv := newFinnhubServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	})
	defer srv.Close()

	md := service.NewFinnhubMarketDataWithBaseURL("test-key", srv.URL)
	_, err := md.GetQuote(context.Background(), "AAPL")
	require.Error(t, err)
}

func TestFinnhubGetBatchQuotes_PartialFailure(t *testing.T) {
	callCount := 0
	srv := newFinnhubServer(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		sym := r.URL.Query().Get("symbol")
		if sym == "AAPL" {
			json.NewEncoder(w).Encode(finnhubQuoteResponse{CurrentPrice: 150.0, PrevClose: 148.0})
		} else {
			// INVALID returns zeros
			json.NewEncoder(w).Encode(finnhubQuoteResponse{})
		}
	})
	defer srv.Close()

	md := service.NewFinnhubMarketDataWithBaseURL("test-key", srv.URL)
	quotes, err := md.GetBatchQuotes(context.Background(), []string{"AAPL", "INVALID"})
	require.NoError(t, err)
	assert.Contains(t, quotes, "AAPL")
	assert.NotContains(t, quotes, "INVALID") // skipped due to zero price
}

// --- CachedMarketData tests ---

func TestCachedMarketData_HitCache(t *testing.T) {
	callCount := 0
	inner := &mockMarketData{
		GetQuoteFn: func(_ context.Context, symbol string) (*models.StockQuote, error) {
			callCount++
			return &models.StockQuote{Symbol: symbol, CurrentPrice: 100.0}, nil
		},
		GetBatchQuotesFn: func(_ context.Context, symbols []string) (map[string]*models.StockQuote, error) {
			callCount++
			result := make(map[string]*models.StockQuote)
			for _, s := range symbols {
				result[s] = &models.StockQuote{Symbol: s, CurrentPrice: 100.0}
			}
			return result, nil
		},
	}

	cached := service.NewCachedMarketData(inner, 60*time.Second)

	// First call: miss
	q1, err := cached.GetQuote(context.Background(), "AAPL")
	require.NoError(t, err)
	assert.Equal(t, 100.0, q1.CurrentPrice)
	assert.Equal(t, 1, callCount)

	// Second call: hit cache, no new inner call
	q2, err := cached.GetQuote(context.Background(), "AAPL")
	require.NoError(t, err)
	assert.Equal(t, q1.CurrentPrice, q2.CurrentPrice)
	assert.Equal(t, 1, callCount) // still 1
}

func TestCachedMarketData_ExpiredCache(t *testing.T) {
	callCount := 0
	inner := &mockMarketData{
		GetQuoteFn: func(_ context.Context, symbol string) (*models.StockQuote, error) {
			callCount++
			return &models.StockQuote{Symbol: symbol, CurrentPrice: 100.0}, nil
		},
		GetBatchQuotesFn: noQuotes,
	}

	// Very short TTL so cache expires quickly
	cached := service.NewCachedMarketData(inner, 1*time.Millisecond)

	cached.GetQuote(context.Background(), "AAPL")
	time.Sleep(5 * time.Millisecond) // let cache expire
	cached.GetQuote(context.Background(), "AAPL")
	assert.Equal(t, 2, callCount)
}

func TestCachedMarketData_BatchPartialCache(t *testing.T) {
	callCount := 0
	inner := &mockMarketData{
		GetQuoteFn: func(_ context.Context, symbol string) (*models.StockQuote, error) {
			return &models.StockQuote{Symbol: symbol, CurrentPrice: 50.0}, nil
		},
		GetBatchQuotesFn: func(_ context.Context, symbols []string) (map[string]*models.StockQuote, error) {
			callCount++
			result := make(map[string]*models.StockQuote)
			for _, s := range symbols {
				result[s] = &models.StockQuote{Symbol: s, CurrentPrice: 50.0}
			}
			return result, nil
		},
	}

	cached := service.NewCachedMarketData(inner, 60*time.Second)

	// Warm AAPL into cache
	cached.GetBatchQuotes(context.Background(), []string{"AAPL"})
	assert.Equal(t, 1, callCount)

	// Request AAPL + MSFT; AAPL is cached, only MSFT should trigger inner call
	quotes, err := cached.GetBatchQuotes(context.Background(), []string{"AAPL", "MSFT"})
	require.NoError(t, err)
	assert.Contains(t, quotes, "AAPL")
	assert.Contains(t, quotes, "MSFT")
	assert.Equal(t, 2, callCount) // one more inner call for MSFT only
}
