package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"finedu-backend/internal/models"
)

// RAGHandler encapsulates RAG pipeline operations
type RAGHandler struct {
	pool          *pgxpool.Pool
	defaultUserID string
}

// NewRAGHandler creates a new RAG handler
func NewRAGHandler(pool *pgxpool.Pool, defaultUserID string) *RAGHandler {
	return &RAGHandler{pool: pool, defaultUserID: defaultUserID}
}

// UploadDocument handles PDF file upload and embedding
// POST /api/documents/upload
func (h *RAGHandler) UploadDocument(c *gin.Context) {
	// Extract multipart form data
	sessionID := c.PostForm("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Validate PDF file
	if filepath.Ext(file.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only PDF files are supported"})
		return
	}

	// Get user ID from context, or use default
	userID, exists := c.Get("userID")
	var userIDStr string
	if exists {
		userIDStr = userID.(string)
	} else {
		userIDStr = h.defaultUserID
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer src.Close()

	pdfData, err := ioutil.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file content"})
		return
	}

	// Generate document ID
	documentID := uuid.New().String()

	ctx := context.Background()

	// Insert document record
	err = h.pool.QueryRow(
		ctx,
		`INSERT INTO document (id, user_id, session_id, filename, filetype, storage_url, upload_status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		documentID, userIDStr, sessionID, file.Filename, "pdf", "", "processing",
	).Scan(&documentID)

	if err != nil {
		log.Printf("Error inserting document: %v", err)
		c.JSON(http.StatusInternalServerError, models.DocumentUploadResponse{
			Success: false,
			Message: fmt.Sprintf("Error storing document: %v", err),
		})
		return
	}

	// Call Python RAG pipeline to process PDF
	resp, err := callPythonRAGPipeline(
		"process_pdf_sync",
		map[string]interface{}{
			"user_id":     userIDStr,
			"document_id": documentID,
			"session_id":  sessionID,
			"filename":    file.Filename,
			"pdf_data":    base64.StdEncoding.EncodeToString(pdfData),
		},
	)

	if err != nil {
		log.Printf("Error calling RAG pipeline: %v", err)
		c.JSON(http.StatusInternalServerError, models.DocumentUploadResponse{
			Success: false,
			Message: fmt.Sprintf("Error processing PDF: %v", err),
		})
		return
	}

	// Parse response from Python
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Printf("Error parsing RAG response: %v", err)
		c.JSON(http.StatusInternalServerError, models.DocumentUploadResponse{
			Success: false,
			Message: "Error parsing pipeline response",
		})
		return
	}

	// Update document upload status to completed
	_, err = h.pool.Exec(
		ctx,
		`UPDATE document SET upload_status = $1 WHERE id = $2`,
		"completed", documentID,
	)
	if err != nil {
		log.Printf("Error updating document status: %v", err)
	}

	chunksCount := 0
	embeddingsStored := 0
	if chunks, ok := result["chunks_count"].(float64); ok {
		chunksCount = int(chunks)
	}
	if embeddings, ok := result["embeddings_stored"].(float64); ok {
		embeddingsStored = int(embeddings)
	}

	c.JSON(http.StatusOK, models.DocumentUploadResponse{
		Success:          true,
		DocumentID:       documentID,
		ChunksCount:      chunksCount,
		EmbeddingsStored: embeddingsStored,
		EmbeddingModel:   "models/embedding-001",
		Message:          "Document processed successfully",
	})
}

// QueryRAG handles RAG queries with similarity search and LLM response
// POST /api/rag/query
func (h *RAGHandler) QueryRAG(c *gin.Context) {
	var req models.RAGQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, models.RAGQueryResponse{
			Success: false,
			Message: "Query cannot be empty",
		})
		return
	}

	// Default top_k to 5 if not provided
	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}

	// Get user ID from context, or use default
	userID, exists := c.Get("userID")
	var userIDStr string
	if exists {
		userIDStr = userID.(string)
	} else {
		userIDStr = h.defaultUserID
	}

	// Call Python RAG pipeline to process query
	resp, err := callPythonRAGPipeline(
		"query_documents",
		map[string]interface{}{
			"user_id":    userIDStr,
			"query_text": req.Query,
			"top_k":      topK,
		},
	)

	if err != nil {
		log.Printf("Error calling RAG query: %v", err)
		c.JSON(http.StatusInternalServerError, models.RAGQueryResponse{
			Success: false,
			Message: fmt.Sprintf("Error processing query: %v", err),
		})
		return
	}

	// Parse response from Python
	var result models.RAGQueryResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		log.Printf("Error parsing query response: %v", err)
		c.JSON(http.StatusInternalServerError, models.RAGQueryResponse{
			Success: false,
			Message: "Error parsing query response",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// callPythonRAGPipeline calls a Python function from the RAG module
func callPythonRAGPipeline(funcName string, args map[string]interface{}) ([]byte, error) {
	// Prepare JSON arguments
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("error marshaling args: %w", err)
	}

	// Save arguments to temporary JSON file
	argsFile, err := ioutil.TempFile("", "rag_args_*.json")
	if err != nil {
		return nil, fmt.Errorf("error creating temp args file: %w", err)
	}
	defer os.Remove(argsFile.Name())

	if _, err := argsFile.Write(argsJSON); err != nil {
		argsFile.Close()
		return nil, fmt.Errorf("error writing args file: %w", err)
	}
	argsFile.Close()

	// Get the backend directory
	backendDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %w", err)
	}

	// RAG pipeline directory
	ragPipelineDir := filepath.Join(backendDir, "internal", "ragPipeline")

	// Create a temporary Python script that reads from the args file
	scriptContent := fmt.Sprintf(`
import sys
import json
# Add ragPipeline directory to path
sys.path.insert(0, r'%s')
from rag import %s
import os
os.environ.setdefault('DATABASE_URL', '%s')
os.environ.setdefault('GEMINI_API_KEY', '%s')

# Read arguments from file
with open(r'%s', 'r') as f:
    args = json.load(f)

result = %s(**args)
print(json.dumps(result))
`, ragPipelineDir, funcName, os.Getenv("DATABASE_URL"), os.Getenv("GEMINI_API_KEY"), argsFile.Name(), funcName)

	// Write script to temporary file
	scriptFile, err := ioutil.TempFile("", "rag_script_*.py")
	if err != nil {
		return nil, fmt.Errorf("error creating temp script file: %w", err)
	}
	defer os.Remove(scriptFile.Name())

	if _, err := scriptFile.WriteString(scriptContent); err != nil {
		scriptFile.Close()
		return nil, fmt.Errorf("error writing script file: %w", err)
	}
	scriptFile.Close()

	// Get the venv Python executable with absolute path
	venvPython := filepath.Join(backendDir, "venv", "Scripts", "python.exe")

	// Use system python if venv doesn't exist
	if _, err := os.Stat(venvPython); os.IsNotExist(err) {
		venvPython = "python"
	}

	// Call Python script from temp file with ragPipeline directory as working directory
	cmd := exec.Command(venvPython, scriptFile.Name())
	cmd.Dir = ragPipelineDir

	output, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("python error: %s", string(ee.Stderr))
		}
		return nil, fmt.Errorf("error executing python: %w", err)
	}

	return output, nil
}
