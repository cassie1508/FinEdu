package routes

import (
	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"

	"finedu-backend/internal/config"
	"finedu-backend/internal/handlers"
	"finedu-backend/internal/middleware"
)

func Register(r *gin.Engine, jwks keyfunc.Keyfunc, cfg config.Config) {
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", middleware.RateLimit(8), handlers.SignUp(cfg))
		}

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
		portfolio.Use(middleware.RequireAuth(jwks))
		{
			portfolio.GET("/holdings", handlers.ListHoldings)
			portfolio.POST("/holdings", handlers.AddHolding)
			portfolio.DELETE("/holdings/:id", handlers.RemoveHolding)
			portfolio.GET("/risk", handlers.GetPortfolioRisk)
		}
	}
}
