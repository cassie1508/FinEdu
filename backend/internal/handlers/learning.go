package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"finedu-backend/internal/models"
	"finedu-backend/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	finnhubBaseURL = "https://finnhub.io/api/v1"
)

var finnhubHTTPClient = &http.Client{Timeout: 10 * time.Second}

var (
	flashcardService *services.FlashcardService
	defaultUserID    string
)

// InitFlashcards wires the flashcard service into the handlers.
// TODO: replace defaultUserID with the authenticated user's id once Supabase Auth is wired in.
func InitFlashcards(service *services.FlashcardService, userID string) {
	flashcardService = service
	defaultUserID = userID
}

// ChatWithAI handles RAG-based finance chatbot conversations.
// Owner: Hiếu (AI Learning Center)
func ChatWithAI(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "ChatWithAI not yet implemented"})
}
func getFinnhubAPIKey() string {
	apiKey := strings.TrimSpace(os.Getenv("FINNHUB_API_KEY"))
	if apiKey != "" {
		return apiKey
	}

	return strings.TrimSpace(os.Getenv("MARKET_DATA_API_KEY"))
}

// financeNewsCategories are the Finnhub /news categories that are relevant to
// personal finance and markets. "general" alone can include unrelated tech/world
// news, so we merge it with the more market-focused categories.
var financeNewsCategories = []string{"general", "forex", "crypto", "merger", "investment", "earnings", "economic", "commodities", "stock", "portfolio", "bond", "dividend", "funds", "nasdaq", "nyse", "sp500"}

// GetLearningResources returns learning resources from Finnhub market news.
// Endpoint: GET /api/v1/learning_center/resources
func GetLearningResources(c *gin.Context) {
	apiKey := getFinnhubAPIKey()
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "FINNHUB_API_KEY (or MARKET_DATA_API_KEY) is not configured"})
		return
	}

	seen := make(map[int64]bool)
	resources := make([]models.LearningResource, 0)

	for _, category := range financeNewsCategories {
		items, err := fetchFinnhubNews(c, apiKey, category, 0)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch learning resources", "details": err.Error()})
			return
		}

		for _, item := range items {
			if seen[item.ID] {
				continue
			}
			seen[item.ID] = true
			resources = append(resources, mapNewsItemToResource(item))
		}
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].PublishedAt.After(resources[j].PublishedAt)
	})

	c.JSON(http.StatusOK, gin.H{"data": resources})
}

func fetchFinnhubNews(c *gin.Context, apiKey, category string, minID int64) ([]models.FinnhubNewsItem, error) {
	endpoint, err := url.Parse(finnhubBaseURL + "/news")
	if err != nil {
		return nil, fmt.Errorf("build request URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("category", category)
	query.Set("token", apiKey)
	if minID > 0 {
		query.Set("minId", strconv.FormatInt(minID, 10))
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := finnhubHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call finnhub: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read finnhub response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("finnhub returned status %d: %s", resp.StatusCode, string(body))
	}

	var items []models.FinnhubNewsItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("decode finnhub response: %w", err)
	}

	return items, nil
}

func mapNewsItemToResource(item models.FinnhubNewsItem) models.LearningResource {
	resource := models.LearningResource{
		ID:          item.ID,
		Title:       item.Headline,
		Category:    item.Category,
		Summary:     item.Summary,
		Source:      item.Source,
		ImageURL:    item.Image,
		Related:     splitRelatedSymbols(item.Related),
		PublishedAt: time.Unix(item.DateTime, 0).UTC(),
		URL:         item.URL,
	}

	return resource
}

func splitRelatedSymbols(related string) []string {
	related = strings.TrimSpace(related)
	if related == "" {
		return []string{}
	}

	rawParts := strings.Split(related, ",")
	parts := make([]string, 0, len(rawParts))
	for _, p := range rawParts {
		t := strings.TrimSpace(p)
		if t == "" {
			continue
		}
		parts = append(parts, t)
	}

	return parts
}

// GetFlashcards returns flashcards for a given financial topic.
// Owner: Hiếu (AI Learning Center)
func GetFlashcards(c *gin.Context) {
	flashcards, err := flashcardService.GetAllFlashcards(c.Request.Context(), defaultUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch flashcards", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": flashcards})
}

// GetFlashcardByID returns a single flashcard.
func GetFlashcardByID(c *gin.Context) {
	id := c.Param("id")
	flashcard, err := flashcardService.GetFlashcardByID(c.Request.Context(), id, defaultUserID)
	if errors.Is(err, services.ErrFlashcardNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "flashcard not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch flashcard", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": flashcard})
}

// CreateFlashcard creates a new flashcard.
func CreateFlashcard(c *gin.Context) {
	var req models.CreateFlashcardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	flashcard, err := flashcardService.CreateFlashcard(c.Request.Context(), defaultUserID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create flashcard", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": flashcard})
}

// UpdateFlashcard updates an existing flashcard.
func UpdateFlashcard(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateFlashcardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	flashcard, err := flashcardService.UpdateFlashcard(c.Request.Context(), id, defaultUserID, req)
	if errors.Is(err, services.ErrFlashcardNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "flashcard not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update flashcard", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": flashcard})
}

// DeleteFlashcard deletes a flashcard by id.
func DeleteFlashcard(c *gin.Context) {
	id := c.Param("id")
	err := flashcardService.DeleteFlashcard(c.Request.Context(), id, defaultUserID)
	if errors.Is(err, services.ErrFlashcardNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "flashcard not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete flashcard", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "flashcard deleted"})
}

// ReviewFlashcard increments the review count for a flashcard.
func ReviewFlashcard(c *gin.Context) {
	id := c.Param("id")
	flashcard, err := flashcardService.ReviewFlashcard(c.Request.Context(), id, defaultUserID)
	if errors.Is(err, services.ErrFlashcardNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "flashcard not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to review flashcard", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": flashcard})
}
