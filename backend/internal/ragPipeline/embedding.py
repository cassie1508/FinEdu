"""
Embedding module using Google Gemini API.
Handles text-to-embedding conversion for RAG pipeline.
"""

import os
import logging
from typing import List, Dict, Optional
from google import genai

logger = logging.getLogger(__name__)

# Auto-detected embedding dimension (will be set on first embedding)
_EMBEDDING_DIMENSION: Optional[int] = None
_GEMINI_CLIENT = None  # Global Gemini client


def init_gemini(api_key: str = None) -> None:
    """Initialize Gemini API client."""
    global _GEMINI_CLIENT
    if api_key is None:
        api_key = os.getenv("GEMINI_API_KEY")
    if not api_key:
        raise ValueError("GEMINI_API_KEY environment variable is not set")
    _GEMINI_CLIENT = genai.Client(api_key=api_key)

    logger.info("Gemini API initialized")
    return _GEMINI_CLIENT


def get_embedding_dimension() -> int:
    """Get the embedding dimension (auto-detected from Gemini API)."""
    global _EMBEDDING_DIMENSION
    return _EMBEDDING_DIMENSION if _EMBEDDING_DIMENSION else 768


def _validate_embedding_dimension(embedding: List[float]) -> None:
    """Validate and auto-detect embedding dimension from Gemini response."""
    global _EMBEDDING_DIMENSION
    dim = len(embedding)
    
    if _EMBEDDING_DIMENSION is None:
        _EMBEDDING_DIMENSION = dim
        logger.info(f"Auto-detected embedding dimension from Gemini: {dim}")
    elif _EMBEDDING_DIMENSION != dim:
        raise ValueError(
            f"Embedding dimension mismatch. Expected {_EMBEDDING_DIMENSION}, "
            f"got {dim}. Database schema uses {_EMBEDDING_DIMENSION} dimensions."
        )


def embed_text(text: str, model: str = "models/embedding-001") -> List[float]:
    """
    Embed a single text using Gemini embedding model.
    
    Args:
        text: Text to embed
        model: Gemini embedding model name
    
    Returns:
        List of floats representing the embedding vector (auto-detected dimension)
    """
    global _GEMINI_CLIENT
    if _GEMINI_CLIENT is None:
        raise ValueError("Gemini client not initialized. Call init_gemini() first.")
    
    try:
        result = _GEMINI_CLIENT.models.embed_content(
            model=model,
            content=text,
        )
        embedding = result.embedding
        _validate_embedding_dimension(embedding)
        return embedding
    except Exception as e:
        logger.error(f"Error embedding text: {e}")
        raise


def embed_texts(texts: List[str], model: str = "models/embedding-001") -> List[List[float]]:
    """
    Embed multiple texts using Gemini embedding model (batch).
    
    Args:
        texts: List of texts to embed
        model: Gemini embedding model name
    
    Returns:
        List of embedding vectors
    """
    global _GEMINI_CLIENT
    if _GEMINI_CLIENT is None:
        raise ValueError("Gemini client not initialized. Call init_gemini() first.")
    
    try:
        result = _GEMINI_CLIENT.models.embed_content(
            model=model,
            content=texts,
        )
        embeddings = result.embeddings
        return embeddings
    except Exception as e:
        logger.error(f"Error embedding texts: {e}")
        raise


def embed_query(query: str, model: str = "models/embedding-001") -> List[float]:
    """
    Embed a query text using Gemini embedding model.
    
    Args:
        query: Query text to embed
        model: Gemini embedding model name
    
    Returns:
        List of floats representing the embedding vector (auto-detected dimension)
    """
    global _GEMINI_CLIENT
    if _GEMINI_CLIENT is None:
        raise ValueError("Gemini client not initialized. Call init_gemini() first.")
    
    try:
        result = _GEMINI_CLIENT.models.embed_content(
            model=model,
            content=query,
        )
        embedding = result.embedding
        _validate_embedding_dimension(embedding)
        return embedding
    except Exception as e:
        logger.error(f"Error embedding query: {e}")
        raise
