package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"finedu-backend/internal/models"
)

func TestAvParamsFor(t *testing.T) {
	tests := []struct {
		resolution    string
		wantFunction  string
		wantSeriesKey string
	}{
		{"5", "TIME_SERIES_INTRADAY", "Time Series (5min)"},
		{"15", "TIME_SERIES_INTRADAY", "Time Series (15min)"},
		{"60", "TIME_SERIES_INTRADAY", "Time Series (60min)"},
		{"D", "TIME_SERIES_DAILY", "Time Series (Daily)"},
		{"W", "TIME_SERIES_WEEKLY", "Weekly Time Series"},
	}

	for _, tt := range tests {
		params, err := avParamsFor(tt.resolution)
		if err != nil {
			t.Errorf("avParamsFor(%q) returned unexpected error: %v", tt.resolution, err)
			continue
		}
		if params.function != tt.wantFunction {
			t.Errorf("avParamsFor(%q).function = %q, want %q", tt.resolution, params.function, tt.wantFunction)
		}
		if params.seriesKey != tt.wantSeriesKey {
			t.Errorf("avParamsFor(%q).seriesKey = %q, want %q", tt.resolution, params.seriesKey, tt.wantSeriesKey)
		}
	}

	if _, err := avParamsFor("bogus"); err == nil {
		t.Error(`avParamsFor("bogus") expected an error, got nil`)
	}
}

// avParamsFor("D") must request outputsize=compact — Alpha Vantage's free
// tier 403s on outputsize=full for TIME_SERIES_DAILY.
func TestAvParamsFor_DailyUsesCompactOutputSize(t *testing.T) {
	params, err := avParamsFor("D")
	if err != nil {
		t.Fatalf(`avParamsFor("D") returned error: %v`, err)
	}
	if params.outputSize != "compact" {
		t.Errorf("daily outputSize = %q, want %q", params.outputSize, "compact")
	}
}

