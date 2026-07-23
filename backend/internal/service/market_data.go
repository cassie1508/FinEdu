package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"finedu-backend/internal/models"
)

// MarketDataService provides stock price quotes.
type MarketDataService interface {
	GetQuote(ctx context.Context, symbol string) (*models.StockQuote, error)
	GetBatchQuotes(ctx context.Context, symbols []string) (map[string]*models.StockQuote, error)
}

// --- Finnhub Implementation ---

type FinnhubMarketData struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewFinnhubMarketData(apiKey string) *FinnhubMarketData {
	return &FinnhubMarketData{
		apiKey:  apiKey,
		baseURL: "https://finnhub.io",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewFinnhubMarketDataWithBaseURL creates a FinnhubMarketData pointing at a custom base URL (used in tests).
func NewFinnhubMarketDataWithBaseURL(apiKey, baseURL string) *FinnhubMarketData {
	return &FinnhubMarketData{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type finnhubQuoteResponse struct {
	CurrentPrice  float64 `json:"c"`
	Change        float64 `json:"d"`
	ChangePercent float64 `json:"dp"`
	High          float64 `json:"h"`
	Low           float64 `json:"l"`
	Open          float64 `json:"o"`
	PrevClose     float64 `json:"pc"`
}

func (f *FinnhubMarketData) GetQuote(ctx context.Context, symbol string) (*models.StockQuote, error) {
	url := fmt.Sprintf("%s/api/v1/quote?symbol=%s&token=%s", f.baseURL, symbol, f.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching quote for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("finnhub API error (status %d): %s", resp.StatusCode, string(body))
	}

	var fq finnhubQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&fq); err != nil {
		return nil, fmt.Errorf("decoding quote response: %w", err)
	}

	// Finnhub returns zeros for invalid symbols
	if fq.CurrentPrice == 0 && fq.PrevClose == 0 {
		return nil, fmt.Errorf("no quote data available for symbol %s", symbol)
	}

	return &models.StockQuote{
		Symbol:        symbol,
		CurrentPrice:  fq.CurrentPrice,
		Change:        fq.Change,
		ChangePercent: fq.ChangePercent,
	}, nil
}

func (f *FinnhubMarketData) GetBatchQuotes(ctx context.Context, symbols []string) (map[string]*models.StockQuote, error) {
	quotes := make(map[string]*models.StockQuote, len(symbols))
	for _, sym := range symbols {
		quote, err := f.GetQuote(ctx, sym)
		if err != nil {
			// Skip symbols that fail — return what we can
			continue
		}
		quotes[sym] = quote
	}
	return quotes, nil
}

// --- Cached Wrapper ---

type cacheEntry struct {
	quote     *models.StockQuote
	fetchedAt time.Time
}

type CachedMarketData struct {
	inner MarketDataService
	ttl   time.Duration
	mu    sync.RWMutex
	cache map[string]cacheEntry
}

func NewCachedMarketData(inner MarketDataService, ttl time.Duration) *CachedMarketData {
	return &CachedMarketData{
		inner: inner,
		ttl:   ttl,
		cache: make(map[string]cacheEntry),
	}
}

func (c *CachedMarketData) GetQuote(ctx context.Context, symbol string) (*models.StockQuote, error) {
	c.mu.RLock()
	entry, ok := c.cache[symbol]
	c.mu.RUnlock()

	if ok && time.Since(entry.fetchedAt) < c.ttl {
		return entry.quote, nil
	}

	quote, err := c.inner.GetQuote(ctx, symbol)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[symbol] = cacheEntry{quote: quote, fetchedAt: time.Now()}
	c.mu.Unlock()

	return quote, nil
}

func (c *CachedMarketData) GetBatchQuotes(ctx context.Context, symbols []string) (map[string]*models.StockQuote, error) {
	results := make(map[string]*models.StockQuote, len(symbols))
	var missing []string

	c.mu.RLock()
	for _, sym := range symbols {
		entry, ok := c.cache[sym]
		if ok && time.Since(entry.fetchedAt) < c.ttl {
			results[sym] = entry.quote
		} else {
			missing = append(missing, sym)
		}
	}
	c.mu.RUnlock()

	if len(missing) == 0 {
		return results, nil
	}

	fetched, err := c.inner.GetBatchQuotes(ctx, missing)
	if err != nil {
		return results, err
	}

	c.mu.Lock()
	for sym, quote := range fetched {
		c.cache[sym] = cacheEntry{quote: quote, fetchedAt: time.Now()}
		results[sym] = quote
	}
	c.mu.Unlock()

	return results, nil
}
