# DOCX Text Extraction Implementation Plan

## Overview
This document outlines the implementation plan for DOCX text extraction functionality in the ЖЦП parsing module. The DOCX parser will be responsible for extracting text content from DOCX documents while preserving document structure, formatting, and hierarchy.

## Requirements

### Functional Requirements
- Extract text content from DOCX files up to 50MB
- Preserve document structure (paragraphs, headings, lists)
- Extract tables and their content
- Handle various formatting (bold, italic, underlined text)
- Extract document metadata (author, creation date, etc.)
- Handle embedded objects and images (metadata only)

### Technical Requirements
- Python 3.8+ compatibility
- Cross-platform support (Windows, Linux, macOS)
- Memory efficient processing for large documents
- Error handling for corrupted or invalid DOCX files

## Implementation Approach

### Core DOCX Parser Class
```python
import docx
from docx.document import Document
from docx.table import Table
from docx.text.paragraph import Paragraph
from docx.oxml.table import CT_Tbl
from docx.oxml.text.paragraph import CT_P
from docx.table import _Cell
from typing import Optional, Dict, Any, List
import logging
from xml.etree import ElementTree as ET

class DOCXExtractor:
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def extract_text(self, docx_path: str) -> Dict[str, Any]:
        """
        Extract text and structure from DOCX file
        
        Args:
            docx_path: Path to the DOCX file
            
        Returns:
            Dictionary containing extracted text and metadata
        """
        try:
            doc = docx.Document(docx_path)
            
            result = {
                'text': '',
                'metadata': self._extract_metadata(doc),
                'structure': [],
                'tables': [],
                'paragraphs': [],
                'has_lists': False,
                'word_count': 0
            }
            
            # Extract paragraphs and preserve structure
            for element in doc.element.body:
                if isinstance(element, CT_P):
                    paragraph = Paragraph(element, doc)
                    paragraph_text = paragraph.text.strip()
                    if paragraph_text:
                        result['text'] += paragraph_text + '\n'
                        
                        # Analyze paragraph structure
                        para_info = self._analyze_paragraph(paragraph)
                        result['paragraphs'].append(para_info)
                        result['structure'].append(para_info)
                        
                elif isinstance(element, CT_Tbl):
                    table = Table(element, doc)
                    table_data = self._extract_table(table)
                    result['tables'].append(table_data)
            
            # Count words
            result['word_count'] = len(result['text'].split())
            
            return result
            
        except Exception as e:
            self.logger.error(f"DOCX extraction failed: {str(e)}")
            raise
    
    def _extract_metadata(self, doc: Document) -> Dict[str, Any]:
        """Extract document metadata"""
        core_props = doc.core_properties
        return {
            'title': core_props.title or '',
            'subject': core_props.subject or '',
            'author': core_props.author or '',
            'keywords': core_props.keywords or '',
            'comments': core_props.comments or '',
            'last_modified_by': core_props.last_modified_by or '',
            'revision': core_props.revision,
            'created': str(core_props.created) if core_props.created else '',
            'modified': str(core_props.modified) if core_props.modified else ''
        }
    
    def _analyze_paragraph(self, paragraph: Paragraph) -> Dict[str, Any]:
        """Analyze paragraph structure and formatting"""
        para_info = {
            'text': paragraph.text.strip(),
            'style': paragraph.style.name if paragraph.style else '',
            'level': 0,  # For headings
            'type': 'paragraph',
            'runs': []
        }
        
        # Determine paragraph type and level
        style_name = para_info['style'].lower()
        if 'heading' in style_name or 'заголовок' in style_name:
            para_info['type'] = 'heading'
            # Extract heading level from style name (e.g., "Heading 1")
            for char in style_name:
                if char.isdigit():
                    para_info['level'] = int(char)
                    break
        elif paragraph.text.strip().startswith(('*', '-', '•')):
            para_info['type'] = 'list_item'
        elif any(c.isdigit() for c in paragraph.text.strip()[:10]) and any(c in '.)' for c in paragraph.text.strip()[:10]):
            para_info['type'] = 'numbered_list_item'
        
        # Extract run information (text formatting)
        for run in paragraph.runs:
            run_info = {
                'text': run.text,
                'bold': run.bold or False,
                'italic': run.italic or False,
                'underline': run.underline or False,
                'font_name': run.font.name if run.font else '',
                'font_size': str(run.font.size) if run.font and run.font.size else ''
            }
            para_info['runs'].append(run_info)
        
        return para_info
    
    def _extract_table(self, table: Table) -> Dict[str, Any]:
        """Extract table content and structure"""
        table_data = {
            'rows': len(table.rows),
            'columns': len(table.columns) if table.columns else 0,
            'data': [],
            'headers': []
        }
        
        for i, row in enumerate(table.rows):
            row_data = []
            for cell in row.cells:
                cell_text = cell.text.strip()
                row_data.append(cell_text)
                
                # If this is the first row, consider it as headers
                if i == 0:
                    table_data['headers'].append(cell_text)
            
            table_data['data'].append(row_data)
        
        return table_data
    
    def extract_tables(self, docx_path: str) -> List[Dict[str, Any]]:
        """Extract all tables from DOCX"""
        doc = docx.Document(docx_path)
        tables = []
        
        for i, table in enumerate(doc.tables):
            table_info = {
                'index': i,
                'rows': len(table.rows),
                'columns': len(table.columns),
                'data': []
            }
            
            for row in table.rows:
                row_data = [cell.text.strip() for cell in row.cells]
                table_info['data'].append(row_data)
            
            tables.append(table_info)
        
        return tables
```

