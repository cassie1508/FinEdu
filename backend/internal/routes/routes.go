package routes

import (
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"finedu-backend/internal/config"
	"finedu-backend/internal/handlers"
	"finedu-backend/internal/middleware"
	"finedu-backend/internal/repository"
	"finedu-backend/internal/service"
)

func Register(r *gin.Engine, pool *pgxpool.Pool, jwks keyfunc.Keyfunc, cfg config.Config) {
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", middleware.RateLimit(8), handlers.SignUp(cfg))
		}

		marketData := service.NewCachedMarketData(
			service.NewAlphaVantageMarketData(cfg.AlphaVantageAPIKey),
			1*time.Minute, // intraday candles (unused today — no free-tier intraday source)
			6*time.Hour,   // daily candles (1W/1M ranges)
			24*time.Hour,  // weekly candles (6M/1Y/5Y ranges)
		)
		chartsHandler := handlers.NewChartsHandler(marketData)

		companyRepo := repository.NewCompanyRepository(pool)
		companyHandler := handlers.NewCompanyHandler(companyRepo)

		companies := api.Group("/companies")
		{
			companies.GET("", companyHandler.SearchCompanies)
			companies.GET("/:symbol", companyHandler.GetCompany)
			companies.GET("/:symbol/prices", chartsHandler.GetPriceHistory)
		}

		newsProvider := service.NewCachedNewsProvider(
			service.NewFinnhubNewsProvider(cfg.FinnhubAPIKey),
			15*time.Minute, // Finnhub's free tier (60 req/min) affords a shorter TTL than Alpha Vantage's did
		)
		newsSummaries := service.NewNewsSummaryService(
			newsProvider,
			service.NewGeminiSummarizer(cfg.GeminiAPIKey),
			4*time.Hour,  // regenerate the AI summary at most this often
			48*time.Hour, // consider cached articles from this far back "today's news"
			30,           // cap articles sent to Gemini per summary
		)
		newsHandler := handlers.NewNewsHandler(newsProvider, newsSummaries)

		news := api.Group("/news")
		{
			news.GET("/general", newsHandler.GetGeneralNews)
			news.GET("/:ticker", newsHandler.GetTickerNews)
			news.GET("/:ticker/summary", newsHandler.GetTickerNewsSummary)
		}

		learning := api.Group("/learning")
		{
			learning.POST("/chat", handlers.ChatWithAI)
			learning.GET("/flashcards", handlers.GetFlashcards)
		}

		portfolio := api.Group("/portfolio")
		portfolio.Use(middleware.RequireAuth(jwks))
		{
			portfolio.GET("/holdings", handlers.ListHoldings)
			portfolio.POST("/holdings", handlers.AddHolding)
			portfolio.DELETE("/holdings/:id", handlers.RemoveHolding)
			portfolio.GET("/risk", handlers.GetPortfolioRisk)
		}
	}
}