func newTestServer(t *testing.T, body string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestAlphaVantageMarketData_GetCandles_Daily(t *testing.T) {
	body := `{
		"Meta Data": {"2. Symbol": "AAPL"},
		"Time Series (Daily)": {
			"2026-07-20": {"1. open": "100.0", "2. high": "105.0", "3. low": "99.0", "4. close": "104.0", "5. volume": "1000"},
			"2026-07-21": {"1. open": "104.0", "2. high": "110.0", "3. low": "103.0", "4. close": "108.0", "5. volume": "2000"},
			"2020-01-01": {"1. open": "1.0", "2. high": "1.0", "3. low": "1.0", "4. close": "1.0", "5. volume": "1"}
		}
	}`
	srv := newTestServer(t, body, http.StatusOK)
	defer srv.Close()

	av := NewAlphaVantageMarketData("test-key")
	av.baseURL = srv.URL

	// from/to are intentionally not the window filter here — that's applied
	// later by CachedMarketData — so this returns the whole series
	// regardless of what's passed.
	candles, err := av.GetCandles(context.Background(), "AAPL", "D", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("GetCandles returned error: %v", err)
	}

	if len(candles) != 3 {
		t.Fatalf("got %d candles, want 3 (all bars, unfiltered)", len(candles))
	}

	// Results must be sorted ascending by time — the 2020 bar first.
	for i := 1; i < len(candles); i++ {
		if candles[i-1].Time > candles[i].Time {
			t.Fatal("candles are not sorted ascending by time")
		}
	}
	if candles[0].Open != 1.0 {
		t.Errorf("first (oldest) candle open = %v, want 1.0", candles[0].Open)
	}
	if candles[2].Open != 104.0 || candles[2].Close != 108.0 {
		t.Errorf("last (newest) candle = %+v, values not parsed correctly", candles[2])
	}
}

func TestAlphaVantageMarketData_GetCandles_UpstreamDeclined(t *testing.T) {
	body := `{"Information": "Thank you for using Alpha Vantage! This is a premium endpoint."}`
	srv := newTestServer(t, body, http.StatusOK)
	defer srv.Close()

	av := NewAlphaVantageMarketData("test-key")
	av.baseURL = srv.URL

	_, err := av.GetCandles(context.Background(), "AAPL", "D", time.Now().AddDate(0, -1, 0), time.Now())
	if !errors.Is(err, ErrUpstreamUnavailable) {
		t.Errorf("expected ErrUpstreamUnavailable, got: %v", err)
	}
}

func TestAlphaVantageMarketData_GetCandles_HTTPError(t *testing.T) {
	srv := newTestServer(t, `{"error":"boom"}`, http.StatusInternalServerError)
	defer srv.Close()

	av := NewAlphaVantageMarketData("test-key")
	av.baseURL = srv.URL

	_, err := av.GetCandles(context.Background(), "AAPL", "D", time.Now().AddDate(0, -1, 0), time.Now())
	if err == nil {
		t.Fatal("expected an error for HTTP 500 response, got nil")
	}
	// A transport-level failure is not the same thing as an upstream
	// "declined to serve data" response, and callers treat them differently.
	if errors.Is(err, ErrUpstreamUnavailable) {
		t.Error("HTTP error should not be classified as ErrUpstreamUnavailable")
	}
}

func TestAlphaVantageMarketData_GetCandles_InvalidResolution(t *testing.T) {
	av := NewAlphaVantageMarketData("test-key")
	_, err := av.GetCandles(context.Background(), "AAPL", "bogus", time.Now(), time.Now())
	if err == nil {
		t.Error("expected an error for an unsupported resolution, got nil")
	}
}

// fakeMarketData counts how many times the wrapped service actually fetches,
// so CachedMarketData's caching behavior can be verified without a network
// call.
type fakeMarketData struct {
	calls   int
	candles []models.CandlePoint
	err     error
}

func (f *fakeMarketData) GetCandles(ctx context.Context, symbol, resolution string, from, to time.Time) ([]models.CandlePoint, error) {
	f.calls++
	return f.candles, f.err
}

func TestCachedMarketData_TTLFor(t *testing.T) {
	cached := NewCachedMarketData(&fakeMarketData{}, time.Minute, 6*time.Hour, 24*time.Hour)

	tests := []struct {
		resolution string
		want       time.Duration
	}{
		{"W", 24 * time.Hour},
		{"D", 6 * time.Hour},
		{"M", 6 * time.Hour},
		{"5", time.Minute},
	}

	for _, tt := range tests {
		if got := cached.ttlFor(tt.resolution); got != tt.want {
			t.Errorf("ttlFor(%q) = %v, want %v", tt.resolution, got, tt.want)
		}
	}
}

func TestCachedMarketData_ServesFromCacheWithinTTL(t *testing.T) {
	inner := &fakeMarketData{candles: []models.CandlePoint{{Time: 1, Close: 100}}}
	cached := NewCachedMarketData(inner, time.Minute, time.Hour, time.Hour)

	for range 3 {
		if _, err := cached.GetCandles(context.Background(), "AAPL", "D", time.Now(), time.Now()); err != nil {
			t.Fatalf("GetCandles returned error: %v", err)
		}
	}

	if inner.calls != 1 {
		t.Errorf("inner service was called %d times, want 1 (repeat calls within TTL should be cached)", inner.calls)
	}
}

func TestCachedMarketData_RefetchesAfterTTLExpires(t *testing.T) {
	inner := &fakeMarketData{candles: []models.CandlePoint{{Time: 1, Close: 100}}}
	cached := NewCachedMarketData(inner, time.Minute, time.Hour, time.Hour)

	if _, err := cached.GetCandles(context.Background(), "AAPL", "D", time.Now(), time.Now()); err != nil {
		t.Fatalf("GetCandles returned error: %v", err)
	}

	// Backdate the cache entry past its TTL instead of sleeping, so the test
	// is deterministic.
	cached.mu.Lock()
	entry := cached.cache["AAPL|D"]
	entry.fetchedAt = time.Now().Add(-2 * time.Hour)
	cached.cache["AAPL|D"] = entry
	cached.mu.Unlock()

	if _, err := cached.GetCandles(context.Background(), "AAPL", "D", time.Now(), time.Now()); err != nil {
		t.Fatalf("GetCandles returned error: %v", err)
	}

	if inner.calls != 2 {
		t.Errorf("inner service was called %d times, want 2 (expired entry should trigger a refetch)", inner.calls)
	}
}

func TestFilterByWindow(t *testing.T) {
	base := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	candles := []models.CandlePoint{
		{Time: base.AddDate(-5, 0, 0).Unix()},
		{Time: base.AddDate(-2, 0, 0).Unix()},
		{Time: base.AddDate(0, -6, 0).Unix()},
		{Time: base.Unix()},
	}

	got := filterByWindow(candles, base.AddDate(-1, 0, 0), base)
	if len(got) != 2 {
		t.Fatalf("got %d candles within the 1-year window, want 2", len(got))
	}
}

// Regression test for a real bug: 1Y and 5Y both use resolution "W", so they
// share one cache entry. Filtering had been applied before caching, so
// whichever range hit the API first "won" the cache entry — the other range
// silently got served that first range's window instead of its own.
func TestCachedMarketData_FiltersIndependentlyPerRange(t *testing.T) {
	base := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	full := []models.CandlePoint{
		{Time: base.AddDate(-5, 0, 0).Unix()}, // 5 years back
		{Time: base.AddDate(-2, 0, 0).Unix()}, // 2 years back
		{Time: base.AddDate(0, -6, 0).Unix()}, // 6 months back
	}
	inner := &fakeMarketData{candles: full}
	cached := NewCachedMarketData(inner, time.Minute, time.Hour, time.Hour)

	// 1Y window hits the API first and populates the "AAPL|W" cache entry.
	got1Y, err := cached.GetCandles(context.Background(), "AAPL", "W", base.AddDate(-1, 0, 0), base)
	if err != nil {
		t.Fatalf("1Y GetCandles returned error: %v", err)
	}
	if len(got1Y) != 1 {
		t.Fatalf("1Y window: got %d candles, want 1", len(got1Y))
	}

	// 5Y window shares that cache entry, but must get its own window applied
	// — not the 1Y-filtered result from the call above.
	got5Y, err := cached.GetCandles(context.Background(), "AAPL", "W", base.AddDate(-5, 0, 0), base)
	if err != nil {
		t.Fatalf("5Y GetCandles returned error: %v", err)
	}
	if len(got5Y) != 3 {
		t.Fatalf("5Y window: got %d candles, want 3 (got the 1Y-filtered cache entry instead of its own window)", len(got5Y))
	}

	if inner.calls != 1 {
		t.Errorf("inner service was called %d times, want 1 (both windows should share the one cached fetch)", inner.calls)
	}
}

func TestCachedMarketData_DoesNotCacheErrors(t *testing.T) {
	inner := &fakeMarketData{err: errors.New("boom")}
	cached := NewCachedMarketData(inner, time.Minute, time.Hour, time.Hour)

	for range 2 {
		if _, err := cached.GetCandles(context.Background(), "AAPL", "D", time.Now(), time.Now()); err == nil {
			t.Fatal("expected error to propagate from inner service")
		}
	}

	if inner.calls != 2 {
		t.Errorf("inner service was called %d times, want 2 (failed fetches should not be cached)", inner.calls)
	}
}
