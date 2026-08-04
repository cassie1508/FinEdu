"""
Vector store module for pgvector operations.
Handles embedding storage and similarity search in PostgreSQL.
"""

import os
import logging
from typing import List, Dict, Tuple
import psycopg2
from psycopg2.extras import execute_values
from psycopg2.pool import SimpleConnectionPool
from embedding import get_embedding_dimension

logger = logging.getLogger(__name__)


class VectorStore:
    """PostgreSQL pgvector store for embeddings."""
    
    def __init__(self, database_url: str = None):
        """
        Initialize vector store connection.
        
        Args:
            database_url: PostgreSQL connection string (defaults to DATABASE_URL env var)
        """
        if database_url is None:
            database_url = os.getenv("DATABASE_URL")
        if not database_url:
            raise ValueError("DATABASE_URL environment variable is not set")
        
        self.database_url = database_url
        self.conn = None
        self._connect()
    
    def _connect(self):
        """Establish database connection."""
        try:
            self.conn = psycopg2.connect(self.database_url)
            logger.info("Connected to PostgreSQL")
        except Exception as e:
            logger.error(f"Failed to connect to PostgreSQL: {e}")
            raise
    
    def close(self):
        """Close database connection."""
        if self.conn:
            self.conn.close()
            logger.info("Closed PostgreSQL connection")
    
    def store_embedding(
        self,
        user_id: str,
        document_id: str,
        chunk_id: str,
        embedding: List[float],
        embedding_model: str = "models/embedding-001",
        metadata: Dict = None,
    ) -> str:
        """
        Store embedding in vector_store table.
        
        Args:
            user_id: UUID of the user
            document_id: UUID of the document
            chunk_id: UUID of the chunk
            embedding: Embedding vector (auto-detected dimension)
            embedding_model: Name of the embedding model
            metadata: Optional metadata dictionary
        
        Returns:
            Vector store ID
        """
        try:
            with self.conn.cursor() as cur:
                embedding_str = "[" + ",".join(str(x) for x in embedding) + "]"
                metadata = metadata or {}
                
                cur.execute(
                    """
                    INSERT INTO vector_store (user_id, document_id, chunk_id, embedding, embedding_model, metadata)
                    VALUES (%s, %s, %s, %s::vector, %s, %s)
                    ON CONFLICT (chunk_id) DO UPDATE
                    SET embedding = %s::vector, embedding_model = %s, metadata = %s
                    RETURNING id
                    """,
                    (
                        user_id,
                        document_id,
                        chunk_id,
                        embedding_str,
                        embedding_model,
                        str(metadata),
                        embedding_str,
                        embedding_model,
                        str(metadata),
                    ),
                )
                vector_id = cur.fetchone()[0]
                self.conn.commit()
                logger.info(f"Stored embedding for chunk {chunk_id}")
                return vector_id
        except Exception as e:
            self.conn.rollback()
            logger.error(f"Error storing embedding: {e}")
            raise
    
    def similarity_search(
        self,
        query_embedding: List[float],
        user_id: str,
        top_k: int = 5,
    ) -> List[Dict]:
        """
        Find top-k most similar embeddings using cosine similarity.
        
        Args:
            query_embedding: Auto-detected dimension query embedding vector
            user_id: UUID of the user (for filtering)
            top_k: Number of top results to return
        
        Returns:
            List of dictionaries with chunk_id, chunk_text, similarity_score, document_id, metadata
        """
        try:
            embedding_str = "[" + ",".join(str(x) for x in query_embedding) + "]"
            
            with self.conn.cursor() as cur:
                cur.execute(
                    """
                    SELECT
                        vs.chunk_id,
                        dc.chunk_text,
                        1 - (vs.embedding <=> %s::vector) as similarity_score,
                        vs.document_id,
                        vs.metadata,
                        d.filename
                    FROM vector_store vs
                    JOIN document_chunk dc ON vs.chunk_id = dc.id
                    JOIN document d ON vs.document_id = d.id
                    WHERE vs.user_id = %s
                    ORDER BY vs.embedding <=> %s::vector ASC
                    LIMIT %s
                    """,
                    (embedding_str, user_id, embedding_str, top_k),
                )
                results = cur.fetchall()
                
                return [
                    {
                        "chunk_id": r[0],
                        "chunk_text": r[1],
                        "similarity_score": float(r[2]),
                        "document_id": r[3],
                        "metadata": r[4],
                        "filename": r[5],
                    }
                    for r in results
                ]
        except Exception as e:
            logger.error(f"Error searching embeddings: {e}")
            raise
    