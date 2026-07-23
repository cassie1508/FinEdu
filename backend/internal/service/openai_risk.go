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

// OpenAIRiskService calls the OpenAI Chat Completions API to analyze portfolio risk.
type OpenAIRiskService struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

type OpenAIRiskConfig struct {
	APIKey  string
	BaseURL string // optional, defaults to "https://api.openai.com"
	Model   string // optional, defaults to "gpt-4o-mini"
}

func NewOpenAIRiskService(cfg OpenAIRiskConfig) *OpenAIRiskService {
	baseURL := "https://api.openai.com"
	if cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}
	model := "gpt-4o-mini"
	if cfg.Model != "" {
		model = cfg.Model
	}
	return &OpenAIRiskService{
		apiKey:  cfg.APIKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// --- OpenAI request/response types ---

type openAIChatRequest struct {
	Model          string              `json:"model"`
	Messages       []openAIChatMessage `json:"messages"`
	ResponseFormat openAIResponseFormat `json:"response_format"`
	Temperature    float64             `json:"temperature"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// aiRiskOutput is the JSON schema we instruct the AI to return.
type aiRiskOutput struct {
	HealthScore          float64              `json:"healthScore"`
	RiskLevel            string               `json:"riskLevel"`
	SectorConcentration  []models.SectorEntry `json:"sectorConcentration"`
	DiversificationScore float64              `json:"diversificationScore"`
	Recommendations      []string             `json:"recommendations"`
}

const systemPrompt = `You are a professional portfolio risk analyst AI. Given a list of stock holdings, analyze the portfolio and return ONLY a valid JSON object with exactly these fields:
- healthScore: integer 0-100 (100 = excellent health)
- riskLevel: exactly one of "Low", "Moderate", "High", "Very High"
- sectorConcentration: array of {"sector": string, "percent": number} inferred from the stock symbols
- diversificationScore: integer 0-100 (100 = perfectly diversified across sectors)
- recommendations: array of exactly 3-5 actionable strings to improve the portfolio

Return only the JSON object, no markdown, no explanation.`

func (s *OpenAIRiskService) AnalyzeRisk(ctx context.Context, holdings []HoldingForAnalysis) (*models.PortfolioRisk, error) {
	holdingsJSON, err := json.Marshal(holdings)
	if err != nil {
		return nil, fmt.Errorf("marshaling holdings: %w", err)
	}

	userMessage := fmt.Sprintf("Analyze this portfolio and return risk analysis JSON:\n%s", string(holdingsJSON))

	reqBody := openAIChatRequest{
		Model: s.model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		ResponseFormat: openAIResponseFormat{Type: "json_object"},
		Temperature:    0.3,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		s.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in OpenAI response")
	}

	var output aiRiskOutput
	if err := json.Unmarshal([]byte(chatResp.Choices[0].Message.Content), &output); err != nil {
		return nil, fmt.Errorf("parsing AI JSON output: %w", err)
	}

	// Validate and clamp score ranges
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

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func validateRiskLevel(level string) string {
	switch level {
	case "Low", "Moderate", "High", "Very High":
		return level
	default:
		return "Moderate"
	}
}
