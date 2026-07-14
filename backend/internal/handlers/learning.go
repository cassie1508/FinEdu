package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChatWithAI handles RAG-based finance chatbot conversations.
// Owner: Hiếu (AI Learning Center)
func ChatWithAI(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "ChatWithAI not yet implemented"})
}

// GetFlashcards returns flashcards for a given financial topic.
// Owner: Hiếu (AI Learning Center)
func GetFlashcards(c *gin.Context) {
	topic := c.Query("topic")
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "GetFlashcards not yet implemented",
		"topic":   topic,
	})
}
