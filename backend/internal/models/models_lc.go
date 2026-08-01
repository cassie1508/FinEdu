package models

import "time"

type Flashcard struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	Category            string    `json:"category"`
	WhyItMatters        string    `json:"whyItMatters"`
	Definition          string    `json:"definition"`
	Example             string    `json:"example"`
	CommonMisconception string    `json:"commonMisconceptions"`
	ReviewCount         int       `json:"reviewCount"`
	UpdatedAt           time.Time `json:"updatedAt"`
	CreatedAt           time.Time `json:"createdAt"`
}

type CreateFlashcardRequest struct {
	Title               string `json:"title" binding:"required"`
	Category            string `json:"category" binding:"required"`
	WhyItMatters        string `json:"whyItMatters" binding:"required"`
	Definition          string `json:"definition" binding:"required"`
	Example             string `json:"example"`
	CommonMisconception string `json:"commonMisconceptions"`
}

type UpdateFlashcardRequest struct {
	Title               string `json:"title" binding:"required"`
	Category            string `json:"category" binding:"required"`
	WhyItMatters        string `json:"whyItMatters" binding:"required"`
	Definition          string `json:"definition" binding:"required"`
	Example             string `json:"example"`
	CommonMisconception string `json:"commonMisconceptions"`
}

type Documents struct {
	Filename     string    `json:"filename"`
	Filetype     string    `json:"filetype"`
	StorageURL   string    `json:"storageUrl"`
	UploadStatus string    `json:"uploadStatus"`
	UploadAt     time.Time `json:"uploadAt"`
}

type LearningResource struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Summary     string    `json:"summary"`
	Source      string    `json:"source"`
	ImageURL    string    `json:"imageUrl"`
	Related     []string  `json:"related"`
	PublishedAt time.Time `json:"publishedAt"`
	URL         string    `json:"url"`
}

type FinnhubNewsItem struct {
	ID       int64  `json:"id"`
	Category string `json:"category"`
	DateTime int64  `json:"datetime"`
	Headline string `json:"headline"`
	Image    string `json:"image"`
	Related  string `json:"related"`
	Source   string `json:"source"`
	Summary  string `json:"summary"`
	URL      string `json:"url"`
}

// ========== RAG Pipeline Models ==========

type DocumentUploadRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
	Filename  string `json:"filename" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type DocumentUploadResponse struct {
	Success          bool   `json:"success"`
	DocumentID       string `json:"documentId"`
	ChunksCount      int    `json:"chunksCount"`
	EmbeddingsStored int    `json:"embeddingsStored"`
	EmbeddingModel   string `json:"embeddingModel"`
	Message          string `json:"message"`
}

type RAGQueryRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"topK"`
}

type RetrievedChunk struct {
	Text            string  `json:"text"`
	Filename        string  `json:"filename"`
	SimilarityScore float64 `json:"similarityScore"`
}

type RAGQueryResponse struct {
	Success           bool             `json:"success"`
	Query             string           `json:"query"`
	RetrievedChunks   []RetrievedChunk `json:"retrievedChunks"`
	GeneratedResponse string           `json:"generatedResponse"`
	ChunkCount        int              `json:"chunkCount"`
	Message           string           `json:"message"`
}
