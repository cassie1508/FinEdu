package services

import (
	"context"
	"strings"

	"finedu-backend/internal/db"
	"finedu-backend/internal/models"
)

// ErrFlashcardNotFound is returned when a flashcard cannot be found for the given id/user.
var ErrFlashcardNotFound = db.ErrFlashcardNotFound

// FlashcardService contains the business logic for flashcards and delegates
// persistence to a FlashcardRepository.
type FlashcardService struct {
	repo *db.FlashcardRepository
}

// NewFlashcardService creates a FlashcardService backed by the given repository.
func NewFlashcardService(repo *db.FlashcardRepository) *FlashcardService {
	return &FlashcardService{repo: repo}
}

// GetAllFlashcards returns every flashcard owned by userID.
func (s *FlashcardService) GetAllFlashcards(ctx context.Context, userID string) ([]models.Flashcard, error) {
	return s.repo.GetAllFlashcards(ctx, userID)
}

// GetFlashcardByID returns a single flashcard owned by userID.
func (s *FlashcardService) GetFlashcardByID(ctx context.Context, id, userID string) (models.Flashcard, error) {
	return s.repo.GetFlashcardByID(ctx, id, userID)
}

// CreateFlashcard validates/normalizes the request and creates a new flashcard.
func (s *FlashcardService) CreateFlashcard(ctx context.Context, userID string, req models.CreateFlashcardRequest) (models.Flashcard, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.WhyItMatters = strings.TrimSpace(req.WhyItMatters)
	req.Definition = strings.TrimSpace(req.Definition)
	req.Example = strings.TrimSpace(req.Example)

	return s.repo.CreateFlashcard(ctx, userID, req)
}

// UpdateFlashcard validates/normalizes the request and updates an existing flashcard.
func (s *FlashcardService) UpdateFlashcard(ctx context.Context, id, userID string, req models.UpdateFlashcardRequest) (models.Flashcard, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.WhyItMatters = strings.TrimSpace(req.WhyItMatters)
	req.Definition = strings.TrimSpace(req.Definition)
	req.Example = strings.TrimSpace(req.Example)

	return s.repo.UpdateFlashcard(ctx, id, userID, req)
}

// DeleteFlashcard deletes a flashcard owned by userID.
func (s *FlashcardService) DeleteFlashcard(ctx context.Context, id, userID string) error {
	return s.repo.DeleteFlashcard(ctx, id, userID)
}

// ReviewFlashcard increments the review count for a flashcard owned by userID.
func (s *FlashcardService) ReviewFlashcard(ctx context.Context, id, userID string) (models.Flashcard, error) {
	return s.repo.ReviewFlashcard(ctx, id, userID)
}
