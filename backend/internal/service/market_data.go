package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"finedu-backend/internal/models"
)

// MarketDataService provides historical price candles for a ticker.
type MarketDataService interface {
	GetCandles(ctx context.Context, symbol, resolution string, from, to time.Time) ([]models.CandlePoint, error)
}

const alphaVantageBaseURL = "https://www.alphavantage.co/query"

var ErrUpstreamUnavailable = errors.New("market data provider unavailable")

type AlphaVantageMarketData struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewAlphaVantageMarketData(apiKey string) *AlphaVantageMarketData {
	return &AlphaVantageMarketData{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: alphaVantageBaseURL,
	}
}

// avFunctionParams describes how to call Alpha Vantage and where to find the
// resulting series in the response for a given internal resolution code.
type avFunctionParams struct {
	function   string
	interval   string // only set for TIME_SERIES_INTRADAY
	outputSize string
	seriesKey  string
	timeLayout string
}

func avParamsFor(resolution string) (avFunctionParams, error) {
	switch resolution {
	case "5":
		return avFunctionParams{"TIME_SERIES_INTRADAY", "5min", "full", "Time Series (5min)", "2006-01-02 15:04:05"}, nil
	case "15":
		return avFunctionParams{"TIME_SERIES_INTRADAY", "15min", "full", "Time Series (15min)", "2006-01-02 15:04:05"}, nil
	case "60":
		return avFunctionParams{"TIME_SERIES_INTRADAY", "60min", "full", "Time Series (60min)", "2006-01-02 15:04:05"}, nil
	case "D":
		// outputsize=full is premium-only on this endpoint; compact (latest
		// 100 points, ~4-5 months of trading days) is free and is enough
		// for the ranges that use daily resolution (1D/1W/1M).
		return avFunctionParams{"TIME_SERIES_DAILY", "", "compact", "Time Series (Daily)", "2006-01-02"}, nil
	case "W":
		return avFunctionParams{"TIME_SERIES_WEEKLY", "", "", "Weekly Time Series", "2006-01-02"}, nil
	default:
		return avFunctionParams{}, fmt.Errorf("unsupported resolution %q", resolution)
	}
}

type avBar struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

// GetCandles fetches the full available OHLCV series for symbol at the
// given internal resolution (5, 15, 60, D, W). Alpha Vantage has no
// from/to query params, so it always returns everything it has for that
// resolution (bounded by outputsize) — the from/to window is intentionally
// ignored here and applied later by CachedMarketData, so one upstream fetch
// can serve every UI range that shares a resolution (e.g. 1Y and 5Y both
// use "W").
func (a *AlphaVantageMarketData) GetCandles(ctx context.Context, symbol, resolution string, _, _ time.Time) ([]models.CandlePoint, error) {
	params, err := avParamsFor(resolution)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf(
		"%s?function=%s&symbol=%s&apikey=%s",
		a.baseURL, params.function, url.QueryEscape(symbol), url.QueryEscape(a.apiKey),
	)
	if params.interval != "" {
		reqURL += "&interval=" + params.interval
	}
	if params.outputSize != "" {
		reqURL += "&outputsize=" + params.outputSize
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching candles for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading candle response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alpha vantage API error (status %d): %s", resp.StatusCode, string(body))
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("decoding candle response: %w", err)
	}

	// Alpha Vantage reports errors/rate limits as 200 OK with one of these
	// keys instead of the expected series key.
	for _, errKey := range []string{"Error Message", "Note", "Information"} {
		if msg, ok := raw[errKey]; ok {
			return nil, fmt.Errorf("%w: %s", ErrUpstreamUnavailable, msg)
		}
	}

	seriesRaw, ok := raw[params.seriesKey]
	if !ok {
		return nil, fmt.Errorf("unexpected alpha vantage response: missing %q", params.seriesKey)
	}

	var series map[string]avBar
	if err := json.Unmarshal(seriesRaw, &series); err != nil {
		return nil, fmt.Errorf("decoding candle series: %w", err)
	}

	// Alpha Vantage timestamps have no explicit timezone (they're US/Eastern
	// market time); parsed as UTC here since we only need them for ordering
	// and coarse windowing, not exact market-hour display.
	candles := make([]models.CandlePoint, 0, len(series))
	for ts, bar := range series {
		t, err := time.Parse(params.timeLayout, ts)
		if err != nil {
			continue
		}

		open, _ := strconv.ParseFloat(bar.Open, 64)
		high, _ := strconv.ParseFloat(bar.High, 64)
		low, _ := strconv.ParseFloat(bar.Low, 64)
		closeP, _ := strconv.ParseFloat(bar.Close, 64)
		volume, _ := strconv.ParseFloat(bar.Volume, 64)

		candles = append(candles, models.CandlePoint{
			Time:   t.Unix(),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  closeP,
			Volume: volume,
		})
	}

	sort.Slice(candles, func(i, j int) bool { return candles[i].Time < candles[j].Time })

	return candles, nil
}

// --- Cached Wrapper ---

// Cached by symbol+resolution only (not the exact from/to window), since
// callers always ask for "up to now" — freshness is controlled by TTL
// instead. Intraday resolutions get the shortest TTL; daily next; weekly
// longest, since a week's bar changes even less often than a day's.
//
// The cached value is always the FULL unfiltered series for that
// resolution — GetCandles filters it down to the caller's [from, to] window
// after every cache read (hit or miss). Filtering before caching would mean
// whichever range hits the API first "wins" the cache entry and every other
// range sharing that resolution (e.g. 1Y and 5Y both use "W") would
// silently get served that first range's window instead of its own.
type candleCacheEntry struct {
	candles   []models.CandlePoint
	fetchedAt time.Time
}

type CachedMarketData struct {
	inner       MarketDataService
	intradayTTL time.Duration
	dailyTTL    time.Duration
	weeklyTTL   time.Duration
	mu          sync.RWMutex
	cache       map[string]candleCacheEntry
}

func NewCachedMarketData(inner MarketDataService, intradayTTL, dailyTTL, weeklyTTL time.Duration) *CachedMarketData {
	return &CachedMarketData{
		inner:       inner,
		intradayTTL: intradayTTL,
		dailyTTL:    dailyTTL,
		weeklyTTL:   weeklyTTL,
		cache:       make(map[string]candleCacheEntry),
	}
}

func (c *CachedMarketData) ttlFor(resolution string) time.Duration {
	switch resolution {
	case "W":
		return c.weeklyTTL
	case "D", "M":
		return c.dailyTTL
	default:
		return c.intradayTTL
	}
}

func (c *CachedMarketData) GetCandles(ctx context.Context, symbol, resolution string, from, to time.Time) ([]models.CandlePoint, error) {
	key := symbol + "|" + resolution

	c.mu.RLock()
	entry, ok := c.cache[key]
	c.mu.RUnlock()

	if ok && time.Since(entry.fetchedAt) < c.ttlFor(resolution) {
		return filterByWindow(entry.candles, from, to), nil
	}

	full, err := c.inner.GetCandles(ctx, symbol, resolution, from, to)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[key] = candleCacheEntry{candles: full, fetchedAt: time.Now()}
	c.mu.Unlock()

	return filterByWindow(full, from, to), nil
}

func filterByWindow(candles []models.CandlePoint, from, to time.Time) []models.CandlePoint {
	fromUnix, toUnix := from.Unix(), to.Unix()
	filtered := make([]models.CandlePoint, 0, len(candles))
	for _, candle := range candles {
		if candle.Time < fromUnix || candle.Time > toUnix {
			continue
		}
		filtered = append(filtered, candle)
	}
	return filtered
}
