package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"finedu-backend/internal/models"
)

// NewsProvider fetches raw news articles for a ticker or for a general
// market category.
type NewsProvider interface {
	GetTickerNews(ctx context.Context, ticker string) ([]models.NewsArticle, error)
	GetGeneralNews(ctx context.Context, category string) ([]models.NewsArticle, error)
}

const finnhubBaseURL = "https://finnhub.io/api/v1"

// finnhubCategories maps the API's category query param to Finnhub's /news
// taxonomy (general, forex, crypto, merger). Finnhub has no dedicated
// "earnings" feed, so that one is handled separately: see
// FinnhubNewsProvider.GetGeneralNews. "market" isn't here either — Finnhub's
// "general" category is broad business/economy news (Fed policy, oil prices,
// etc.), not stock-specific, so it also needs client-side filtering.
var finnhubCategories = map[string]string{
	"mergers": "merger",
}

// earningsKeywords approximates an "earnings" news feed by filtering
// Finnhub's general category, since Finnhub's /news endpoint only supports
// general/forex/crypto/merger — there's no earnings-specific category the
// way Alpha Vantage's topics= offered.
var earningsKeywords = []string{
	"earnings", "eps", "quarterly results", "quarterly report",
	"guidance", "beats estimates", "misses estimates", "earnings call",
	"earnings report", "profit forecast", "revenue forecast",
}

// stockKeywords approximates a "stocks / how companies are doing" news feed
// by filtering Finnhub's general category, which otherwise mixes broad
// economy, commodities, and policy stories in alongside actual stock-market
// news. Same heuristic (and same limitations) as earningsKeywords above.
var stockKeywords = []string{
	"stock", "stocks", "shares", "share price", "equity", "equities",
	"nasdaq", "dow jones", "s&p 500", "wall street", "ticker",
	"market cap", "ipo", "buyback", "dividend", "analyst", "price target",
	"rally", "rallies", "sell-off", "selloff", "earnings", "revenue",
	"guidance", "upgrade", "downgrade", "outlook", "quarterly",
}

type FinnhubNewsProvider struct {
	apiKey         string
	httpClient     *http.Client
	baseURL        string
	tickerLookback time.Duration
}

func NewFinnhubNewsProvider(apiKey string) *FinnhubNewsProvider {
	return &FinnhubNewsProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: finnhubBaseURL,
		// Finnhub's company-news endpoint requires an explicit from/to
		// range rather than a limit param, so this bounds how far back one
		// upstream fetch looks — wide enough to cover any since/limit a
		// caller asks for, without pulling a company's entire news history.
		tickerLookback: 30 * 24 * time.Hour,
	}
}

// finnhubNewsItem is shared by both the company-news and market-news
// endpoints — same shape either way.
type finnhubNewsItem struct {
	ID       int64  `json:"id"`
	Headline string `json:"headline"`
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	URL      string `json:"url"`
	Datetime int64  `json:"datetime"` // unix seconds
	Image    string `json:"image"`
}

func (f *FinnhubNewsProvider) GetTickerNews(ctx context.Context, ticker string) ([]models.NewsArticle, error) {
	to := time.Now()
	from := to.Add(-f.tickerLookback)

	reqURL := fmt.Sprintf(
		"%s/company-news?symbol=%s&from=%s&to=%s&token=%s",
		f.baseURL, url.QueryEscape(strings.ToUpper(ticker)),
		from.Format("2006-01-02"), to.Format("2006-01-02"),
		url.QueryEscape(f.apiKey),
	)

	items, err := f.fetch(ctx, reqURL)
	if err != nil {
		return nil, err
	}
	return toNewsArticles(items), nil
}

