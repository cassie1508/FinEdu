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

func Register(r *gin.Engine, pool *pgxpool.Pool, jwks keyfunc.Keyfunc, cfg config.Config, defaultUserID string) {
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
			learning.GET("/flashcards/:id", handlers.GetFlashcardByID)
			learning.POST("/flashcards", handlers.CreateFlashcard)
			learning.PUT("/flashcards/:id", handlers.UpdateFlashcard)
			learning.DELETE("/flashcards/:id", handlers.DeleteFlashcard)
			learning.POST("/flashcards/:id/review", handlers.ReviewFlashcard)
			learning.GET("/resources/podcast", handlers.GetPodcastByListennotes)

		}

		learningCenter := api.Group("/learning_center")
		{
			learningCenter.GET("/resources", handlers.GetLearningResources)
		}

		// RAG Pipeline endpoints
		rag := api.Group("/rag")
		{
			ragHandler := handlers.NewRAGHandler(pool, defaultUserID)
			rag.POST("/query", ragHandler.QueryRAG)
		}

		// AI Recommendation endpoints
		recommendations := api.Group("/recommendations")
		{
			recommendations.POST("", handlers.GetRecommendations(pool, defaultUserID))
		}

		documents := api.Group("/documents")
		{
			ragHandler := handlers.NewRAGHandler(pool, defaultUserID)
			documents.POST("/upload", ragHandler.UploadDocument)
		}

		// Portfolio routes — require authentication
		portfolioRepo := repository.NewPortfolioRepository(pool)
		quoteData := service.NewCachedQuoteData(
			service.NewFinnhubMarketData(cfg.FinnhubAPIKey),
			60*time.Second,
		)
		// Uses Gemini (already required for news summaries) rather than OpenAI
		// for now, so portfolio risk analysis doesn't need a separate paid key.
		aiRiskSvc := service.NewGeminiRiskService(service.GeminiRiskConfig{
			APIKey: cfg.GeminiAPIKey,
		})
		portfolioSvc := service.NewPortfolioService(portfolioRepo, quoteData, aiRiskSvc)
		portfolioHandler := handlers.NewPortfolioHandler(portfolioSvc)

		portfolio := api.Group("/portfolio")
		portfolio.Use(middleware.RequireAuth(jwks))
		{
			portfolio.POST("/portfolios", portfolioHandler.CreatePortfolio)
			portfolio.GET("/portfolios", portfolioHandler.ListPortfolios)
			portfolio.GET("/portfolios/:portfolioId", portfolioHandler.GetPortfolioDetail)
			portfolio.DELETE("/portfolios/:portfolioId", portfolioHandler.DeletePortfolio)

			portfolio.POST("/portfolios/:portfolioId/holdings", portfolioHandler.AddHolding)
			portfolio.PUT("/portfolios/:portfolioId/holdings/:holdingId", portfolioHandler.UpdateHolding)
			portfolio.DELETE("/portfolios/:portfolioId/holdings/:holdingId", portfolioHandler.RemoveHolding)

			portfolio.GET("/portfolios/:portfolioId/summary", portfolioHandler.GetPortfolioSummary)
			portfolio.GET("/portfolios/:portfolioId/risk", portfolioHandler.GetPortfolioRisk)
		}
	}
}
