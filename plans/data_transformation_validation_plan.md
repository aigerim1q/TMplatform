# Data Transformation and Validation Modules Plan for Go Implementation

## Overview

This document outlines the plan for implementing data transformation and validation functionality in Go, based on the existing Python implementation.

## Data Transformation Module (internal/transformers/json_transformer.go)

### Dependencies
- `encoding/json` - for JSON parsing and marshaling
- `github.com/go-playground/validator/v10` - for data validation
- `golang.org/x/text` - for text processing
- `github.com/shopspring/decimal` - for precise decimal calculations (if needed)

### Core Types

#### Transformation Status
```go
type TransformationStatus string

const (
    TransformationStatusSuccess       TransformationStatus = "success"
    TransformationStatusPartial       TransformationStatus = "partial"
    TransformationStatusFailed        TransformationStatus = "failed"
    TransformationStatusValidationError TransformationStatus = "validation_error"
)
```

#### Transformation Result
```go
type TransformationResult struct {
    TransformedData    *ProjectStructure
    Status             TransformationStatus
    ConfidenceScore    float64
    ValidationErrors   []string
    ProcessingNotes    []string
    TokensUsed         TokenUsage
}
```

#### Project Structure Models
```go
type ProjectStructure struct {
    Project  Project  `json:"project" validate:"required"`
    Metadata Metadata `json:"metadata,omitempty"`
}

type Project struct {
    Title       string  `json:"title" validate:"required"`
    Description string  `json:"description" validate:"required"`
    Phases      []Phase `json:"phases" validate:"required,min=1"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type Phase struct {
    ID          string `json:"id" validate:"required"`
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
    StartDate   string `json:"start_date,omitempty" validate:"omitempty,date_format"`
    EndDate     string `json:"end_date,omitempty" validate:"omitempty,date_format,date_after_start"`
    Tasks       []Task `json:"tasks"`
}

type Task struct {
    ID                string               `json:"id" validate:"required"`
    Name              string               `json:"name" validate:"required"`
    Description       string               `json:"description"`
    StartDate         string               `json:"start_date,omitempty" validate:"omitempty,date_format"`
    EndDate           string               `json:"end_date,omitempty" validate:"omitempty,date_format,date_after_start"`
    ResponsiblePersons []ResponsiblePerson `json:"responsible_persons"`
    Dependencies      []string             `json:"dependencies"`
    Status            string               `json:"status" validate:"oneof=planned in_progress completed"`
}

type ResponsiblePerson struct {
    Name    string `json:"name" validate:"required"`
    Role    string `json:"role"`
    Contact string `json:"contact"`
}

type Metadata struct {
    SourceDocument   string            `json:"source_document"`
    ExtractionDate   string            `json:"extraction_date"`
    ConfidenceScore  float64           `json:"confidence_score"`
    ProcessingTime   float64           `json:"processing_time"`
    ValidationResults ValidationResults `json:"validation_results,omitempty"`
}
```

### DataTransformer Implementation
```go
type DataTransformer struct {
    logger        *Logger
    datePatterns  []string
    validator     *validator.Validate
}

func NewDataTransformer() *DataTransformer
func (dt *DataTransformer) Transform(llmResponse string) *TransformationResult
func (dt *DataTransformer) normalizeData(rawData map[string]interface{}) *ProjectStructure
func (dt *DataTransformer) normalizePhases(rawPhases []interface{}) []Phase
func (dt *DataTransformer) normalizeTasks(rawTasks []interface{}, phaseID string) []Task
func (dt *DataTransformer) normalizeResponsibles(rawResponsibles []interface{}) []ResponsiblePerson
func (dt *DataTransformer) normalizeDate(dateValue interface{}) *string
func (dt *DataTransformer) normalizeText(text interface{}) string
func (dt *DataTransformer) normalizeStatus(status interface{}) string
func (dt *DataTransformer) validateData(data *ProjectStructure) *ValidationResult
func (dt *DataTransformer) businessValidation(data *ProjectStructure) *ValidationResult
func (dt *DataTransformer) createPartialTransformation(data *ProjectStructure, validation *ValidationResult) *ProjectStructure
func (dt *DataTransformer) calculateConfidenceScore(data *ProjectStructure, validation *ValidationResult) float64
```

## Data Enricher Module (internal/transformers/data_enricher.go)

### DataEnricher Implementation
```go
type DataEnricher struct {
    logger *Logger
}

func NewDataEnricher() *DataEnricher
func (de *DataEnricher) EnrichData(transformedData *ProjectStructure) *ProjectStructure
func (de *DataEnricher) calculateDerivedFields(data *ProjectStructure) map[string]interface{}
func (de *DataEnricher) calculateDataQuality(data *ProjectStructure) float64
func (de *DataEnricher) calculateComplexityMetrics(data *ProjectStructure) map[string]int
func (de *DataEnricher) calculateDuration(startDate, endDate string) int
```

## Data Validation Module (internal/validators/data_validator.go)

### Dependencies
- `github.com/go-playground/validator/v10` - for field-level validation
- `golang.org/x/text` - for text processing

### Core Types
```go
type ValidationResult struct {
    IsValid         bool                   `json:"is_valid"`
    Issues          []string               `json:"issues"`
    QualityScore    float64                `json:"quality_score"`
    Warnings        []string               `json:"warnings"`
    Suggestions     []string               `json:"suggestions"`
    ValidationStages map[string]interface{} `json:"validation_stages"`
}

