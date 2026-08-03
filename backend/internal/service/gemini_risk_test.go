package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"finedu-backend/internal/service"
)

func validGeminiRiskContent() string {
	b, _ := json.Marshal(map[string]interface{}{
		"healthScore":          85,
		"riskLevel":            "Moderate",
		"sectorConcentration":  []map[string]interface{}{{"sector": "Technology", "percent": 100}},
		"diversificationScore": 40,
		"recommendations":      []string{"Diversify sectors", "Add bonds", "Consider international"},
	})
	return string(b)
}

func geminiRiskResponse(text string) map[string]any {
	return map[string]any{
		"candidates": []map[string]any{
			{"content": map[string]any{"parts": []map[string]any{{"text": text}}}},
		},
	}
}

func newGeminiRiskServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestGeminiRisk_SuccessfulAnalysis(t *testing.T) {
	srv := newGeminiRiskServer(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "test-key", r.Header.Get("x-goog-api-key"))
		json.NewEncoder(w).Encode(geminiRiskResponse(validGeminiRiskContent()))
	})
	defer srv.Close()

	svc := service.NewGeminiRiskService(service.GeminiRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL + "?model=%s",
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

func TestGeminiRisk_APIError500(t *testing.T) {
	srv := newGeminiRiskServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": {"message": "server error"}}`))
	})
	defer srv.Close()

	svc := service.NewGeminiRiskService(service.GeminiRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL + "?model=%s",
	})

	_, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Gemini API error")
}

func TestGeminiRisk_NoCandidates(t *testing.T) {
	srv := newGeminiRiskServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"candidates": []map[string]any{}})
	})
	defer srv.Close()

	svc := service.NewGeminiRiskService(service.GeminiRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL + "?model=%s",
	})

	_, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no candidates")
}

func TestGeminiRisk_MalformedAIContent(t *testing.T) {
	srv := newGeminiRiskServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(geminiRiskResponse("not valid json content"))
	})
	defer srv.Close()

	svc := service.NewGeminiRiskService(service.GeminiRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL + "?model=%s",
	})

	_, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing AI JSON output")
}

func TestGeminiRisk_ScoresClamped(t *testing.T) {
	outOfRangeContent, _ := json.Marshal(map[string]interface{}{
		"healthScore":          150,
		"riskLevel":            "Low",
		"sectorConcentration":  []interface{}{},
		"diversificationScore": -10,
		"recommendations":      []string{"Tip"},
	})

	srv := newGeminiRiskServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(geminiRiskResponse(string(outOfRangeContent)))
	})
	defer srv.Close()

	svc := service.NewGeminiRiskService(service.GeminiRiskConfig{
		APIKey:  "test-key",
		BaseURL: srv.URL + "?model=%s",
	})

	risk, err := svc.AnalyzeRisk(context.Background(), testHoldings())
	require.NoError(t, err)
	require.NotNil(t, risk.HealthScore)
	assert.Equal(t, 100.0, *risk.HealthScore) // clamped from 150
	require.NotNil(t, risk.DiversificationScore)
	assert.Equal(t, 0.0, *risk.DiversificationScore) // clamped from -10
}

func TestGeminiRisk_DefaultModelAndBaseURL(t *testing.T) {
	svc := service.NewGeminiRiskService(service.GeminiRiskConfig{
		APIKey:  "key",
		BaseURL: "", // should default to Gemini's generateContent endpoint
		Model:   "", // should default to gemini-2.0-flash
	})
	assert.NotNil(t, svc)
}
