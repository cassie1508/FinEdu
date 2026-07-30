package models

import "time"

type Company struct {
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	Sector         string  `json:"sector"`
	Industry       string  `json:"industry"`
	MarketCap      float64 `json:"marketCap"`
	PERatio        float64 `json:"peRatio"`
	DividendYield  float64 `json:"dividendYield"`
	WeekHigh52     float64 `json:"weekHigh52"`
	WeekLow52      float64 `json:"weekLow52"`
}

type PortfolioHolding struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"userId"`
	Symbol             string    `json:"symbol"`
	Shares             float64   `json:"shares"`
	AverageCost        float64   `json:"averageCost"`
	CreatedAt          time.Time `json:"createdAt"`
}

// CandlePoint is a single OHLCV bar for a given period.
type CandlePoint struct {
	Time   int64   `json:"t"` // unix seconds
	Open   float64 `json:"o"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Close  float64 `json:"c"`
	Volume float64 `json:"v"`
}

// PriceHistory is the chart data returned for a ticker + time range.
type PriceHistory struct {
	Symbol         string        `json:"symbol"`
	Range          string        `json:"range"`
	Resolution     string        `json:"resolution"`
	Candles        []CandlePoint `json:"candles"`
	DelayedMinutes int           `json:"delayedMinutes"`
}

// NewsArticle is a single raw news item, either for a ticker or from the
// general market feed.
type NewsArticle struct {
	ID          string    `json:"id"`
	Headline    string    `json:"headline"`
	Summary     string    `json:"summary"`
	Source      string    `json:"source"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
	ImageURL    string    `json:"imageUrl"`
}

// NewsDailySummary is the AI-generated daily digest for a ticker, built from
// its cached raw articles and regenerated on a schedule rather than per
// request.
type NewsDailySummary struct {
	Ticker           string   `json:"ticker"`
	Date             string   `json:"date"` // YYYY-MM-DD
	DailySummary     string   `json:"dailySummary"`
	Sentiment        string   `json:"sentiment"` // bullish|neutral|bearish
	SentimentScore   float64  `json:"sentimentScore"`
	PotentialImpact  string   `json:"potentialImpact"`
	SourceArticleIDs []string `json:"sourceArticleIds"`
}
