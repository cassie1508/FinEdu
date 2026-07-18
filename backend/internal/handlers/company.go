package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"finedu-backend/internal/models"
	"finedu-backend/internal/repository"
)

type CompanyHandler struct {
	repo *repository.CompanyRepository
}

func NewCompanyHandler(repo *repository.CompanyRepository) *CompanyHandler {
	return &CompanyHandler{repo: repo}
}

// GetCompany returns market dashboard data for a single ticker symbol.
// Owner: Nhien (Company Search & Market Dashboard)
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	symbol := c.Param("symbol")
	company, err := h.repo.GetBySymbol(c.Request.Context(), symbol)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, company)
}

// SearchCompanies returns companies matching a search query.
// Owner: Nhien (Company Search & Market Dashboard)
func (h *CompanyHandler) SearchCompanies(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, gin.H{"companies": []models.Company{}})
		return
	}
	companies, err := h.repo.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if companies == nil {
		companies = []models.Company{}
	}
	c.JSON(http.StatusOK, gin.H{"companies": companies})
}
