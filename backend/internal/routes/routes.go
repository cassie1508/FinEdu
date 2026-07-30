package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"finedu-backend/internal/config"
	"finedu-backend/internal/handlers"
	"finedu-backend/internal/repository"
	"finedu-backend/internal/service"
)

func Register(r *gin.Engine, pool *pgxpool.Pool, cfg config.Config) {
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api/v1")
	{
		companyRepo := repository.NewCompanyRepository(pool)
		companyDataService := service.NewCompanyDataService(cfg.MarketDataAPIKey)
		companyHandler := handlers.NewCompanyHandler(companyRepo, companyDataService)

		companies := api.Group("/companies")
		{
			companies.GET("", companyHandler.SearchCompanies)
			companies.GET("/:symbol", companyHandler.GetCompany)
			companies.GET("/:symbol/news", handlers.GetCompanyNews)
			companies.GET("/:symbol/prices", handlers.GetPriceHistory)
		}

		learning := api.Group("/learning")
		{
			learning.POST("/chat", handlers.ChatWithAI)
			learning.GET("/flashcards", handlers.GetFlashcards)
		}

		portfolio := api.Group("/portfolio")
		{
			portfolio.GET("/holdings", handlers.ListHoldings)
			portfolio.POST("/holdings", handlers.AddHolding)
			portfolio.DELETE("/holdings/:id", handlers.RemoveHolding)
			portfolio.GET("/risk", handlers.GetPortfolioRisk)
		}
	}
}