## Advanced Features

### Text Preprocessing for LLM Input
```python
import re

class DOCXTextPreprocessor:
    @staticmethod
    def clean_text(text: str) -> str:
        """Clean extracted DOCX text for LLM processing"""
        # Remove excessive whitespace
        text = re.sub(r'\s+', ' ', text)
        
        # Remove special characters that might be formatting artifacts
        text = re.sub(r'\x0b', ' ', text)  # Vertical tab
        text = re.sub(r'\x0c', ' ', text)  # Form feed
        
        # Normalize line breaks
        text = re.sub(r'\n\s*\n', '\n\n', text)
        
        # Remove extra spaces around punctuation
        text = re.sub(r'\s+([,.!?;:])', r'\1', text)
        
        return text.strip()
    
    @staticmethod
    def preserve_document_structure(paragraphs: List[Dict]) -> str:
        """Preserve document structure while creating text for LLM"""
        structured_text = []
        
        for para in paragraphs:
            text = para['text']
            para_type = para['type']
            
            if para_type == 'heading':
                # Add markers for headings
                level = para.get('level', 1)
                header_marker = '#' * level
                structured_text.append(f"{header_marker} {text}")
            elif para_type in ['list_item', 'numbered_list_item']:
                # Preserve list structure
                structured_text.append(f"• {text}")
            else:
                # Regular paragraph
                structured_text.append(text)
        
        return '\n\n'.join(structured_text)
```

## Enhanced Parser with Formatting Awareness

