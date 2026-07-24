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
			learning.GET("/flashcards/:id", handlers.GetFlashcardByID)
			learning.POST("/flashcards", handlers.CreateFlashcard)
			learning.PUT("/flashcards/:id", handlers.UpdateFlashcard)
			learning.DELETE("/flashcards/:id", handlers.DeleteFlashcard)
			learning.POST("/flashcards/:id/review", handlers.ReviewFlashcard)

		}

		learningCenter := api.Group("/learning_center")
		{
			learningCenter.GET("/resources", handlers.GetLearningResources)
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
