package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"finedu-backend/internal/models"
)

type FlashcardRepository struct {
	pool *pgxpool.Pool
}

func NewFlashcardRepository(pool *pgxpool.Pool) *FlashcardRepository {
	return &FlashcardRepository{pool: pool}
}

// GetAllFlashcards retrieves all flashcards for a user
func (r *FlashcardRepository) GetAllFlashcards(ctx context.Context, userID string) ([]models.Flashcard, error) {
	query := `
		SELECT id, user_id, title, category, why_it_matters, definition, example, 
		       common_misconceptions, review_count, created_at, updated_at
		FROM flashcards 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flashcards []models.Flashcard
	for rows.Next() {
		var card models.Flashcard

		err := rows.Scan(
			&card.ID,
			nil, // user_id (not needed in response)
			&card.Title,
			&card.Category,
			&card.WhyItMatters,
			&card.Definition,
			&card.Example,
			&card.CommonMisconception,
			&card.ReviewCount,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		flashcards = append(flashcards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return flashcards, nil
}

// GetFlashcardByID retrieves a specific flashcard by ID
func (r *FlashcardRepository) GetFlashcardByID(ctx context.Context, id, userID string) (models.Flashcard, error) {
	query := `
		SELECT id, user_id, title, category, why_it_matters, definition, example, 
		       common_misconceptions, review_count, created_at, updated_at
		FROM flashcards 
		WHERE id = $1 AND user_id = $2
	`

	var card models.Flashcard

	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&card.ID,
		nil, // user_id
		&card.Title,
		&card.Category,
		&card.WhyItMatters,
		&card.Definition,
		&card.Example,
		&card.CommonMisconception,
		&card.ReviewCount,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Flashcard{}, errors.New("flashcard not found")
		}
		return models.Flashcard{}, err
	}

	return card, nil
}

// CreateFlashcard inserts a new flashcard
func (r *FlashcardRepository) CreateFlashcard(ctx context.Context, userID string, req models.CreateFlashcardRequest) (models.Flashcard, error) {
	query := `
		INSERT INTO flashcards (user_id, title, category, why_it_matters, definition, example, 
		                        common_misconceptions, review_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 0, NOW(), NOW())
		RETURNING id, title, category, why_it_matters, definition, example, 
		          common_misconceptions, review_count, created_at, updated_at
	`

	var card models.Flashcard

	err := r.pool.QueryRow(
		ctx, query,
		userID,
		req.Title,
		req.Category,
		req.WhyItMatters,
		req.Definition,
		req.Example,
		req.CommonMisconception,
	).Scan(
		&card.ID,
		&card.Title,
		&card.Category,
		&card.WhyItMatters,
		&card.Definition,
		&card.Example,
		&card.CommonMisconception,
		&card.ReviewCount,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		return models.Flashcard{}, err
	}

	return card, nil
}

// UpdateFlashcard updates an existing flashcard
func (r *FlashcardRepository) UpdateFlashcard(ctx context.Context, id, userID string, req models.UpdateFlashcardRequest) (models.Flashcard, error) {
	query := `
		UPDATE flashcards 
		SET title = $1, category = $2, why_it_matters = $3, definition = $4, example = $5, 
		    common_misconceptions = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING id, title, category, why_it_matters, definition, example, 
		          common_misconceptions, review_count, created_at, updated_at
	`

	var card models.Flashcard

	err := r.pool.QueryRow(
		ctx, query,
		req.Title,
		req.Category,
		req.WhyItMatters,
		req.Definition,
		req.Example,
		req.CommonMisconception,
		id,
		userID,
	).Scan(
		&card.ID,
		&card.Title,
		&card.Category,
		&card.WhyItMatters,
		&card.Definition,
		&card.Example,
		&card.CommonMisconception,
		&card.ReviewCount,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Flashcard{}, errors.New("flashcard not found")
		}
		return models.Flashcard{}, err
	}

	return card, nil
}

// DeleteFlashcard deletes a flashcard
func (r *FlashcardRepository) DeleteFlashcard(ctx context.Context, id, userID string) error {
	query := `
		DELETE FROM flashcards 
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("flashcard not found")
	}

	return nil
}

// ReviewFlashcard increments the review count for a flashcard
func (r *FlashcardRepository) ReviewFlashcard(ctx context.Context, id, userID string) (models.Flashcard, error) {
	query := `
		UPDATE flashcards 
		SET review_count = review_count + 1, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING id, title, category, why_it_matters, definition, example, 
		          common_misconceptions, review_count, created_at, updated_at
	`

	var card models.Flashcard

	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&card.ID,
		&card.Title,
		&card.Category,
		&card.WhyItMatters,
		&card.Definition,
		&card.Example,
		&card.CommonMisconception,
		&card.ReviewCount,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Flashcard{}, errors.New("flashcard not found")
		}
		return models.Flashcard{}, err
	}

	return card, nil
}
