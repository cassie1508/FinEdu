package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCompanyNews returns AI-summarized news for a ticker symbol.
// Owner: Nhi (News Aggregation Summary + Interactive Stock Charts)
func GetCompanyNews(c *gin.Context) {
	symbol := c.Param("symbol")
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "GetCompanyNews not yet implemented",
		"symbol":  symbol,
	})
}

// GetPriceHistory returns historical price points for a ticker symbol and time range.
// Owner: Nhi (News Aggregation Summary + Interactive Stock Charts)
func GetPriceHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	rangeParam := c.DefaultQuery("range", "1M")
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "GetPriceHistory not yet implemented",
		"symbol":  symbol,
		"range":   rangeParam,
	})
}