type ValidationPipeline struct {
    DocumentValidator    *DocumentValidator
    StructureValidator   *StructureValidator
    ConsistencyValidator *ConsistencyValidator
    ErrorHandler         *ErrorHandler
    logger               *Logger
}
```

### Document Validator (internal/validators/document_validator.go)
```go
type DocumentValidator struct{}

func (dv *DocumentValidator) ValidateDocumentContent(content, docType string) *ValidationResult
```

### Structure Validator (internal/validators/structure_validator.go)
```go
type StructureValidator struct {
    requiredFields map[string][]string
}

func NewStructureValidator() *StructureValidator
func (sv *StructureValidator) ValidateStructure(structure *ProjectStructure) *ValidationResult
func (sv *StructureValidator) calculateCompletenessScore(structure *ProjectStructure) float64
```

### Consistency Validator (internal/validators/consistency_validator.go)
```go
type ConsistencyValidator struct{}

func (cv *ConsistencyValidator) ValidateConsistency(structure *ProjectStructure) *ValidationResult
```

### Validation Pipeline (internal/validators/data_validator.go)
```go
func NewValidationPipeline() *ValidationPipeline
func (vp *ValidationPipeline) ValidateComplete(data map[string]interface{}, documentPath string) *ValidationResult
func (vp *ValidationPipeline) calculateConfidenceAdjustment(issues, warnings []string) float64
```

## Custom Validation Tags

### Date Format Validation
```go
func validateDateFormat(fl validator.FieldLevel) bool
func validateDateAfterStart(fl validator.FieldLevel) bool
```

### Register Custom Validators
```go
func registerCustomValidators(v *validator.Validate) error
```

## Text Processing Utilities

### Date Processing
```go
type DateProcessor struct {
    datePatterns []string
}

func NewDateProcessor() *DateProcessor
func (dp *DateProcessor) NormalizeDate(dateValue interface{}) *string
func (dp *DateProcessor) ParseDate(dateStr string) (*time.Time, error)
func (dp *DateProcessor) FormatDate(date time.Time) string
```

### Text Normalization
```go
type TextNormalizer struct{}

func (tn *TextNormalizer) NormalizeText(text interface{}) string
func (tn *TextNormalizer) NormalizeStatus(status interface{}) string
func (tn *TextNormalizer) CleanText(text string) string
```

## Error Handling Integration

### Transformation Errors
```go
type TransformationError struct {
    Category ErrorCategory
    Message  string
    Details  map[string]interface{}
}

func NewTransformationError(message string, details map[string]interface{}) *TransformationError
```

### Validation Errors
```go
type ValidationError struct {
    Category ErrorCategory
    Message  string
    Field    string
    Value    interface{}
}

func NewValidationError(message, field string, value interface{}) *ValidationError
```

## Performance Considerations

### Memory Management
- Use pointers for large data structures
- Implement proper garbage collection considerations
- Avoid unnecessary data copying

### Processing Efficiency
- Optimize JSON parsing operations
- Cache validation results when appropriate
- Use efficient string operations

### Concurrency
- Make transformer and validator components thread-safe
- Consider parallel validation for independent checks
- Use appropriate synchronization primitives

## Validation Rules

### Document Content Validation
- Check for minimum content length
- Detect document corruption (null bytes, etc.)
- Verify expected language content (Cyrillic detection)

### Structure Validation
- Required fields presence
- Data type validation
- Minimum structure requirements (phases, tasks)

### Consistency Validation
- Duplicate ID detection
- Date relationship validation
- Dependency reference validation
- Business rule validation

## Data Quality Scoring

### Quality Metrics
- Completeness score (ratio of filled required fields)
- Consistency score (ratio of passed validation checks)
- Confidence adjustment based on validation results

### Scoring Algorithm
- Base score of 1.0
- Deduct points for validation errors
- Deduct points for warnings
- Ensure scores stay within 0.0-1.0 range

## Testing Strategy

### Unit Tests
- Test individual transformation functions
- Test validation logic for each field
- Test date parsing and normalization
- Test error handling scenarios

### Integration Tests
- Test full transformation pipeline
- Test validation pipeline with real data
- Test edge cases and malformed inputs
- Test performance with large datasets

### Test Data
- Valid project structures
- Invalid project structures with various errors
- Edge cases (empty fields, malformed dates)
- Large structures to test performance

## Configuration Integration

### Validation Configuration
- Configurable validation rules
- Adjustable quality thresholds
- Custom validation patterns

### Error Tolerance
- Configurable error thresholds
- Partial result acceptance criteria
- Validation strictness levels

This comprehensive plan covers all aspects of data transformation and validation in Go, ensuring robust handling of LLM responses and quality assurance of extracted project data.