```python
class AdvancedDOCXExtractor:
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def extract_with_formatting(self, docx_path: str) -> Dict[str, Any]:
        """Extract content with detailed formatting information"""
        doc = docx.Document(docx_path)
        
        result = {
            'content': {
                'text': '',
                'formatted_elements': [],
                'structure': []
            },
            'metadata': self._extract_metadata(doc),
            'tables': self._extract_all_tables(doc),
            'lists': [],
            'images': []  # Metadata only
        }
        
        # Process document elements
        for i, element in enumerate(doc.element.body):
            if isinstance(element, CT_P):
                paragraph = Paragraph(element, doc)
                formatted_para = self._process_paragraph(paragraph, i)
                result['content']['formatted_elements'].append(formatted_para)
                result['content']['structure'].append({
                    'type': 'paragraph',
                    'element': formatted_para
                })
                result['content']['text'] += paragraph.text + '\n'
                
            elif isinstance(element, CT_Tbl):
                table = Table(element, doc)
                table_data = self._process_table(table, i)
                result['content']['structure'].append({
                    'type': 'table',
                    'element': table_data
                })
        
        return result
    
    def _process_paragraph(self, paragraph: Paragraph, index: int) -> Dict[str, Any]:
        """Process paragraph with detailed formatting"""
        para_info = {
            'index': index,
            'text': paragraph.text,
            'style': paragraph.style.name if paragraph.style else '',
            'properties': {
                'alignment': self._get_alignment(paragraph),
                'indentation': self._get_indentation(paragraph)
            },
            'runs': self._extract_runs(paragraph.runs),
            'type': self._classify_paragraph_type(paragraph)
        }
        
        return para_info
    
    def _get_alignment(self, paragraph: Paragraph) -> str:
        """Get paragraph alignment"""
        try:
            return paragraph.alignment.name if paragraph.alignment else 'LEFT'
        except:
            return 'LEFT'
    
    def _get_indentation(self, paragraph: Paragraph) -> Dict[str, Any]:
        """Get paragraph indentation properties"""
        pPr = paragraph._element.pPr
        if pPr is not None:
            ind = pPr.ind
            if ind is not None:
                return {
                    'left': ind.left,
                    'right': ind.right,
                    'first_line': ind.firstLine,
                    'hanging': ind.hanging
                }
        return {}
    
    def _extract_runs(self, runs) -> List[Dict[str, Any]]:
        """Extract detailed run information"""
        run_list = []
        for run in runs:
            run_info = {
                'text': run.text,
                'formatting': {
                    'bold': run.bold or False,
                    'italic': run.italic or False,
                    'underline': run.underline is not None,
                    'font': {
                        'name': run.font.name if run.font and run.font.name else '',
                        'size': str(run.font.size) if run.font and run.font.size else '',
                        'color': str(run.font.color.rgb) if run.font and run.font.color and run.font.color.rgb else ''
                    }
                }
            }
            run_list.append(run_info)
        return run_list
    
    def _classify_paragraph_type(self, paragraph: Paragraph) -> str:
        """Classify paragraph type based on content and style"""
        text = paragraph.text.strip().lower()
        style = paragraph.style.name.lower() if paragraph.style else ''
        
        if 'heading' in style or 'заголовок' in style:
            return 'heading'
        elif text.startswith(('*', '-', '•', '+')):
            return 'list_item'
        elif re.match(r'^\d+[\.\)]', text):
            return 'numbered_list_item'
        elif 'таблица' in text or 'табл.' in text:
            return 'table_reference'
        else:
            return 'paragraph'
    
    def _extract_all_tables(self, doc: Document) -> List[Dict[str, Any]]:
        """Extract all tables with detailed structure"""
        tables = []
        for i, table in enumerate(doc.tables):
            table_info = {
                'index': i,
                'rows': len(table.rows),
                'columns': len(table.columns),
                'header_row': self._extract_header_row(table),
                'data_rows': self._extract_data_rows(table),
                'cell_structure': self._analyze_cell_structure(table)
            }
            tables.append(table_info)
        return tables
    
    def _extract_header_row(self, table: Table) -> List[str]:
        """Extract header row if it exists"""
        if len(table.rows) > 0:
            first_row = table.rows[0]
            headers = [cell.text.strip() for cell in first_row.cells]
            # Simple heuristic: if most cells in first row have text, consider it header
            if sum(1 for h in headers if h) > len(headers) // 2:
                return headers
        return []
    
    def _extract_data_rows(self, table: Table) -> List[List[str]]:
        """Extract data rows (excluding header if identified)"""
        start_row = 1 if self._extract_header_row(table) else 0
        data_rows = []
        for row in table.rows[start_row:]:
            row_data = [cell.text.strip() for cell in row.cells]
            data_rows.append(row_data)
        return data_rows
    
    def _analyze_cell_structure(self, table: Table) -> Dict[str, Any]:
        """Analyze cell structure including merges"""
        structure = {
            'has_merged_cells': False,
            'merged_ranges': [],
            'cell_widths': []
        }
        
        # Analyze first row to get column widths
        if table.rows:
            first_row = table.rows[0]
            for cell in first_row.cells:
                # This is a simplified analysis
                pass
        
        return structure
```

## Error Handling and Validation

