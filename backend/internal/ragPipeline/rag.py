"""
RAG (Retrieval-Augmented Generation) Pipeline orchestration.
Handles document processing, chunking, embedding, and LLM-powered Q&A.
"""

import os
import logging
from typing import List, Dict
import uuid
import json
from datetime import datetime
from google import genai
import psycopg2
import base64

from embedding import init_gemini, embed_text, embed_texts, embed_query
from vector_store import VectorStore
from pdf_parser import extract_text_from_pdf, validate_pdf

logger = logging.getLogger(__name__)


class RAGPipeline:
    """Orchestrates the RAG pipeline for document Q&A."""
    
    def __init__(self, database_url: str = None, gemini_api_key: str = None):
        """
        Initialize RAG pipeline.
        
        Args:
            database_url: PostgreSQL connection string
            gemini_api_key: Gemini API key for embeddings and LLM
        """
        self.vector_store = VectorStore(database_url)
        init_gemini(gemini_api_key)
        self.embedding_model = "models/embedding-001"
        self.llm_model = os.getenv("LLM_MODEL", "gemini-2.5-flash")
        self.conn = psycopg2.connect(database_url or os.getenv("DATABASE_URL"))
        logger.info("RAG Pipeline initialized")
    
    def close(self):
        """Close all connections."""
        if self.vector_store:
            self.vector_store.close()
        if self.conn:
            self.conn.close()
    
    def chunk_text(self, text: str, chunk_size: int = 500, overlap: int = 50) -> List[str]:
        """
        Split text into overlapping chunks.
        
        Args:
            text: Text to chunk
            chunk_size: Size of each chunk (in characters)
            overlap: Character overlap between chunks
        
        Returns:
            List of text chunks
        """
        chunks = []
        start = 0
        while start < len(text):
            end = start + chunk_size
            chunks.append(text[start:end])
            start = end - overlap
        return chunks if chunks else [text]
    
    def process_text(
        self,
        user_id: str,
        document_id: str,
        session_id: str,
        filename: str,
        content: str,
    ) -> Dict:
        """
        Process text content: chunk, embed, and store in vector database.
        
        Args:
            user_id: UUID of the user
            document_id: UUID of the document
            session_id: Session ID for grouping
            filename: Original filename
            content: Text content to process
        
        Returns:
            Dictionary with processing results
        """
        try:
            logger.info(f"Processing document {filename} for user {user_id}")
            
            # Chunk the text
            chunks = self.chunk_text(content)
            logger.info(f"Created {len(chunks)} chunks from content")
            
            # Embed and store each chunk
            stored_count = 0
            for i, chunk in enumerate(chunks):
                try:
                    # Embed the chunk
                    embedding = embed_text(chunk, self.embedding_model)
                    
                    # Store in vector database
                    chunk_id = str(uuid.uuid4())
                    self.vector_store.add_chunk(
                        user_id=user_id,
                        document_id=document_id,
                        chunk_id=chunk_id,
                        filename=filename,
                        chunk_text=chunk,
                        embedding=embedding,
                        session_id=session_id,
                    )
                    stored_count += 1
                except Exception as e:
                    logger.warning(f"Failed to embed chunk {i}: {e}")
            
            logger.info(f"Successfully stored {stored_count} embeddings")
            return {
                "success": True,
                "document_id": document_id,
                "filename": filename,
                "chunks_count": len(chunks),
                "embeddings_stored": stored_count,
            }
        except Exception as e:
            logger.error(f"Error processing text: {e}")
            raise
    

    def process_pdf(
        self,
        user_id: str,
        document_id: str,
        session_id: str,
        filename: str,
        pdf_data: str,
    ) -> Dict:
        """
        Process a PDF file: extract text, chunk, embed, and store in vector database.
        
        Args:
            user_id: UUID of the user
            document_id: UUID of the document
            session_id: Session ID for grouping
            filename: Original filename
            pdf_data: Base64-encoded PDF data string
        
        Returns:
            Dictionary with processing results (chunks_count, embeddings_stored, etc.)
        """
        try:
            logger.info(f"Processing PDF {filename} for user {user_id}")
            
            # Decode PDF data
            pdf_bytes = base64.b64decode(pdf_data)
            
            # Validate PDF
            if not validate_pdf(pdf_bytes):
                raise ValueError("Invalid PDF file")
            
            # Extract text from PDF
            content = extract_text_from_pdf(pdf_bytes)
            logger.info(f"Extracted {len(content)} characters from PDF")
            
            # Reuse process_text to handle chunking, embedding, and storage
            return self.process_text(user_id, document_id, session_id, filename, content)
        
        except Exception as e:
            logger.error(f"Error processing PDF: {e}")
            raise
    
    def query(self, user_id: str, query_text: str, top_k: int = 5) -> Dict:
        """
        Query documents with RAG: retrieve similar chunks and generate response.
        
        Args:
            user_id: UUID of the user
            query_text: Question/query text
            top_k: Number of top similar chunks to retrieve
        
        Returns:
            Dictionary with retrieved_chunks, generated_response, and metadata
        """
        try:
            logger.info(f"Processing query for user {user_id}: {query_text[:50]}...")
            
            # Embed the query
            query_embedding = embed_query(query_text, self.embedding_model)
            logger.info("Query embedding generated")
            
            # Search for similar chunks
            similar_chunks = self.vector_store.similarity_search(
                query_embedding, user_id, top_k
            )
            logger.info(f"Retrieved {len(similar_chunks)} similar chunks")
            
            if not similar_chunks:
                return {
                    "success": True,
                    "query": query_text,
                    "retrieved_chunks": [],
                    "generated_response": "No relevant documents found.",
                    "chunk_count": 0,
                }
            
            # Prepare context from retrieved chunks
            context = "\n\n---\n\n".join([
                f"[{c['filename']} - Chunk {i}]\n{c['chunk_text']}"
                for i, c in enumerate(similar_chunks)
            ])
            
            # Generate response using Gemini LLM
            prompt = f"""You are a financial expert AI assistant. 
            Based on the following documents, answer the user's question.

            Documents:
            {context}

            User Question: {query_text}

            Provide a clear, concise answer based on the documents. 
            Rules:
            - Use the context as the primary source.
            - Do not invent facts that are not supported by the context.
            - If the context does not contain enough information, say so.
            - Explain financial concepts clearly and simply."""
            
            response = genai.models.generate_content(model=self.llm_model, prompt=prompt)
            generated_response = response.text
            logger.info("Generated LLM response")
            
            return {
                "success": True,
                "query": query_text,
                "retrieved_chunks": [
                    {
                        "text": c["chunk_text"],
                        "filename": c["filename"],
                        "similarity_score": round(c["similarity_score"], 4),
                    }
                    for c in similar_chunks
                ],
                "generated_response": generated_response,
                "chunk_count": len(similar_chunks),
            }
        except Exception as e:
            logger.error(f"Error processing query: {e}")
            raise