func (f *FinnhubNewsProvider) GetGeneralNews(ctx context.Context, category string) ([]models.NewsArticle, error) {
	switch category {
	case "earnings":
		items, err := f.fetchCategory(ctx, "general")
		if err != nil {
			return nil, err
		}
		return toNewsArticles(filterEarningsRelated(items)), nil
	case "market":
		items, err := f.fetchCategory(ctx, "general")
		if err != nil {
			return nil, err
		}
		return toNewsArticles(filterStockRelated(items)), nil
	}

	finnhubCategory, ok := finnhubCategories[category]
	if !ok {
		return nil, fmt.Errorf("unsupported news category %q", category)
	}

	items, err := f.fetchCategory(ctx, finnhubCategory)
	if err != nil {
		return nil, err
	}
	return toNewsArticles(items), nil
}

func (f *FinnhubNewsProvider) fetchCategory(ctx context.Context, category string) ([]finnhubNewsItem, error) {
	reqURL := fmt.Sprintf("%s/news?category=%s&token=%s", f.baseURL, url.QueryEscape(category), url.QueryEscape(f.apiKey))
	return f.fetch(ctx, reqURL)
}

func (f *FinnhubNewsProvider) fetch(ctx context.Context, reqURL string) ([]finnhubNewsItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching news: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading news response: %w", err)
	}

	// Finnhub returns 401/403 for a bad/unauthorized key and 429 once the
	// per-minute call limit is exceeded — all "try again later", not a
	// caller/transport bug.
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("%w: finnhub API error (status %d): %s", ErrUpstreamUnavailable, resp.StatusCode, string(body))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("finnhub API error (status %d): %s", resp.StatusCode, string(body))
	}

	var items []finnhubNewsItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decoding news response: %w", err)
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Datetime > items[j].Datetime })

	return items, nil
}

func toNewsArticles(items []finnhubNewsItem) []models.NewsArticle {
	articles := make([]models.NewsArticle, 0, len(items))
	for _, item := range items {
		if item.Headline == "" || item.URL == "" {
			continue
		}
		articles = append(articles, models.NewsArticle{
			ID:          strconv.FormatInt(item.ID, 10),
			Headline:    item.Headline,
			Summary:     item.Summary,
			Source:      item.Source,
			URL:         item.URL,
			PublishedAt: time.Unix(item.Datetime, 0).UTC(),
			ImageURL:    item.Image,
		})
	}
	return articles
}

