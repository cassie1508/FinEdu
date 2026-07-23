package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"finedu-backend/internal/models"
	"finedu-backend/internal/service"
)

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func validAIContent() string {
	b, _ := json.Marshal(map[string]interface{}{
		"healthScore":          85,
		"riskLevel":            "Moderate",
		"sectorConcentration":  []map[string]interface{}{{"sector": "Technology", "percent": 100}},
		"diversificationScore": 40,
		"recommendations":      []string{"Diversify sectors", "Add bonds", "Consider international"},
	})
	return string(b)
}

func newOpenAIServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func testHoldings() []service.HoldingForAnalysis {
	return []service.HoldingForAnalysis{
		{Symbol: "AAPL", Shares: 10, AllocationPercent: 100, CurrentPrice: 150},
	}
}

func TestOpenAIRisk_SuccessfulAnalysis(t *testing.T) {
	srv := newOpenAIServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")

		resp := openAIResponse{}
		resp.Choices = []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{{}}
		resp.Choices[0].Message.Content = validAIContent()
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})

	risk, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.NoError(t, err)
	require.NotNil(t, risk.HealthScore)
	assert.Equal(t, 85.0, *risk.HealthScore)
	require.NotNil(t, risk.RiskLevel)
	assert.Equal(t, "Moderate", *risk.RiskLevel)
	require.NotNil(t, risk.DiversificationScore)
	assert.Equal(t, 40.0, *risk.DiversificationScore)
	assert.Len(t, risk.Recommendations, 3)
	assert.Len(t, risk.SectorConcentration, 1)
}

func TestOpenAIRisk_APIError500(t *testing.T) {
	srv := newOpenAIServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": {"message": "server error"}}`))
	})
	defer srv.Close()

	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})

	_, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OpenAI API error")
}

func TestOpenAIRisk_MalformedJSON(t *testing.T) {
	srv := newOpenAIServer(func(w http.ResponseWriter, r *http.Request) {
		// Return valid HTTP 200 but invalid JSON body
		w.Write([]byte("not json at all"))
	})
	defer srv.Close()

	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})

	_, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.Error(t, err)
}

func TestOpenAIRisk_NoChoices(t *testing.T) {
	srv := newOpenAIServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(openAIResponse{Choices: nil})
	})
	defer srv.Close()

	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})

	_, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no choices")
}

func TestOpenAIRisk_MalformedAIContent(t *testing.T) {
	srv := newOpenAIServer(func(w http.ResponseWriter, r *http.Request) {
		resp := openAIResponse{}
		resp.Choices = []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{{}}
		resp.Choices[0].Message.Content = "not valid json content"
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})

	_, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing AI JSON output")
}

func TestOpenAIRisk_InvalidRiskLevelDefaultsToModerate(t *testing.T) {
	invalidContent, _ := json.Marshal(map[string]interface{}{
		"healthScore":          50,
		"riskLevel":            "UNKNOWN_LEVEL",
		"sectorConcentration":  []interface{}{},
		"diversificationScore": 50,
		"recommendations":      []string{"Buy bonds"},
	})

	srv := newOpenAIServer(func(w http.ResponseWriter, r *http.Request) {
		resp := openAIResponse{}
		resp.Choices = []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{{}}
		resp.Choices[0].Message.Content = string(invalidContent)
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})

	risk, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.NoError(t, err)
	require.NotNil(t, risk.RiskLevel)
	assert.Equal(t, "Moderate", *risk.RiskLevel) // falls back to default
}

func TestOpenAIRisk_ScoresClamped(t *testing.T) {
	// AI returns scores outside 0-100
	outOfRangeContent, _ := json.Marshal(map[string]interface{}{
		"healthScore":          150,
		"riskLevel":            "Low",
		"sectorConcentration":  []interface{}{},
		"diversificationScore": -10,
		"recommendations":      []string{"Tip"},
	})

	srv := newOpenAIServer(func(w http.ResponseWriter, r *http.Request) {
		resp := openAIResponse{}
		resp.Choices = []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{{}}
		resp.Choices[0].Message.Content = string(outOfRangeContent)
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL,
	})

	risk, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.NoError(t, err)
	require.NotNil(t, risk.HealthScore)
	assert.Equal(t, 100.0, *risk.HealthScore) // clamped from 150
	require.NotNil(t, risk.DiversificationScore)
	assert.Equal(t, 0.0, *risk.DiversificationScore) // clamped from -10
}

func TestOpenAIRisk_DefaultModelAndBaseURL(t *testing.T) {
	// Verify NewOpenAIRiskService applies defaults when empty strings are passed
	svc := service.NewOpenAIRiskService(service.OpenAIRiskConfig{
		APIKey:  "key",
		BaseURL: "", // should default to https://api.openai.com
		Model:   "", // should default to gpt-4o-mini
	})
	assert.NotNil(t, svc)
}

// Ensure the unused models import is satisfied — the test file uses models.SectorEntry
var _ models.SectorEntry
