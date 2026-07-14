package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListHoldings returns all holdings in the authenticated user's portfolio.
// Owner: Quang (Portfolio Management + Risk)
func ListHoldings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "ListHoldings not yet implemented"})
}

// AddHolding adds a new holding to the authenticated user's portfolio.
// Owner: Quang (Portfolio Management + Risk)
func AddHolding(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "AddHolding not yet implemented"})
}

// RemoveHolding removes a holding from the authenticated user's portfolio.
// Owner: Quang (Portfolio Management + Risk)
func RemoveHolding(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "RemoveHolding not yet implemented",
		"id":      id,
	})
}

// GetPortfolioRisk returns the portfolio health score and risk breakdown.
// Owner: Quang (Portfolio Management + Risk)
func GetPortfolioRisk(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GetPortfolioRisk not yet implemented"})
}
