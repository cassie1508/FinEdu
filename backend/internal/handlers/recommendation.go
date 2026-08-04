package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetRecommendationsRequest represents a recommendation request.
type GetRecommendationsRequest struct {
	UserID             string `json:"user_id"`
	MaxResultsPerTitle int    `json:"max_results_per_title"`
}

// RecommendationResource represents a single resource.
type RecommendationResource struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

// RecommendationResponse represents recommendation results.
type RecommendationResponse struct {
	FlashcardTitle string                   `json:"flashcard_title"`
	Resources      []RecommendationResource `json:"resources"`
}

// GetRecommendations retrieves AI recommendations in real-time based on user's flashcard titles.
// POST /api/recommendations
func GetRecommendations(pool *pgxpool.Pool, defaultUserID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetRecommendationsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Use default user if not provided
		if req.UserID == "" {
			req.UserID = defaultUserID
		}

		if req.MaxResultsPerTitle == 0 {
			req.MaxResultsPerTitle = 3
		}

		// 1. Fetch all flashcard titles for the user
		titles, err := getFlashcardTitles(c.Request.Context(), pool, req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch flashcard titles: %v", err)})
			return
		}

		if len(titles) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"success":         true,
				"user_id":         req.UserID,
				"titles_count":    0,
				"recommendations": []interface{}{},
			})
			return
		}

		// 2. Call Python module to get recommendations from Tavily
		recommendations, err := callPythonRecommendations(titles, req.MaxResultsPerTitle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get recommendations: %v", err)})
			return
		}

		// 3. Return real-time results
		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"user_id":         req.UserID,
			"titles_count":    len(titles),
			"recommendations": recommendations,
		})
	}
}

// getFlashcardTitles fetches all flashcard titles for a user from the database.
func getFlashcardTitles(ctx context.Context, pool *pgxpool.Pool, userID string) ([]string, error) {
	rows, err := pool.Query(ctx, "SELECT title FROM flashcards WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return titles, nil
}

// callPythonRecommendations calls the Python recommendation module to get Tavily results.
func callPythonRecommendations(titles []string, maxResults int) ([]RecommendationResponse, error) {
	titlesJSON, err := json.Marshal(titles)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal titles: %w", err)
	}

	// Call the Python script directly (path is relative to backend/ directory)
	cmd := exec.Command("py", "internal/AIRecommendation/resources_suggestion.py", string(titlesJSON), fmt.Sprintf("%d", maxResults))

	// Inherit all environment variables (includes TAVILY_API_KEY from .env)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("python execution failed: %w, output: %s", err, string(output))
	}

	// Remove any debug logs from output (stderr goes to output with CombinedOutput)
	// Only parse the last line which should be the JSON
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var jsonLine string
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" && lines[i][0] == '{' {
			jsonLine = lines[i]
			break
		}
	}

	if jsonLine == "" {
		return nil, fmt.Errorf("no valid JSON output from python, full output: %s", string(output))
	}

	// Parse Python output: { "title1": [...resources...], "title2": [...], ... }
	var result map[string][]RecommendationResource
	if err := json.Unmarshal([]byte(jsonLine), &result); err != nil {
		return nil, fmt.Errorf("failed to parse python output (line was: %s): %w", jsonLine, err)
	}

	// Convert to response format
	var recommendations []RecommendationResponse
	for _, title := range titles {
		if resources, ok := result[title]; ok {
			recommendations = append(recommendations, RecommendationResponse{
				FlashcardTitle: title,
				Resources:      resources,
			})
		} else {
			recommendations = append(recommendations, RecommendationResponse{
				FlashcardTitle: title,
				Resources:      []RecommendationResource{},
			})
		}
	}

	return recommendations, nil
}
