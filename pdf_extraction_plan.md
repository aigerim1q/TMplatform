# PDF Text Extraction Implementation Plan

## Overview
This document outlines the implementation plan for PDF text extraction functionality in the ЖЦП parsing module. The PDF parser will be responsible for extracting text content from PDF documents while preserving document structure and hierarchy.

## Requirements

### Functional Requirements
- Extract text content from PDF files up to 100MB
- Handle various PDF layouts (text blocks, tables, headers, footers)
- Preserve document hierarchy and structure
- Handle different fonts and encodings (including Cyrillic)
- Extract page numbers and document metadata
- Handle password-protected PDFs (optional)

### Technical Requirements
- Python 3.8+ compatibility
- Cross-platform support (Windows, Linux, macOS)
- Memory efficient processing for large documents
- Error handling for corrupted or invalid PDFs

## Implementation Options

### Option 1: PyPDF2 (Legacy but stable)
```python
import PyPDF2

def extract_text_pypdf2(pdf_path):
    with open(pdf_path, 'rb') as file:
        reader = PyPDF2.PdfReader(file)
        text = ""
        for page in reader.pages:
            text += page.extract_text()
    return text
```

### Option 2: pdfplumber (Recommended)
```python
import pdfplumber

def extract_text_pdfplumber(pdf_path):
    text = ""
    with pdfplumber.open(pdf_path) as pdf:
        for page in pdf.pages:
            text += page.extract_text() or ""
    return text
```

### Option 3: pymupdf (fitz) - Most Advanced
```python
import fitz  # pymupdf

def extract_text_pymupdf(pdf_path):
    doc = fitz.open(pdf_path)
    text = ""
    for page in doc:
        text += page.get_text()
    doc.close()
    return text
```

## Recommended Implementation: Hybrid Approach

### Core PDF Parser Class
```python
import pdfplumber
import fitz  # PyMuPDF
from typing import Optional, Dict, Any
import logging

class PDFExtractor:
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def extract_text(self, pdf_path: str) -> Dict[str, Any]:
        """
        Extract text from PDF with fallback mechanisms
        
        Args:
            pdf_path: Path to the PDF file
            
        Returns:
            Dictionary containing extracted text and metadata
        """
        result = {
            'text': '',
            'metadata': {},
            'page_count': 0,
            'has_tables': False,
            'structure': []
        }
        
        try:
            # Try PyMuPDF first (best for complex layouts)
            result = self._extract_with_pymupdf(pdf_path)
            
            # If PyMuPDF fails or returns empty text, try pdfplumber
            if not result['text'].strip():
                result = self._extract_with_pdfplumber(pdf_path)
                
        except Exception as e:
            self.logger.error(f"PDF extraction failed: {str(e)}")
            raise
        
        return result
    
    def _extract_with_pymupdf(self, pdf_path: str) -> Dict[str, Any]:
        """Extract text using PyMuPDF"""
        doc = fitz.open(pdf_path)
        
        result = {
            'text': '',
            'metadata': doc.metadata,
            'page_count': doc.page_count,
            'has_tables': False,
            'structure': []
        }
        
        for page_num in range(doc.page_count):
            page = doc.load_page(page_num)
            text = page.get_text()
            result['text'] += text + '\n'
            
            # Extract structure information
            blocks = page.get_text("dict")["blocks"]
            for block in blocks:
                if "lines" in block:
                    structure_info = {
                        'page': page_num + 1,
                        'type': 'text',
                        'content': '',
                        'bbox': block['bbox']
                    }
                    
                    # Combine lines in the block
                    for line in block.get("lines", []):
                        for span in line.get("spans", []):
                            structure_info['content'] += span["text"]
                    
                    result['structure'].append(structure_info)
        
        doc.close()
        return result
    
    def _extract_with_pdfplumber(self, pdf_path: str) -> Dict[str, Any]:
        """Extract text using pdfplumber"""
        with pdfplumber.open(pdf_path) as pdf:
            result = {
                'text': '',
                'metadata': {},
                'page_count': len(pdf.pages),
                'has_tables': False,
                'structure': []
            }
            
            for page_num, page in enumerate(pdf.pages):
                # Extract text
                page_text = page.extract_text() or ""
                result['text'] += page_text + '\n'
                
                # Check for tables
                tables = page.extract_tables()
                if tables:
                    result['has_tables'] = True
                
                # Extract structure
                chars = page.chars if hasattr(page, 'chars') else []
                structure_info = {
                    'page': page_num + 1,
                    'type': 'text',
                    'content': page_text[:100] + '...' if len(page_text) > 100 else page_text,
                    'char_count': len(chars)
                }
                result['structure'].append(structure_info)
        
        return result
    
    def extract_tables(self, pdf_path: str) -> list:
        """Extract tables from PDF"""
        tables = []
        with pdfplumber.open(pdf_path) as pdf:
            for page in pdf.pages:
                page_tables = page.extract_tables()
                for table in page_tables:
                    tables.append({
                        'page': page.page_number,
                        'data': table
                    })
        return tables
```

