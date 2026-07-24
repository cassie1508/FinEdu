package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"finedu-backend/internal/db"
	"finedu-backend/internal/models"
)

var (
	ErrFlashcardNotFound = errors.New("flashcard not found")
	ErrInvalidTopic      = errors.New("topic is required")
	ErrInvalidCount      = errors.New("count must be greater than or equal to 0")
)

func GetAllFlashcards() []models.Flashcard {
	db.FlashcardStoreMu.Lock()
	defer db.FlashcardStoreMu.Unlock()

	cards := make([]models.Flashcard, len(db.FlashcardStore))
	copy(cards, db.FlashcardStore)
	return cards
}

func GetFlashcardByID(id string) (models.Flashcard, error) {
	db.FlashcardStoreMu.Lock()
	defer db.FlashcardStoreMu.Unlock()

	for _, flashcard := range db.FlashcardStore {
		if flashcard.ID == id {
			return flashcard, nil
		}
	}

	return models.Flashcard{}, ErrFlashcardNotFound
}

func CreateFlashcard(req models.CreateFlashcardRequest) models.Flashcard {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.WhyItMatters = strings.TrimSpace(req.WhyItMatters)
	req.Definition = strings.TrimSpace(req.Definition)
	req.Example = strings.TrimSpace(req.Example)

	db.FlashcardStoreMu.Lock()
	defer db.FlashcardStoreMu.Unlock()

	now := time.Now().UTC()
	card := models.Flashcard{
		ID:                  fmt.Sprintf("fc-%d", db.FlashcardSeq),
		CreatedAt:           now,
		UpdatedAt:           now,
		Title:               req.Title,
		Category:            req.Category,
		WhyItMatters:        req.WhyItMatters,
		Definition:          req.Definition,
		Example:             req.Example,
		CommonMisconception: req.CommonMisconception,
		ReviewCount:         0,
	}
	db.FlashcardSeq++
	db.FlashcardStore = append(db.FlashcardStore, card)

	return card
}

func UpdateFlashcard(id string, req models.UpdateFlashcardRequest) (models.Flashcard, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Category = strings.TrimSpace(req.Category)
	req.WhyItMatters = strings.TrimSpace(req.WhyItMatters)
	req.Definition = strings.TrimSpace(req.Definition)
	req.Example = strings.TrimSpace(req.Example)

	db.FlashcardStoreMu.Lock()
	defer db.FlashcardStoreMu.Unlock()

	for idx, card := range db.FlashcardStore {
		if card.ID != id {
			continue
		}

		card.Title = req.Title
		card.Category = req.Category
		card.WhyItMatters = req.WhyItMatters
		card.Definition = req.Definition
		card.Example = req.Example
		card.CommonMisconception = req.CommonMisconception
		card.UpdatedAt = time.Now().UTC()
		db.FlashcardStore[idx] = card
		return card, nil
	}

	return models.Flashcard{}, ErrFlashcardNotFound
}

func DeleteFlashcard(id string) error {
	db.FlashcardStoreMu.Lock()
	defer db.FlashcardStoreMu.Unlock()

	for idx, card := range db.FlashcardStore {
		if card.ID != id {
			continue
		}

		db.FlashcardStore = append(db.FlashcardStore[:idx], db.FlashcardStore[idx+1:]...)
		return nil
	}

	return ErrFlashcardNotFound
}

func ReviewFlashcard(id string) (models.Flashcard, error) {
	db.FlashcardStoreMu.Lock()
	defer db.FlashcardStoreMu.Unlock()

	for idx, card := range db.FlashcardStore {
		if card.ID != id {
			continue
		}

		card.ReviewCount++
		card.UpdatedAt = time.Now().UTC()
		db.FlashcardStore[idx] = card
		return card, nil
	}

	return models.Flashcard{}, ErrFlashcardNotFound
}
