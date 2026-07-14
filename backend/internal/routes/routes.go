package routes

import (
	"github.com/gin-gonic/gin"

	"finedu-backend/internal/handlers"
)

func Register(r *gin.Engine) {
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

		portfolio := api.Group("/portfolio")
		{
			portfolio.GET("/holdings", handlers.ListHoldings)
			portfolio.POST("/holdings", handlers.AddHolding)
			portfolio.DELETE("/holdings/:id", handlers.RemoveHolding)
			portfolio.GET("/risk", handlers.GetPortfolioRisk)
		}
	}
}
