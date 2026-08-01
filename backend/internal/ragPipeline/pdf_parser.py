"""
PDF Parser module for RAG pipeline.
Extracts text from PDF data and prepares for embedding.
"""

import logging
from typing import List
import io

try:
    import PyPDF2
except ImportError:
    PyPDF2 = None

logger = logging.getLogger(__name__)


def extract_text_from_pdf(pdf_data: bytes) -> str:
    """
    Extract text from PDF binary data.
    
    Args:
        pdf_data: PDF file content as bytes
    
    Returns:
        Extracted text from PDF
    """
    if PyPDF2 is None:
        raise ImportError("PyPDF2 is not installed. Install with: pip install PyPDF2")
    
    try:
        pdf_file = io.BytesIO(pdf_data)
        pdf_reader = PyPDF2.PdfReader(pdf_file)
        
        text = ""
        for page_num in range(len(pdf_reader.pages)):
            page = pdf_reader.pages[page_num]
            text += page.extract_text()
            text += "\n\n"
        
        logger.info(f"Extracted {len(text)} characters from PDF")
        return text.strip()
    
    except Exception as e:
        logger.error(f"Error extracting text from PDF: {e}")
        raise


def validate_pdf(pdf_data: bytes) -> bool:
    """
    Validate if the data is a valid PDF.
    
    Args:
        pdf_data: PDF file content as bytes
    
    Returns:
        True if valid PDF, False otherwise
    """
    try:
        if PyPDF2 is None:
            logger.warning("PyPDF2 not available for validation, skipping")
            return True
        
        pdf_file = io.BytesIO(pdf_data)
        pdf_reader = PyPDF2.PdfReader(pdf_file)
        return len(pdf_reader.pages) > 0
    except Exception as e:
        logger.error(f"PDF validation failed: {e}")
        return False
