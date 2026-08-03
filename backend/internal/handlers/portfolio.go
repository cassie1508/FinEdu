package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"finedu-backend/internal/models"
	"finedu-backend/internal/service"
)

type PortfolioHandler struct {
	svc service.PortfolioServicer
}

func NewPortfolioHandler(svc service.PortfolioServicer) *PortfolioHandler {
	return &PortfolioHandler{svc: svc}
}

// CreatePortfolio handles POST /api/v1/portfolio/portfolios
func (h *PortfolioHandler) CreatePortfolio(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	portfolio, err := h.svc.CreatePortfolio(c.Request.Context(), userID, req.Name, req.Description)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

// ListPortfolios handles GET /api/v1/portfolio/portfolios
func (h *PortfolioHandler) ListPortfolios(c *gin.Context) {
	userID := c.GetString("userID")

	portfolios, err := h.svc.ListPortfolios(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list portfolios"})
		return
	}

	if portfolios == nil {
		portfolios = []models.PortfolioListItem{}
	}

	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})
}

// GetPortfolioDetail handles GET /api/v1/portfolio/portfolios/:portfolioId
func (h *PortfolioHandler) GetPortfolioDetail(c *gin.Context) {
	userID := c.GetString("userID")
	portfolioID := c.Param("portfolioId")

	detail, err := h.svc.GetPortfolioDetail(c.Request.Context(), userID, portfolioID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, detail)
}

// DeletePortfolio handles DELETE /api/v1/portfolio/portfolios/:portfolioId
func (h *PortfolioHandler) DeletePortfolio(c *gin.Context) {
	userID := c.GetString("userID")
	portfolioID := c.Param("portfolioId")

	if err := h.svc.DeletePortfolio(c.Request.Context(), userID, portfolioID); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// AddHolding handles POST /api/v1/portfolio/portfolios/:portfolioId/holdings
func (h *PortfolioHandler) AddHolding(c *gin.Context) {
	userID := c.GetString("userID")
	portfolioID := c.Param("portfolioId")

	var req models.AddHoldingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	holding, err := h.svc.AddHolding(c.Request.Context(), userID, portfolioID, req.Symbol, req.Shares, req.AverageCost)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, holding)
}

// UpdateHolding handles PUT /api/v1/portfolio/portfolios/:portfolioId/holdings/:holdingId
func (h *PortfolioHandler) UpdateHolding(c *gin.Context) {
	userID := c.GetString("userID")
	portfolioID := c.Param("portfolioId")
	holdingID := c.Param("holdingId")

	var req models.UpdateHoldingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	holding, err := h.svc.UpdateHolding(c.Request.Context(), userID, portfolioID, holdingID, req.Shares, req.AverageCost)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, holding)
}

// RemoveHolding handles DELETE /api/v1/portfolio/portfolios/:portfolioId/holdings/:holdingId
func (h *PortfolioHandler) RemoveHolding(c *gin.Context) {
	userID := c.GetString("userID")
	portfolioID := c.Param("portfolioId")
	holdingID := c.Param("holdingId")

	if err := h.svc.RemoveHolding(c.Request.Context(), userID, portfolioID, holdingID); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPortfolioSummary handles GET /api/v1/portfolio/portfolios/:portfolioId/summary
func (h *PortfolioHandler) GetPortfolioSummary(c *gin.Context) {
	userID := c.GetString("userID")
	portfolioID := c.Param("portfolioId")

	summary, err := h.svc.GetPortfolioSummary(c.Request.Context(), userID, portfolioID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetPortfolioRisk handles GET /api/v1/portfolio/portfolios/:portfolioId/risk
func (h *PortfolioHandler) GetPortfolioRisk(c *gin.Context) {
	userID := c.GetString("userID")
	portfolioID := c.Param("portfolioId")

	risk, err := h.svc.GetPortfolioRisk(c.Request.Context(), userID, portfolioID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, risk)
}

func (h *PortfolioHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
	case errors.Is(err, service.ErrDuplicate):
		c.JSON(http.StatusConflict, gin.H{"error": "resource already exists"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
