package service

import (
	"context"

	"finedu-backend/internal/models"
)

// AIRiskAnalyzer generates AI-powered portfolio risk analysis.
type AIRiskAnalyzer interface {
	AnalyzeRisk(ctx context.Context, holdings []HoldingForAnalysis) (*models.PortfolioRisk, error)
}

// HoldingForAnalysis is the minimal info the AI service needs per holding.
type HoldingForAnalysis struct {
	Symbol            string
	Shares            float64
	AllocationPercent float64
	CurrentPrice      float64
}