# Convenience functions for use from Go backend

_pipeline = None


def init_pipeline(database_url: str = None, gemini_api_key: str = None):
    """Initialize global RAG pipeline."""
    global _pipeline
    if database_url is None:
        database_url = os.getenv("DATABASE_URL")
    if gemini_api_key is None:
        gemini_api_key = os.getenv("GEMINI_API_KEY")
    _pipeline = RAGPipeline(database_url, gemini_api_key)


def process_document_sync(
    user_id: str,
    document_id: str,
    session_id: str,
    filename: str,
    content: str,
) -> Dict:
    """Process a document text (synchronous wrapper)."""
    if not _pipeline:
        init_pipeline()
    return _pipeline.process_text(user_id, document_id, session_id, filename, content)


def query_documents(user_id: str, query_text: str, top_k: int = 5) -> Dict:
    """Query documents (synchronous wrapper)."""
    if not _pipeline:
        init_pipeline()
    return _pipeline.query(user_id, query_text, top_k)


def process_pdf_sync(
    user_id: str,
    document_id: str,
    session_id: str,
    filename: str,
    pdf_data: str,
) -> Dict:
    """Process a PDF file (synchronous wrapper)."""
    if not _pipeline:
        init_pipeline()
    return _pipeline.process_pdf(user_id, document_id, session_id, filename, pdf_data)
