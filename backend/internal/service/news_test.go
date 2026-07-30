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

func TestFinnhubNewsProvider_GetTickerNews(t *testing.T) {
	body := `[
		{"id": 1, "headline": "AAPL rallies", "summary": "shares up", "source": "Reuters", "url": "https://example.com/a", "datetime": 1784800000, "image": "https://example.com/a.png"},
		{"id": 2, "headline": "AAPL dips", "summary": "shares down", "source": "Bloomberg", "url": "https://example.com/b", "datetime": 1784900000}
	]`
	srv := newTestServer(t, body, http.StatusOK)
	defer srv.Close()

	fh := NewFinnhubNewsProvider("test-key")
	fh.baseURL = srv.URL

	articles, err := fh.GetTickerNews(context.Background(), "aapl")
	if err != nil {
		t.Fatalf("GetTickerNews returned error: %v", err)
	}
	if len(articles) != 2 {
		t.Fatalf("got %d articles, want 2", len(articles))
	}

	// Newest first.
	if articles[0].Headline != "AAPL dips" {
		t.Errorf("articles[0].Headline = %q, want %q (newest first)", articles[0].Headline, "AAPL dips")
	}
	if articles[0].ID != "2" {
		t.Errorf("articles[0].ID = %q, want %q (Finnhub's numeric id)", articles[0].ID, "2")
	}
	if articles[1].ImageURL != "https://example.com/a.png" {
		t.Errorf("articles[1].ImageURL = %q, want %q", articles[1].ImageURL, "https://example.com/a.png")
	}
}

func TestFinnhubNewsProvider_GetGeneralNews_MarketAndMergersFetchExpectedFinnhubCategory(t *testing.T) {
	var gotCategory string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCategory = r.URL.Query().Get("category")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	fh := NewFinnhubNewsProvider("test-key")
	fh.baseURL = srv.URL

	// "market" is client-side filtered (see the test below) but still reads
	// from Finnhub's "general" category — there's no dedicated stocks feed.
	if _, err := fh.GetGeneralNews(context.Background(), "market"); err != nil {
		t.Fatalf("GetGeneralNews(market) returned error: %v", err)
	}
	if gotCategory != "general" {
		t.Errorf("category=market sent Finnhub category %q, want %q", gotCategory, "general")
	}

	if _, err := fh.GetGeneralNews(context.Background(), "mergers"); err != nil {
		t.Fatalf("GetGeneralNews(mergers) returned error: %v", err)
	}
	if gotCategory != "merger" {
		t.Errorf("category=mergers sent Finnhub category %q, want %q", gotCategory, "merger")
	}
}

// Finnhub's "general" category mixes broad economy/commodities/policy
// stories in with actual stock-market news, so GetGeneralNews("market")
// keyword-filters it client-side — same approach as the "earnings" filter.
func TestFinnhubNewsProvider_GetGeneralNews_MarketFiltersToStockRelated(t *testing.T) {
	body := `[
		{"id": 1, "headline": "Fed signals potential rate cut later this year", "summary": "the Federal Reserve hinted at possible rate cuts", "source": "CNBC", "url": "https://example.com/a", "datetime": 1784800000},
		{"id": 2, "headline": "Oil prices dip on demand concerns", "summary": "crude oil prices fell", "source": "MarketWatch", "url": "https://example.com/b", "datetime": 1784850000},
		{"id": 3, "headline": "Tech stocks rally as inflation cools", "summary": "shares of major tech companies gained", "source": "Reuters", "url": "https://example.com/c", "datetime": 1784900000}
	]`
	srv := newTestServer(t, body, http.StatusOK)
	defer srv.Close()

	fh := NewFinnhubNewsProvider("test-key")
	fh.baseURL = srv.URL

	articles, err := fh.GetGeneralNews(context.Background(), "market")
	if err != nil {
		t.Fatalf("GetGeneralNews(market) returned error: %v", err)
	}
	if len(articles) != 1 || articles[0].ID != "3" {
		t.Errorf("got %+v, want only the stock-related article", articles)
	}
}

// Finnhub's /news endpoint has no "earnings" category (only
// general/forex/crypto/merger), so GetGeneralNews("earnings") fetches
// "general" and keyword-filters it client-side.
func TestFinnhubNewsProvider_GetGeneralNews_EarningsFiltersGeneralFeed(t *testing.T) {
	body := `[
		{"id": 1, "headline": "Tech stocks rally on Fed comments", "summary": "broad market move", "source": "Reuters", "url": "https://example.com/a", "datetime": 1784800000},
		{"id": 2, "headline": "Acme Corp beats earnings estimates", "summary": "EPS of $1.20 vs $1.05 expected", "source": "Bloomberg", "url": "https://example.com/b", "datetime": 1784900000}
	]`
	srv := newTestServer(t, body, http.StatusOK)
	defer srv.Close()

	fh := NewFinnhubNewsProvider("test-key")
	fh.baseURL = srv.URL

	articles, err := fh.GetGeneralNews(context.Background(), "earnings")
	if err != nil {
		t.Fatalf("GetGeneralNews(earnings) returned error: %v", err)
	}
	if len(articles) != 1 || articles[0].ID != "2" {
		t.Errorf("got %+v, want only the earnings-related article", articles)
	}
}

