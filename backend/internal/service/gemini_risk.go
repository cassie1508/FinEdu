package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"finedu-backend/internal/models"
)

// GeminiRiskService calls the Gemini API to analyze portfolio risk — the
// same provider already used for news summaries, so no separate AI API key
// is needed beyond GEMINI_API_KEY.
type GeminiRiskService struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

type GeminiRiskConfig struct {
	APIKey  string
	BaseURL string // optional, defaults to Gemini's generateContent endpoint template
	Model   string // optional, defaults to "gemini-2.0-flash"
}

func NewGeminiRiskService(cfg GeminiRiskConfig) *GeminiRiskService {
	baseURL := geminiGenerateContentURL
	if cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}
	model := "gemini-2.0-flash"
	if cfg.Model != "" {
		model = cfg.Model
	}
	return &GeminiRiskService{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *GeminiRiskService) AnalyzeRisk(ctx context.Context, holdings []HoldingForAnalysis) (*models.PortfolioRisk, error) {
	holdingsJSON, err := json.Marshal(holdings)
	if err != nil {
		return nil, fmt.Errorf("marshaling holdings: %w", err)
	}

	prompt := fmt.Sprintf(
		"%s\n\nHoldings:\n%s",
		systemPrompt, string(holdingsJSON),
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
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	url := fmt.Sprintf(s.baseURL, s.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling Gemini: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
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
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	if len(completion.Candidates) == 0 || len(completion.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no candidates in Gemini response")
	}

	var output aiRiskOutput
	if err := json.Unmarshal([]byte(completion.Candidates[0].Content.Parts[0].Text), &output); err != nil {
		return nil, fmt.Errorf("parsing AI JSON output: %w", err)
	}

	healthScore := clamp(output.HealthScore, 0, 100)
	divScore := clamp(output.DiversificationScore, 0, 100)
	riskLevel := validateRiskLevel(output.RiskLevel)

	if output.SectorConcentration == nil {
		output.SectorConcentration = []models.SectorEntry{}
	}
	if output.Recommendations == nil {
		output.Recommendations = []string{}
	}

	return &models.PortfolioRisk{
		HealthScore:          &healthScore,
		RiskLevel:            &riskLevel,
		SectorConcentration:  output.SectorConcentration,
		DiversificationScore: &divScore,
		Recommendations:      output.Recommendations,
	}, nil
}