```python
class DOCXValidator:
    @staticmethod
    def validate_docx(docx_path: str) -> Dict[str, Any]:
        """Validate DOCX file before processing"""
        import os
        from zipfile import ZipFile
        
        validation_result = {
            'is_valid': False,
            'file_size': 0,
            'errors': [],
            'is_zip_file': False
        }
        
        try:
            # Check file exists
            if not os.path.exists(docx_path):
                validation_result['errors'].append('File does not exist')
                return validation_result
            
            # Check file size
            file_size = os.path.getsize(docx_path)
            validation_result['file_size'] = file_size
            
            if file_size > 50 * 1024 * 1024:  # 50MB
                validation_result['errors'].append('File size exceeds 50MB limit')
            
            # Check file extension
            if not docx_path.lower().endswith('.docx'):
                validation_result['errors'].append('File is not a DOCX')
            
            # Check if it's a valid zip file (DOCX is a zip archive)
            try:
                with ZipFile(docx_path, 'r') as zip_file:
                    # Check if it contains required DOCX files
                    required_files = ['word/document.xml', 'word/_rels/document.xml.rels']
                    docx_files = zip_file.namelist()
                    
                    if all(req_file in docx_files for req_file in required_files):
                        validation_result['is_valid'] = True
                        validation_result['is_zip_file'] = True
                    else:
                        validation_result['errors'].append('Not a valid DOCX structure')
            except:
                validation_result['errors'].append('File is not a valid zip archive')
                
        except Exception as e:
            validation_result['errors'].append(f'Validation error: {str(e)}')
        
        return validation_result
```

## Integration with Main Parser

```python
class ЖЦПParser:
    def __init__(self):
        self.pdf_extractor = PDFExtractor()
        self.docx_extractor = AdvancedDOCXExtractor()
        self.pdf_validator = PDFValidator()
        self.docx_validator = DOCXValidator()
        self.text_preprocessor = TextPreprocessor()
        self.docx_preprocessor = DOCXTextPreprocessor()
    
    def parse_docx(self, docx_path: str) -> Dict[str, Any]:
        """Main method to parse DOCX ЖЦП document"""
        
        # Validate DOCX
        validation = self.docx_validator.validate_docx(docx_path)
        if not validation['is_valid']:
            raise ValueError(f"Invalid DOCX: {', '.join(validation['errors'])}")
        
        # Extract content with formatting
        extraction_result = self.docx_extractor.extract_with_formatting(docx_path)
        
        # Create structured text preserving document hierarchy
        structured_text = self.docx_preprocessor.preserve_document_structure(
            extraction_result['content']['formatted_elements']
        )
        
        # Clean text for LLM processing
        cleaned_text = self.docx_preprocessor.clean_text(structured_text)
        
        # Add tables content to text if needed
        for table in extraction_result['tables']:
            table_text = self._table_to_text(table)
            cleaned_text += f"\n\nTABLE: {table_text}"
        
        # Return structured result
        return {
            'original_path': docx_path,
            'extracted_content': cleaned_text,
            'metadata': extraction_result['metadata'],
            'structure': extraction_result['content']['structure'],
            'tables': extraction_result['tables'],
            'word_count': sum(1 for word in cleaned_text.split() if word.strip())
        }
    
    def _table_to_text(self, table: Dict[str, Any]) -> str:
        """Convert table structure to text representation"""
        text_parts = []
        
        # Add headers if they exist
        if table['header_row']:
            text_parts.append(" | ".join(table['header_row']))
            text_parts.append("-" * 50)
        
        # Add data rows
        for row in table['data_rows']:
            text_parts.append(" | ".join(row))
        
        return "\n".join(text_parts)
```

## Testing Strategy

### Unit Tests
- Test with various DOCX formats and layouts
- Test with documents containing Cyrillic text
- Test error handling for invalid DOCXs
- Test performance with large documents
- Test table extraction accuracy

### Integration Tests
- Test end-to-end DOCX processing pipeline
- Test with real ЖЦП documents
- Test edge cases (empty tables, complex formatting, etc.)

## Dependencies
- `python-docx>=0.8.11`
- `lxml>=4.6.0` (dependency of python-docx)

## Performance Considerations
- For large documents, implement streaming processing
- Use memory-efficient processing for text extraction
- Implement caching for frequently accessed documents
- Consider parallel processing for multiple documents

## Special Considerations for Russian Documents
- Handle Cyrillic character encoding properly
- Recognize Russian document structure patterns
- Handle Russian date formats in text
- Preserve Russian typography and formatting

This implementation plan provides a comprehensive approach to DOCX text extraction that preserves document structure and formatting while ensuring reliable extraction for the ЖЦП parsing module.