package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"finedu-backend/internal/models"
	"finedu-backend/internal/service"
)

// ChartsHandler serves historical price data for the Interactive Stock
// Charts feature.
// Owner: Nhi (News Aggregation Summary + Interactive Stock Charts)
type ChartsHandler struct {
	marketData service.MarketDataService
}

func NewChartsHandler(marketData service.MarketDataService) *ChartsHandler {
	return &ChartsHandler{marketData: marketData}
}

// marketDataDelayMinutes reflects the ~15min delay on free-tier market data.
const marketDataDelayMinutes = 15

// rangeToParams maps a UI time-range pill to an internal candle resolution
// (D, W) and lookback window.
//
// There's no free-tier intraday data source (Finnhub and Alpha Vantage both
// 403 it as premium-only), so there's no "1D" range — with only one bar per
// calendar day available, a true 1-day window would be a single point, and
// widening it to look like a chart would just duplicate what "1W" already
// shows.
//
// Daily candles are also capped at the latest ~100 points on the free tier
// (outputsize=full is premium-only), which covers 1W/1M but not 6M/1Y — so
// those two use weekly bars instead, which have no such cap.
func rangeToParams(rangeParam string) (resolution string, from, to time.Time, ok bool) {
	to = time.Now()
	switch rangeParam {
	case "1W":
		return "D", to.AddDate(0, 0, -14), to, true
	case "1M":
		return "D", to.AddDate(0, -1, 0), to, true
	case "6M":
		return "W", to.AddDate(0, -6, 0), to, true
	case "1Y":
		return "W", to.AddDate(-1, 0, 0), to, true
	case "5Y":
		return "W", to.AddDate(-5, 0, 0), to, true
	default:
		return "", time.Time{}, time.Time{}, false
	}
}

// GetPriceHistory returns historical price candles for a ticker symbol and
// time range (1W|1M|6M|1Y|5Y).
func (h *ChartsHandler) GetPriceHistory(c *gin.Context) {
	symbol := strings.ToUpper(c.Param("symbol"))
	rangeParam := c.DefaultQuery("range", "1M")

	resolution, from, to, ok := rangeToParams(rangeParam)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid range: expected one of 1W, 1M, 6M, 1Y, 5Y",
		})
		return
	}

	candles, err := h.marketData.GetCandles(c.Request.Context(), symbol, resolution, from, to)
	if err != nil {
		log.Printf("GetPriceHistory: %v", err)

		if errors.Is(err, service.ErrUpstreamUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "market data provider is temporarily unavailable",
			})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to fetch price history",
		})
		return
	}

	c.JSON(http.StatusOK, models.PriceHistory{
		Symbol:         symbol,
		Range:          rangeParam,
		Resolution:     resolution,
		Candles:        candles,
		DelayedMinutes: marketDataDelayMinutes,
	})
}

// NewsHandler serves raw and AI-summarized news for the News Aggregation
// Summary feature.
// Owner: Nhi (News Aggregation Summary + Interactive Stock Charts)
type NewsHandler struct {
	articles  service.NewsProvider
	summaries service.SummaryService
}

func NewNewsHandler(articles service.NewsProvider, summaries service.SummaryService) *NewsHandler {
	return &NewsHandler{articles: articles, summaries: summaries}
}

const defaultNewsLimit = 20

// parseLimitQuery reads the `limit` query param, defaulting when absent.
func parseLimitQuery(c *gin.Context) (int, bool) {
	raw := c.Query("limit")
	if raw == "" {
		return defaultNewsLimit, true
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit: expected a positive integer"})
		return 0, false
	}
	return limit, true
}

// GetTickerNews returns raw news articles for a ticker.
// GET /api/v1/news/:ticker?limit=20&since=<date>
func (h *NewsHandler) GetTickerNews(c *gin.Context) {
	ticker := strings.ToUpper(c.Param("ticker"))

	limit, ok := parseLimitQuery(c)
	if !ok {
		return
	}

	var since time.Time
	if raw := c.Query("since"); raw != "" {
		parsed, err := time.Parse("2006-01-02", raw)
		if err != nil {
			parsed, err = time.Parse(time.RFC3339, raw)
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid since: expected YYYY-MM-DD or RFC3339"})
			return
		}
		since = parsed
	}

	articles, err := h.articles.GetTickerNews(c.Request.Context(), ticker)
	if err != nil {
		log.Printf("GetTickerNews: %v", err)
		if errors.Is(err, service.ErrUpstreamUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "news provider is temporarily unavailable"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch news"})
		return
	}

	c.JSON(http.StatusOK, service.FilterArticles(articles, since, limit))
}

// GetTickerNewsSummary returns the cached AI-generated daily digest for a
// ticker.
// GET /api/v1/news/:ticker/summary
func (h *NewsHandler) GetTickerNewsSummary(c *gin.Context) {
	ticker := strings.ToUpper(c.Param("ticker"))

	summary, err := h.summaries.GetSummary(c.Request.Context(), ticker)
	if err != nil {
		log.Printf("GetTickerNewsSummary: %v", err)
		if errors.Is(err, service.ErrNoArticles) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no recent news available to summarize for this ticker"})
			return
		}
		if errors.Is(err, service.ErrUpstreamUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "news summary provider is temporarily unavailable"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to generate news summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

var generalNewsCategories = map[string]bool{"market": true, "earnings": true, "mergers": true}

// GetGeneralNews returns market-wide news that isn't tied to a specific
// ticker.
// GET /api/v1/news/general?category=market|earnings|mergers&limit=20
func (h *NewsHandler) GetGeneralNews(c *gin.Context) {
	category := c.DefaultQuery("category", "market")
	if !generalNewsCategories[category] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category: expected one of market, earnings, mergers"})
		return
	}

	limit, ok := parseLimitQuery(c)
	if !ok {
		return
	}

	articles, err := h.articles.GetGeneralNews(c.Request.Context(), category)
	if err != nil {
		log.Printf("GetGeneralNews: %v", err)
		if errors.Is(err, service.ErrUpstreamUnavailable) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "news provider is temporarily unavailable"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch news"})
		return
	}

	c.JSON(http.StatusOK, service.FilterArticles(articles, time.Time{}, limit))
}
