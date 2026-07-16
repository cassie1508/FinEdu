package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"finedu-backend/internal/config"
	"finedu-backend/internal/handlers"
	"finedu-backend/internal/middleware"
	"finedu-backend/internal/repository"
	"finedu-backend/internal/service"
)

func Register(r *gin.Engine, pool *pgxpool.Pool, cfg config.Config) {
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api/v1")
	{
		companies := api.Group("/companies")
		{
			companies.GET("", handlers.SearchCompanies)
			companies.GET("/:symbol", handlers.GetCompany)
			companies.GET("/:symbol/news", handlers.GetCompanyNews)
			companies.GET("/:symbol/prices", handlers.GetPriceHistory)
		}

		learning := api.Group("/learning")
		{
			learning.POST("/chat", handlers.ChatWithAI)
			learning.GET("/flashcards", handlers.GetFlashcards)
		}

		// Portfolio routes — require authentication
		portfolioRepo := repository.NewPortfolioRepository(pool)
		marketData := service.NewCachedMarketData(
			service.NewFinnhubMarketData(cfg.MarketDataAPIKey),
			60*time.Second,
		)
		portfolioSvc := service.NewPortfolioService(portfolioRepo, marketData)
		portfolioHandler := handlers.NewPortfolioHandler(portfolioSvc)

		portfolio := api.Group("/portfolio")
		portfolio.Use(middleware.RequireAuth(cfg.SupabaseJWTSecret))
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
