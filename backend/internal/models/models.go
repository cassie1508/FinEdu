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

type PortfolioHolding struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Symbol      string    `json:"symbol"`
	Shares      float64   `json:"shares"`
	AverageCost float64   `json:"averageCost"`
	CreatedAt   time.Time `json:"createdAt"`
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
