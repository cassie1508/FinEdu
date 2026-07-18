package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"finedu-backend/internal/models"
)

type CompanyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{pool: pool}
}

func (r *CompanyRepository) GetBySymbol(ctx context.Context, symbol string) (*models.Company, error) {
	var c models.Company
	err := r.pool.QueryRow(ctx,
		`SELECT symbol, name, sector, industry, market_cap, revenue, eps,
		        pe_ratio, dividend_yield, week_high_52, week_low_52
		 FROM companies WHERE symbol = $1`,
		symbol,
	).Scan(&c.Symbol, &c.Name, &c.Sector, &c.Industry, &c.MarketCap,
		&c.Revenue, &c.EPS, &c.PERatio, &c.DividendYield, &c.WeekHigh52, &c.WeekLow52)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *CompanyRepository) Search(ctx context.Context, query string) ([]models.Company, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT symbol, name, sector, industry, market_cap, revenue, eps,
		        pe_ratio, dividend_yield, week_high_52, week_low_52
		 FROM companies
		 WHERE symbol ILIKE $1 OR name ILIKE $1
		 LIMIT 10`,
		"%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Company
	for rows.Next() {
		var c models.Company
		if err := rows.Scan(&c.Symbol, &c.Name, &c.Sector, &c.Industry, &c.MarketCap,
			&c.Revenue, &c.EPS, &c.PERatio, &c.DividendYield, &c.WeekHigh52, &c.WeekLow52); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil
}