func TestFinnhubNewsProvider_GetGeneralNews_InvalidCategory(t *testing.T) {
	fh := NewFinnhubNewsProvider("test-key")
	if _, err := fh.GetGeneralNews(context.Background(), "bogus"); err == nil {
		t.Error("expected an error for an unsupported category, got nil")
	}
}

func TestFinnhubNewsProvider_UpstreamDeclined(t *testing.T) {
	for _, status := range []int{http.StatusTooManyRequests, http.StatusUnauthorized, http.StatusForbidden} {
		srv := newTestServer(t, `{"error":"rate limited"}`, status)

		fh := NewFinnhubNewsProvider("test-key")
		fh.baseURL = srv.URL

		_, err := fh.GetTickerNews(context.Background(), "AAPL")
		srv.Close()

		if !errors.Is(err, ErrUpstreamUnavailable) {
			t.Errorf("status %d: expected ErrUpstreamUnavailable, got: %v", status, err)
		}
	}
}

func TestFilterArticles(t *testing.T) {
	base := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	articles := []models.NewsArticle{
		{ID: "1", PublishedAt: base},
		{ID: "2", PublishedAt: base.AddDate(0, 0, -1)},
		{ID: "3", PublishedAt: base.AddDate(0, 0, -5)},
	}

	got := FilterArticles(articles, base.AddDate(0, 0, -2), 0)
	if len(got) != 2 {
		t.Fatalf("got %d articles within the window, want 2", len(got))
	}

	limited := FilterArticles(articles, time.Time{}, 1)
	if len(limited) != 1 || limited[0].ID != "1" {
		t.Errorf("limit=1 got %+v, want just the first (newest) article", limited)
	}
}

// fakeNewsProvider counts calls per method, so CachedNewsProvider's caching
// behavior can be verified without a network call.
type fakeNewsProvider struct {
	tickerCalls  int
	generalCalls int
	articles     []models.NewsArticle
	err          error
}

func (f *fakeNewsProvider) GetTickerNews(ctx context.Context, ticker string) ([]models.NewsArticle, error) {
	f.tickerCalls++
	return f.articles, f.err
}

func (f *fakeNewsProvider) GetGeneralNews(ctx context.Context, category string) ([]models.NewsArticle, error) {
	f.generalCalls++
	return f.articles, f.err
}

func TestCachedNewsProvider_ServesFromCacheWithinTTL(t *testing.T) {
	inner := &fakeNewsProvider{articles: []models.NewsArticle{{ID: "1"}}}
	cached := NewCachedNewsProvider(inner, time.Hour)

	for range 3 {
		if _, err := cached.GetTickerNews(context.Background(), "AAPL"); err != nil {
			t.Fatalf("GetTickerNews returned error: %v", err)
		}
	}

	if inner.tickerCalls != 1 {
		t.Errorf("inner provider was called %d times, want 1 (repeat calls within TTL should be cached)", inner.tickerCalls)
	}
}

func TestCachedNewsProvider_TickerAndCategoryCachedIndependently(t *testing.T) {
	inner := &fakeNewsProvider{articles: []models.NewsArticle{{ID: "1"}}}
	cached := NewCachedNewsProvider(inner, time.Hour)

	if _, err := cached.GetTickerNews(context.Background(), "AAPL"); err != nil {
		t.Fatalf("GetTickerNews returned error: %v", err)
	}
	if _, err := cached.GetGeneralNews(context.Background(), "market"); err != nil {
		t.Fatalf("GetGeneralNews returned error: %v", err)
	}

	if inner.tickerCalls != 1 || inner.generalCalls != 1 {
		t.Errorf("tickerCalls=%d generalCalls=%d, want 1 and 1 (distinct cache keys)", inner.tickerCalls, inner.generalCalls)
	}
}

func TestCachedNewsProvider_RefetchesAfterTTLExpires(t *testing.T) {
	inner := &fakeNewsProvider{articles: []models.NewsArticle{{ID: "1"}}}
	cached := NewCachedNewsProvider(inner, time.Hour)

	if _, err := cached.GetTickerNews(context.Background(), "AAPL"); err != nil {
		t.Fatalf("GetTickerNews returned error: %v", err)
	}

	cached.mu.Lock()
	entry := cached.cache["ticker|AAPL"]
	entry.fetchedAt = time.Now().Add(-2 * time.Hour)
	cached.cache["ticker|AAPL"] = entry
	cached.mu.Unlock()

	if _, err := cached.GetTickerNews(context.Background(), "AAPL"); err != nil {
		t.Fatalf("GetTickerNews returned error: %v", err)
	}

	if inner.tickerCalls != 2 {
		t.Errorf("inner provider was called %d times, want 2 (expired entry should trigger a refetch)", inner.tickerCalls)
	}
}

