package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCompany returns market dashboard data for a single ticker symbol.
// Owner: Nhien (Company Search & Market Dashboard)
func GetCompany(c *gin.Context) {
	symbol := c.Param("symbol")
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "GetCompany not yet implemented",
		"symbol":  symbol,
	})
}

// SearchCompanies returns companies matching a search query.
// Owner: Nhien (Company Search & Market Dashboard)
func SearchCompanies(c *gin.Context) {
	query := c.Query("q")
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "SearchCompanies not yet implemented",
		"query":   query,
	})
}
