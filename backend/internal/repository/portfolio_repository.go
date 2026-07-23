package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"finedu-backend/internal/models"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

// PortfolioRepo defines the persistence contract for portfolios and holdings.
type PortfolioRepo interface {
	CreatePortfolio(ctx context.Context, userID, name, description string) (*models.Portfolio, error)
	GetPortfolioByID(ctx context.Context, portfolioID string) (*models.Portfolio, error)
	ListPortfoliosByUser(ctx context.Context, userID string) ([]models.PortfolioListItem, error)
	DeletePortfolio(ctx context.Context, portfolioID string) error
	AddHolding(ctx context.Context, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error)
	GetHoldingByID(ctx context.Context, holdingID string) (*models.PortfolioHolding, error)
	GetHoldingsByPortfolio(ctx context.Context, portfolioID string) ([]models.PortfolioHolding, error)
	UpdateHolding(ctx context.Context, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error)
	DeleteHolding(ctx context.Context, holdingID string) error
}

type PortfolioRepository struct {
	pool *pgxpool.Pool
}

func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{pool: pool}
}

// CreatePortfolio inserts a new portfolio for the given user.
func (r *PortfolioRepository) CreatePortfolio(ctx context.Context, userID, name, description string) (*models.Portfolio, error) {
	var p models.Portfolio
	err := r.pool.QueryRow(ctx,
		`INSERT INTO portfolios (user_id, name, description)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, name, description, created_at, updated_at`,
		userID, name, description,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &p, nil
}

// GetPortfolioByID retrieves a single portfolio by its ID.
func (r *PortfolioRepository) GetPortfolioByID(ctx context.Context, portfolioID string) (*models.Portfolio, error) {
	var p models.Portfolio
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, description, created_at, updated_at
		 FROM portfolios WHERE id = $1`,
		portfolioID,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// ListPortfoliosByUser returns all portfolios belonging to a user with holdings count.
func (r *PortfolioRepository) ListPortfoliosByUser(ctx context.Context, userID string) ([]models.PortfolioListItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.id, p.name, p.description, COUNT(h.id) AS holdings_count, p.created_at
		 FROM portfolios p
		 LEFT JOIN portfolio_holdings h ON h.portfolio_id = p.id
		 WHERE p.user_id = $1
		 GROUP BY p.id
		 ORDER BY p.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PortfolioListItem
	for rows.Next() {
		var item models.PortfolioListItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.HoldingsCount, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// DeletePortfolio removes a portfolio (cascade deletes holdings).
func (r *PortfolioRepository) DeletePortfolio(ctx context.Context, portfolioID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM portfolios WHERE id = $1`, portfolioID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// AddHolding inserts a new holding into a portfolio.
func (r *PortfolioRepository) AddHolding(ctx context.Context, portfolioID, symbol string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	var h models.PortfolioHolding
	err := r.pool.QueryRow(ctx,
		`INSERT INTO portfolio_holdings (portfolio_id, symbol, shares, average_cost)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, portfolio_id, symbol, shares, average_cost, created_at, updated_at`,
		portfolioID, symbol, shares, averageCost,
	).Scan(&h.ID, &h.PortfolioID, &h.Symbol, &h.Shares, &h.AverageCost, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &h, nil
}

// GetHoldingByID retrieves a single holding.
func (r *PortfolioRepository) GetHoldingByID(ctx context.Context, holdingID string) (*models.PortfolioHolding, error) {
	var h models.PortfolioHolding
	err := r.pool.QueryRow(ctx,
		`SELECT id, portfolio_id, symbol, shares, average_cost, created_at, updated_at
		 FROM portfolio_holdings WHERE id = $1`,
		holdingID,
	).Scan(&h.ID, &h.PortfolioID, &h.Symbol, &h.Shares, &h.AverageCost, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &h, nil
}

// GetHoldingsByPortfolio returns all holdings for a given portfolio.
func (r *PortfolioRepository) GetHoldingsByPortfolio(ctx context.Context, portfolioID string) ([]models.PortfolioHolding, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, portfolio_id, symbol, shares, average_cost, created_at, updated_at
		 FROM portfolio_holdings WHERE portfolio_id = $1
		 ORDER BY created_at ASC`,
		portfolioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holdings []models.PortfolioHolding
	for rows.Next() {
		var h models.PortfolioHolding
		if err := rows.Scan(&h.ID, &h.PortfolioID, &h.Symbol, &h.Shares, &h.AverageCost, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		holdings = append(holdings, h)
	}
	return holdings, rows.Err()
}

// UpdateHolding modifies shares and average cost for an existing holding.
func (r *PortfolioRepository) UpdateHolding(ctx context.Context, holdingID string, shares, averageCost float64) (*models.PortfolioHolding, error) {
	var h models.PortfolioHolding
	err := r.pool.QueryRow(ctx,
		`UPDATE portfolio_holdings
		 SET shares = $2, average_cost = $3, updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, portfolio_id, symbol, shares, average_cost, created_at, updated_at`,
		holdingID, shares, averageCost,
	).Scan(&h.ID, &h.PortfolioID, &h.Symbol, &h.Shares, &h.AverageCost, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &h, nil
}

// DeleteHolding removes a holding by ID.
func (r *PortfolioRepository) DeleteHolding(ctx context.Context, holdingID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM portfolio_holdings WHERE id = $1`, holdingID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// isDuplicateError checks for PostgreSQL unique constraint violation (code 23505).
func isDuplicateError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23505")
}
