package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"finedu-backend/internal/models"
	"finedu-backend/internal/repository"
	"finedu-backend/internal/service"
)

type CompanyRepo interface {
	GetBySymbol(ctx context.Context, symbol string) (*models.Company, error)
	Search(ctx context.Context, query string) ([]models.Company, error)
	Upsert(ctx context.Context, c *models.Company) error
}

type CompanyHandler struct {
	repo        CompanyRepo
	dataService *service.CompanyDataService
}

func NewCompanyHandler(repo CompanyRepo, dataService *service.CompanyDataService) *CompanyHandler {
	return &CompanyHandler{repo: repo, dataService: dataService}
}

// GetCompany returns market dashboard data for a single ticker symbol.
// Owner: Nhien (Company Search & Market Dashboard)
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	symbol := c.Param("symbol")
	ctx := c.Request.Context()

	// 1. Thử lấy từ DB trước (cache)
	company, err := h.repo.GetBySymbol(ctx, symbol)
	if err == nil {
		c.JSON(http.StatusOK, company)
		return
	}
	if !errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 2. DB chưa có → gọi Finnhub
	fetched, ferr := h.dataService.FetchCompany(ctx, symbol)
	if ferr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	// 3. Lưu vào DB để lần sau đọc nhanh (cache)
	_ = h.repo.Upsert(ctx, fetched)

	c.JSON(http.StatusOK, fetched)
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