## Advanced Features

### Text Preprocessing for LLM Input
```python
import re

class TextPreprocessor:
    @staticmethod
    def clean_text(text: str) -> str:
        """Clean extracted text for LLM processing"""
        # Remove excessive whitespace
        text = re.sub(r'\s+', ' ', text)
        
        # Handle common PDF artifacts
        text = text.replace('\n', ' ').strip()
        
        # Normalize Cyrillic characters
        text = text.replace('ё', 'е')
        
        # Remove page numbers if they appear as standalone numbers
        text = re.sub(r'\b\d{1,3}\b\s*$', '', text, flags=re.MULTILINE)
        
        return text
    
    @staticmethod
    def preserve_structure(text: str) -> str:
        """Preserve document structure while cleaning"""
        # Maintain paragraph breaks
        text = re.sub(r'\n\s*\n', '\n\n', text)
        
        # Identify headers and sections
        lines = text.split('\n')
        processed_lines = []
        
        for line in lines:
            stripped = line.strip()
            if stripped and all(c.isupper() or not c.isalnum() for c in stripped[:50]):
                # Likely a header
                processed_lines.append(f"[HEADER] {stripped} [HEADER]")
            else:
                processed_lines.append(stripped)
        
        return '\n'.join(processed_lines)
```

## Error Handling and Validation

```python
class PDFValidator:
    @staticmethod
    def validate_pdf(pdf_path: str) -> Dict[str, Any]:
        """Validate PDF file before processing"""
        import os
        
        validation_result = {
            'is_valid': False,
            'file_size': 0,
            'errors': []
        }
        
        try:
            # Check file exists
            if not os.path.exists(pdf_path):
                validation_result['errors'].append('File does not exist')
                return validation_result
            
            # Check file size
            file_size = os.path.getsize(pdf_path)
            validation_result['file_size'] = file_size
            
            if file_size > 100 * 1024 * 1024:  # 100MB
                validation_result['errors'].append('File size exceeds 100MB limit')
            
            # Check file extension
            if not pdf_path.lower().endswith('.pdf'):
                validation_result['errors'].append('File is not a PDF')
            
            # Attempt to open with PyMuPDF to check if it's a valid PDF
            import fitz
            try:
                doc = fitz.open(pdf_path)
                doc.close()
                validation_result['is_valid'] = True
            except:
                validation_result['errors'].append('File is not a valid PDF')
                
        except Exception as e:
            validation_result['errors'].append(f'Validation error: {str(e)}')
        
        return validation_result
```

## Integration with Main Parser

```python
class ЖЦПParser:
    def __init__(self):
        self.pdf_extractor = PDFExtractor()
        self.text_preprocessor = TextPreprocessor()
        self.validator = PDFValidator()
    
    def parse_pdf(self, pdf_path: str) -> Dict[str, Any]:
        """Main method to parse PDF ЖЦП document"""
        
        # Validate PDF
        validation = self.validator.validate_pdf(pdf_path)
        if not validation['is_valid']:
            raise ValueError(f"Invalid PDF: {', '.join(validation['errors'])}")
        
        # Extract text
        extraction_result = self.pdf_extractor.extract_text(pdf_path)
        
        # Preprocess text
        cleaned_text = self.text_preprocessor.clean_text(extraction_result['text'])
        
        # Return structured result
        return {
            'original_path': pdf_path,
            'extracted_content': cleaned_text,
            'metadata': extraction_result['metadata'],
            'structure': extraction_result['structure'],
            'page_count': extraction_result['page_count'],
            'has_tables': extraction_result['has_tables']
        }
```

## Testing Strategy

### Unit Tests
- Test with various PDF formats and layouts
- Test with documents containing Cyrillic text
- Test error handling for invalid PDFs
- Test performance with large documents

### Integration Tests
- Test end-to-end PDF processing pipeline
- Test with real ЖЦП documents
- Test edge cases (empty pages, scanned documents, etc.)

## Dependencies
- `pdfplumber>=0.7.0`
- `PyMuPDF>=1.23.0` (fitz)
- `pymupdf>=1.23.0`

## Performance Considerations
- For large documents, implement chunked processing
- Use memory-efficient processing for text extraction
- Implement caching for frequently accessed documents
- Consider parallel processing for multiple documents

This implementation plan provides a robust, multi-engine approach to PDF text extraction that handles various document layouts and ensures reliable extraction for the ЖЦП parsing module.