# Document Parsing Modules Plan for Go Implementation

## Overview

This document outlines the plan for implementing PDF and DOCX parsing functionality in Go, based on the existing Python implementation.

## PDF Parsing Module (internal/parsers/pdf/)

### Dependencies
- `github.com/pdfcpu/pdfcpu/pkg/api` - for PDF parsing and text extraction
- `github.com/unidoc/unipdf/v3` - alternative PDF library for complex layouts
- `golang.org/x/text` - for text processing and encoding

### Components

#### 1. PDFExtractor
```go
type PDFExtractor struct {
    logger *Logger
}

func (e *PDFExtractor) ExtractText(pdfPath string) (*PDFExtractionResult, error)

type PDFExtractionResult struct {
    Text      string
    Metadata  map[string]interface{}
    PageCount int
    HasTables bool
    Structure []StructureInfo
}

type StructureInfo struct {
    Page     int
    Type     string  // "text", "table", "image", etc.
    Content  string
    BBox     [4]float64  // bounding box coordinates
}
```

#### 2. PDFValidator
```go
type PDFValidator struct{}

func (v *PDFValidator) ValidatePDF(pdfPath string) (*ValidationResult, error)

type ValidationResult struct {
    IsValid  bool
    FileSize int64
    Errors   []string
}

func (v *PDFValidator) ValidatePDF(pdfPath string) (*ValidationResult, error)
```

#### 3. TextPreprocessor
```go
type TextPreprocessor struct{}

func (tp *TextPreprocessor) CleanText(text string) string
func (tp *TextPreprocessor) PreserveStructure(text string) string
```

### Implementation Strategy
1. Use pdfcpu for primary text extraction (handles most PDFs well)
2. Fallback to unidoc/unipdf for complex layouts if pdfcpu fails
3. Implement table detection and extraction
4. Preserve document structure information
5. Handle Cyrillic text properly

## DOCX Parsing Module (internal/parsers/docx/)

### Dependencies
- `github.com/baliance/sdcoffice` - for DOCX parsing
- `golang.org/x/text` - for text processing and encoding

### Components

#### 1. DOCXExtractor
```go
type DOCXExtractor struct {
    logger *Logger
}

func (e *DOCXExtractor) ExtractWithFormatting(docxPath string) (*DOCXExtractionResult, error)

type DOCXExtractionResult struct {
    Content struct {
        Text              string
        FormattedElements []FormattedElement
        Structure         []StructureElement
    }
    Metadata map[string]interface{}
    Tables   []TableInfo
    Lists    []ListInfo
    Images   []ImageInfo  // metadata only
}

type FormattedElement struct {
    Index      int
    Text       string
    Style      string
    Properties ElementProperties
    Runs       []RunInfo
    Type       string  // "heading", "paragraph", "list_item", etc.
}

type ElementProperties struct {
    Alignment  string
    Indentation IndentationInfo
}

type RunInfo struct {
    Text       string
    Formatting RunFormatting
}

type RunFormatting struct {
    Bold     bool
    Italic   bool
    Underline bool
    Font     FontInfo
}

type FontInfo struct {
    Name  string
    Size  float64
    Color string
}

type TableInfo struct {
    Index         int
    Rows          int
    Columns       int
    HeaderRow     []string
    DataRows      [][]string
    CellStructure CellStructureInfo
}

type StructureElement struct {
    Type    string
    Element interface{}
}
```

#### 2. DOCXValidator
```go
type DOCXValidator struct{}

func (v *DOCXValidator) ValidateDOCX(docxPath string) (*ValidationResult, error)
```

#### 3. DOCXTextPreprocessor
```go
type DOCXTextPreprocessor struct{}

func (tp *DOCXTextPreprocessor) CleanText(text string) string
func (tp *DOCXTextPreprocessor) PreserveDocumentStructure(elements []FormattedElement) string
```

### Implementation Strategy
1. Use baliance/sdcoffice for comprehensive DOCX parsing
2. Extract formatting information (styles, fonts, alignment)
3. Preserve document structure (headings, lists, tables)
4. Handle embedded objects and images (metadata only)
5. Extract tables with proper structure

## Text Preprocessing Module (internal/parsers/text_preprocessor.go)

### Components
```go
type TextPreprocessor struct{}

// CleanText removes artifacts and normalizes text
func (tp *TextPreprocessor) CleanText(text string) string

// PreserveStructure maintains document hierarchy for LLM processing
func (tp *TextPreprocessor) PreserveStructure(text string) string

// NormalizeCyrillic handles common Cyrillic character variations
func (tp *TextPreprocessor) NormalizeCyrillic(text string) string
```

## Common Interfaces

### Parser Interface
```go
type DocumentParser interface {
    Extract(path string) (*ExtractionResult, error)
    Validate(path string) (*ValidationResult, error)
}

type ExtractionResult struct {
    Text      string
    Metadata  map[string]interface{}
    Structure []interface{}
    Tables    []interface{}
    HasTables bool
    PageCount int
}
```

## Error Handling

### PDF/DOCX Specific Errors
```go
type DocumentError struct {
    Category ErrorCategory
    Message  string
    Details  map[string]interface{}
    Path     string
}

func NewDocumentError(category ErrorCategory, message string, path string) *DocumentError
```

## Performance Considerations

1. **Memory Management**: 
   - Process large documents in chunks
   - Use streaming when possible
   - Implement proper cleanup of resources

2. **Processing Speed**:
   - Cache frequently used patterns
   - Optimize regex operations
   - Use efficient string operations

3. **Large Document Handling**:
   - Implement progress tracking
   - Add cancellation support via context
   - Provide status updates during processing

## Testing Strategy

1. **Unit Tests**:
   - Test text extraction with various PDF/DOCX formats
   - Test validation logic
   - Test text preprocessing functions

2. **Integration Tests**:
   - End-to-end parsing of sample documents
   - Test error handling scenarios
   - Performance tests with large documents

3. **Sample Documents**:
   - Various PDF layouts (scanned, text-based, mixed)
   - DOCX documents with different structures
   - Edge cases (corrupted files, password-protected, etc.)

## Fallback Mechanisms

1. **PDF Parsing Fallbacks**:
   - Primary: pdfcpu
   - Secondary: unidoc/unipdf
   - Tertiary: command-line tools (if available)

2. **DOCX Parsing Fallbacks**:
   - Primary: baliance/sdcoffice
   - Secondary: xml package (manual parsing)
   - Tertiary: external converters

## Security Considerations

1. **File Validation**:
   - Verify file headers match expected format
   - Check file size limits
   - Validate file structure before processing

2. **Injection Prevention**:
   - Sanitize extracted text before LLM processing
   - Validate extracted data formats

3. **Resource Limits**:
   - Set memory limits for document processing
   - Implement timeouts for parsing operations