func filterEarningsRelated(items []finnhubNewsItem) []finnhubNewsItem {
	filtered := make([]finnhubNewsItem, 0, len(items))
	for _, item := range items {
		text := strings.ToLower(item.Headline + " " + item.Summary)
		for _, kw := range earningsKeywords {
			if strings.Contains(text, kw) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}

func filterStockRelated(items []finnhubNewsItem) []finnhubNewsItem {
	filtered := make([]finnhubNewsItem, 0, len(items))
	for _, item := range items {
		text := strings.ToLower(item.Headline + " " + item.Summary)
		for _, kw := range stockKeywords {
			if strings.Contains(text, kw) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}

// FilterArticles returns articles published at or after `since` (a zero
// value means no lower bound), capped at `limit` (<=0 means no cap).
// Assumes articles is already sorted newest-first.
func FilterArticles(articles []models.NewsArticle, since time.Time, limit int) []models.NewsArticle {
	filtered := make([]models.NewsArticle, 0, len(articles))
	for _, a := range articles {
		if !since.IsZero() && a.PublishedAt.Before(since) {
			continue
		}
		filtered = append(filtered, a)
		if limit > 0 && len(filtered) >= limit {
			break
		}
	}
	return filtered
}

// --- Cached wrapper ---

type newsCacheEntry struct {
	articles  []models.NewsArticle
	fetchedAt time.Time
}

// CachedNewsProvider caches the full article set per ticker/category and
// leaves since/limit filtering to the caller on every read — same rationale
// as CachedMarketData's window filtering: caching a pre-filtered result
// would let whichever request hits the API first decide the window for
// every other request sharing that key.
type CachedNewsProvider struct {
	inner NewsProvider
	ttl   time.Duration
	mu    sync.RWMutex
	cache map[string]newsCacheEntry
}

func NewCachedNewsProvider(inner NewsProvider, ttl time.Duration) *CachedNewsProvider {
	return &CachedNewsProvider{
		inner: inner,
		ttl:   ttl,
		cache: make(map[string]newsCacheEntry),
	}
}

func (c *CachedNewsProvider) GetTickerNews(ctx context.Context, ticker string) ([]models.NewsArticle, error) {
	return c.get("ticker|"+strings.ToUpper(ticker), func() ([]models.NewsArticle, error) {
		return c.inner.GetTickerNews(ctx, ticker)
	})
}

func (c *CachedNewsProvider) GetGeneralNews(ctx context.Context, category string) ([]models.NewsArticle, error) {
	return c.get("category|"+category, func() ([]models.NewsArticle, error) {
		return c.inner.GetGeneralNews(ctx, category)
	})
}

func (c *CachedNewsProvider) get(key string, fetch func() ([]models.NewsArticle, error)) ([]models.NewsArticle, error) {
	c.mu.RLock()
	entry, ok := c.cache[key]
	c.mu.RUnlock()

	if ok && time.Since(entry.fetchedAt) < c.ttl {
		return entry.articles, nil
	}

	articles, err := fetch()
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[key] = newsCacheEntry{articles: articles, fetchedAt: time.Now()}
	c.mu.Unlock()

	return articles, nil
}

// --- AI daily summary ---

var ErrNoArticles = errors.New("no articles available to summarize")

// Summarizer turns a ticker's raw articles into an AI-generated daily
// digest.
type Summarizer interface {
	Summarize(ctx context.Context, ticker string, articles []models.NewsArticle) (models.NewsDailySummary, error)
}

const geminiGenerateContentURL = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent"

type GeminiSummarizer struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewGeminiSummarizer(apiKey string) *GeminiSummarizer {
	return &GeminiSummarizer{
		apiKey: apiKey,
		model:  "gemini-2.0-flash",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: geminiGenerateContentURL,
	}
}

// geminiSummaryJSON is the shape asked of the model — sourceArticleIds is
// deliberately not part of it; that's derived by us from the input articles
// rather than left for the model to transcribe (and possibly hallucinate).
type geminiSummaryJSON struct {
	DailySummary    string  `json:"dailySummary"`
	Sentiment       string  `json:"sentiment"`
	SentimentScore  float64 `json:"sentimentScore"`
	PotentialImpact string  `json:"potentialImpact"`
}

func (g *GeminiSummarizer) Summarize(ctx context.Context, ticker string, articles []models.NewsArticle) (models.NewsDailySummary, error) {
	if len(articles) == 0 {
		return models.NewsDailySummary{}, ErrNoArticles
	}

	var headlines strings.Builder
	ids := make([]string, 0, len(articles))
	for _, a := range articles {
		fmt.Fprintf(&headlines, "- [%s] (%s) %s\n", a.ID, a.Source, a.Headline)
		ids = append(ids, a.ID)
	}

	prompt := fmt.Sprintf(
		"You are a financial news analyst. Based only on the following recent headlines for %s, "+
			"write a concise 2-3 sentence daily summary, an overall sentiment "+
			"(one of \"bullish\", \"neutral\", \"bearish\"), a sentiment score from -1.0 (very bearish) "+
			"to 1.0 (very bullish), and a short note on potential price impact. "+
			"Respond with ONLY a JSON object with keys: dailySummary, sentiment, sentimentScore, potentialImpact.\n\nHeadlines:\n%s",
		ticker, headlines.String(),
	)

	reqBody, err := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]any{
			"responseMimeType": "application/json",
			"temperature":      0.3,
		},
	})
	if err != nil {
		return models.NewsDailySummary{}, fmt.Errorf("encoding Gemini request: %w", err)
	}

	url := fmt.Sprintf(g.baseURL, g.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return models.NewsDailySummary{}, fmt.Errorf("creating Gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return models.NewsDailySummary{}, fmt.Errorf("calling Gemini: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.NewsDailySummary{}, fmt.Errorf("reading Gemini response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.NewsDailySummary{}, fmt.Errorf("%w: Gemini API error (status %d): %s", ErrUpstreamUnavailable, resp.StatusCode, string(body))
	}

	var completion struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &completion); err != nil {
		return models.NewsDailySummary{}, fmt.Errorf("decoding Gemini response: %w", err)
	}
	if len(completion.Candidates) == 0 || len(completion.Candidates[0].Content.Parts) == 0 {
		return models.NewsDailySummary{}, fmt.Errorf("Gemini response had no candidates")
	}

	var parsed geminiSummaryJSON
	if err := json.Unmarshal([]byte(completion.Candidates[0].Content.Parts[0].Text), &parsed); err != nil {
		return models.NewsDailySummary{}, fmt.Errorf("decoding Gemini summary JSON: %w", err)
	}

	return models.NewsDailySummary{
		Ticker:           strings.ToUpper(ticker),
		Date:             time.Now().UTC().Format("2006-01-02"),
		DailySummary:     parsed.DailySummary,
		Sentiment:        parsed.Sentiment,
		SentimentScore:   parsed.SentimentScore,
		PotentialImpact:  parsed.PotentialImpact,
		SourceArticleIDs: ids,
	}, nil
}

// --- Cached summary service (regenerated on a schedule, not per request) ---

// SummaryService generates the AI daily digest for a ticker. Implemented by
// NewsSummaryService; an interface here lets handlers be tested against a
// stub instead of a real Gemini-backed instance.
type SummaryService interface {
	GetSummary(ctx context.Context, ticker string) (models.NewsDailySummary, error)
}

type summaryCacheEntry struct {
	summary     models.NewsDailySummary
	generatedAt time.Time
}

// NewsSummaryService builds the AI daily summary from cached raw articles
// and caches the result itself, regenerating at most every
// `regenerateEvery` rather than per request. `articleWindow` bounds how far
// back cached articles are considered "today's news" when building a
// summary; `articleLimit` caps how many are sent to Gemini, bounding prompt
// size and cost.
type NewsSummaryService struct {
	articles        NewsProvider
	summarizer      Summarizer
	regenerateEvery time.Duration
	articleWindow   time.Duration
	articleLimit    int
	mu              sync.RWMutex
	cache           map[string]summaryCacheEntry
}

func NewNewsSummaryService(articles NewsProvider, summarizer Summarizer, regenerateEvery, articleWindow time.Duration, articleLimit int) *NewsSummaryService {
	return &NewsSummaryService{
		articles:        articles,
		summarizer:      summarizer,
		regenerateEvery: regenerateEvery,
		articleWindow:   articleWindow,
		articleLimit:    articleLimit,
		cache:           make(map[string]summaryCacheEntry),
	}
}

func (s *NewsSummaryService) GetSummary(ctx context.Context, ticker string) (models.NewsDailySummary, error) {
	key := strings.ToUpper(ticker)

	s.mu.RLock()
	entry, ok := s.cache[key]
	s.mu.RUnlock()

	if ok && time.Since(entry.generatedAt) < s.regenerateEvery {
		return entry.summary, nil
	}

	all, err := s.articles.GetTickerNews(ctx, ticker)
	if err != nil {
		return models.NewsDailySummary{}, err
	}

	recent := FilterArticles(all, time.Now().Add(-s.articleWindow), s.articleLimit)

	summary, err := s.summarizer.Summarize(ctx, ticker, recent)
	if err != nil {
		return models.NewsDailySummary{}, err
	}

	s.mu.Lock()
	s.cache[key] = summaryCacheEntry{summary: summary, generatedAt: time.Now()}
	s.mu.Unlock()

	return summary, nil
}
