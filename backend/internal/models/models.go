package models

import "time"

type Company struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Sector        string  `json:"sector"`
	Industry      string  `json:"industry"`
	MarketCap     float64 `json:"marketCap"`
	Revenue       float64 `json:"revenue"`
	EPS           float64 `json:"eps"`
	PERatio       float64 `json:"peRatio"`
	DividendYield float64 `json:"dividendYield"`
	WeekHigh52    float64 `json:"weekHigh52"`
	WeekLow52     float64 `json:"weekLow52"`
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

// --- Portfolio Management Models ---

type Portfolio struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PortfolioHolding struct {
	ID          string    `json:"id"`
	PortfolioID string    `json:"portfolioId"`
	Symbol      string    `json:"symbol"`
	Shares      float64   `json:"shares"`
	AverageCost float64   `json:"averageCost"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// HoldingDetail is a holding enriched with live market data (not stored in DB).
type HoldingDetail struct {
	ID                        string  `json:"id"`
	PortfolioID               string  `json:"portfolioId"`
	Symbol                    string  `json:"symbol"`
	Shares                    float64 `json:"shares"`
	AverageCost               float64 `json:"averageCost"`
	CurrentPrice              float64 `json:"currentPrice"`
	CurrentValue              float64 `json:"currentValue"`
	UnrealizedGainLoss        float64 `json:"unrealizedGainLoss"`
	UnrealizedGainLossPercent float64 `json:"unrealizedGainLossPercent"`
	AllocationPercent         float64 `json:"allocationPercent"`
}

type PortfolioDetail struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	Holdings             []HoldingDetail `json:"holdings"`
	TotalValue           float64         `json:"totalValue"`
	TotalCost            float64         `json:"totalCost"`
	TotalGainLoss        float64         `json:"totalGainLoss"`
	TotalGainLossPercent float64         `json:"totalGainLossPercent"`
}

type PortfolioListItem struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	HoldingsCount int       `json:"holdingsCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PortfolioSummary struct {
	TotalValue           float64           `json:"totalValue"`
	TotalCost            float64           `json:"totalCost"`
	TotalGainLoss        float64           `json:"totalGainLoss"`
	TotalGainLossPercent float64           `json:"totalGainLossPercent"`
	Allocations          []AllocationEntry `json:"allocations"`
}

type AllocationEntry struct {
	Symbol  string  `json:"symbol"`
	Percent float64 `json:"percent"`
}

// --- Market Data Models ---

type StockQuote struct {
	Symbol        string  `json:"symbol"`
	CurrentPrice  float64 `json:"currentPrice"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

// --- Portfolio Risk Models (Week 3 placeholder) ---

type PortfolioRisk struct {
	HealthScore          *float64      `json:"healthScore"`
	RiskLevel            *string       `json:"riskLevel"`
	SectorConcentration  []SectorEntry `json:"sectorConcentration"`
	DiversificationScore *float64      `json:"diversificationScore"`
	Recommendations      []string      `json:"recommendations"`
	Message              string        `json:"message,omitempty"`
}

type SectorEntry struct {
	Sector  string  `json:"sector"`
	Percent float64 `json:"percent"`
}

// --- Request DTOs ---

type CreatePortfolioRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description"`
}

type AddHoldingRequest struct {
	Symbol      string  `json:"symbol" binding:"required,max=10"`
	Shares      float64 `json:"shares" binding:"required,gt=0"`
	AverageCost float64 `json:"averageCost" binding:"required,gte=0"`
}

type UpdateHoldingRequest struct {
	Shares      float64 `json:"shares" binding:"required,gt=0"`
	AverageCost float64 `json:"averageCost" binding:"required,gte=0"`
}