func TestCachedNewsProvider_DoesNotCacheErrors(t *testing.T) {
	inner := &fakeNewsProvider{err: errors.New("boom")}
	cached := NewCachedNewsProvider(inner, time.Hour)

	for range 2 {
		if _, err := cached.GetTickerNews(context.Background(), "AAPL"); err == nil {
			t.Fatal("expected error to propagate from inner provider")
		}
	}

	if inner.tickerCalls != 2 {
		t.Errorf("inner provider was called %d times, want 2 (failed fetches should not be cached)", inner.tickerCalls)
	}
}

// fakeSummarizer records the articles it was called with and returns a
// canned summary/error.
type fakeSummarizer struct {
	calls        int
	lastArticles []models.NewsArticle
	summary      models.NewsDailySummary
	err          error
}

func (f *fakeSummarizer) Summarize(ctx context.Context, ticker string, articles []models.NewsArticle) (models.NewsDailySummary, error) {
	f.calls++
	f.lastArticles = articles
	// Mirrors OpenAISummarizer's real contract: nothing to summarize means
	// ErrNoArticles, not a canned success.
	if len(articles) == 0 {
		return models.NewsDailySummary{}, ErrNoArticles
	}
	return f.summary, f.err
}

func TestNewsSummaryService_GetSummary_CachesAcrossRegenerateWindow(t *testing.T) {
	articles := &fakeNewsProvider{articles: []models.NewsArticle{{ID: "1", PublishedAt: time.Now()}}}
	summarizer := &fakeSummarizer{summary: models.NewsDailySummary{Ticker: "AAPL", Sentiment: "bullish"}}
	svc := NewNewsSummaryService(articles, summarizer, time.Hour, 48*time.Hour, 30)

	for range 3 {
		summary, err := svc.GetSummary(context.Background(), "aapl")
		if err != nil {
			t.Fatalf("GetSummary returned error: %v", err)
		}
		if summary.Sentiment != "bullish" {
			t.Errorf("summary.Sentiment = %q, want %q", summary.Sentiment, "bullish")
		}
	}

	if summarizer.calls != 1 {
		t.Errorf("summarizer was called %d times, want 1 (repeat calls within the regenerate window should be cached)", summarizer.calls)
	}
}

func TestNewsSummaryService_GetSummary_RegeneratesAfterWindowExpires(t *testing.T) {
	articles := &fakeNewsProvider{articles: []models.NewsArticle{{ID: "1", PublishedAt: time.Now()}}}
	summarizer := &fakeSummarizer{summary: models.NewsDailySummary{Ticker: "AAPL"}}
	svc := NewNewsSummaryService(articles, summarizer, time.Hour, 48*time.Hour, 30)

	if _, err := svc.GetSummary(context.Background(), "AAPL"); err != nil {
		t.Fatalf("GetSummary returned error: %v", err)
	}

	svc.mu.Lock()
	entry := svc.cache["AAPL"]
	entry.generatedAt = time.Now().Add(-2 * time.Hour)
	svc.cache["AAPL"] = entry
	svc.mu.Unlock()

	if _, err := svc.GetSummary(context.Background(), "AAPL"); err != nil {
		t.Fatalf("GetSummary returned error: %v", err)
	}

	if summarizer.calls != 2 {
		t.Errorf("summarizer was called %d times, want 2 (expired entry should trigger regeneration)", summarizer.calls)
	}
}

func TestNewsSummaryService_GetSummary_FiltersToArticleWindowAndLimit(t *testing.T) {
	now := time.Now()
	articles := &fakeNewsProvider{articles: []models.NewsArticle{
		{ID: "recent", PublishedAt: now},
		{ID: "stale", PublishedAt: now.AddDate(0, 0, -10)},
	}}
	summarizer := &fakeSummarizer{}
	svc := NewNewsSummaryService(articles, summarizer, time.Hour, 48*time.Hour, 30)

	if _, err := svc.GetSummary(context.Background(), "AAPL"); err != nil {
		t.Fatalf("GetSummary returned error: %v", err)
	}

	if len(summarizer.lastArticles) != 1 || summarizer.lastArticles[0].ID != "recent" {
		t.Errorf("summarizer got %+v, want only the article inside the 48h window", summarizer.lastArticles)
	}
}

func TestNewsSummaryService_GetSummary_PropagatesProviderError(t *testing.T) {
	articles := &fakeNewsProvider{err: ErrUpstreamUnavailable}
	svc := NewNewsSummaryService(articles, &fakeSummarizer{}, time.Hour, 48*time.Hour, 30)

	_, err := svc.GetSummary(context.Background(), "AAPL")
	if !errors.Is(err, ErrUpstreamUnavailable) {
		t.Errorf("expected ErrUpstreamUnavailable, got: %v", err)
	}
}

func TestNewsSummaryService_GetSummary_NoRecentArticlesPropagatesErrNoArticles(t *testing.T) {
	articles := &fakeNewsProvider{articles: []models.NewsArticle{
		{ID: "stale", PublishedAt: time.Now().AddDate(0, 0, -10)},
	}}
	svc := NewNewsSummaryService(articles, &fakeSummarizer{}, time.Hour, 48*time.Hour, 30)

	_, err := svc.GetSummary(context.Background(), "AAPL")
	if !errors.Is(err, ErrNoArticles) {
		t.Errorf("expected ErrNoArticles, got: %v", err)
	}
}
