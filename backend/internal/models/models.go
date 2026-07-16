package models

import "time"

type Company struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Sector        string  `json:"sector"`
	Industry      string  `json:"industry"`
	MarketCap     float64 `json:"marketCap"`
	PERatio       float64 `json:"peRatio"`
	DividendYield float64 `json:"dividendYield"`
	WeekHigh52    float64 `json:"weekHigh52"`
	WeekLow52     float64 `json:"weekLow52"`
}

type NewsSummary struct {
	CompanySymbol   string    `json:"companySymbol"`
	Headline        string    `json:"headline"`
	Summary         string    `json:"summary"`
	Sentiment       string    `json:"sentiment"`
	PotentialImpact string    `json:"potentialImpact"`
	ArticleURL      string    `json:"articleUrl"`
	PublishedAt     time.Time `json:"publishedAt"`
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
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	Description          string         `json:"description"`
	Holdings             []HoldingDetail `json:"holdings"`
	TotalValue           float64        `json:"totalValue"`
	TotalCost            float64        `json:"totalCost"`
	TotalGainLoss        float64        `json:"totalGainLoss"`
	TotalGainLossPercent float64        `json:"totalGainLossPercent"`
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
