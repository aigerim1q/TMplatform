# Main Application Structure Plan for Go Implementation

## Overview

This document outlines the plan for implementing the main application structure in Go, which will orchestrate all components of the ЖЦП Parser system.

## Main Parser Structure (internal/parser/parser.go)

### ZhcpParser Implementation
```go
type ZhcpParser struct {
    config           *Config
    pdfExtractor     *PDFExtractor
    pdfValidator     *PDFValidator
    docxExtractor    *DOCXExtractor
    docxValidator    *DOCXValidator
    textPreprocessor *TextPreprocessor
    llmManager       *LLMManager
    promptManager    *PromptManager
    dataTransformer  *DataTransformer
    dataEnricher     *DataEnricher
    validationPipeline *ValidationPipeline
    errorHandler     *ErrorHandler
    logger           *Logger
    mu               sync.RWMutex  // For thread safety
}
```

### Constructor
```go
func NewZhcpParser(config *Config) (*ZhcpParser, error)
func (p *ZhcpParser) initializeComponents() error
```

### Core Method - ParseDocument
```go
type ParseResult struct {
    Success           bool                    `json:"success"`
    ProjectStructure  *ProjectStructure       `json:"project_structure,omitempty"`
    ExtractionMetadata ExtractionMetadata     `json:"extraction_metadata"`
    ValidationError  []string                `json:"validation_errors,omitempty"`
    ProcessingNotes   []string                `json:"processing_notes,omitempty"`
    Error            *ErrorInfo              `json:"error,omitempty"`
}

type ExtractionMetadata struct {
    Confidence       float64                 `json:"confidence"`
    Status           string                  `json:"status"`
    ProcessingTime   float64                 `json:"processing_time"`
    ValidationResults *ValidationResult       `json:"validation_results,omitempty"`
}

func (p *ZhcpParser) ParseDocument(documentPath string, validate bool, enrich bool) (*ParseResult, error)
```

## Main Parsing Workflow

### ParseDocument Implementation
```go
func (p *ZhcpParser) ParseDocument(documentPath string, validate bool, enrich bool) (*ParseResult, error) {
    startTime := time.Now()
    
    // 1. Determine document type and validate
    docType, err := p.getDocumentType(documentPath)
    if err != nil {
        return p.createErrorResult(err, documentPath, startTime)
    }
    
    // 2. Validate document based on type
    if docType == "pdf" {
        validation := p.pdfValidator.ValidatePDF(documentPath)
        if !validation.IsValid {
            return p.createErrorResult(
                NewParsingError("PDF validation failed: "+strings.Join(validation.Errors, ", "), documentPath, nil),
                documentPath, startTime)
        }
    } else { // docx
        validation := p.docxValidator.ValidateDOCX(documentPath)
        if !validation.IsValid {
            return p.createErrorResult(
                NewParsingError("DOCX validation failed: "+strings.Join(validation.Errors, ", "), documentPath, nil),
                documentPath, startTime)
        }
    }
    
    // 3. Extract content based on document type
    var extractionResult *ExtractionResult
    if docType == "pdf" {
        extractionResult, err = p.parsePDF(documentPath)
    } else {
        extractionResult, err = p.parseDOCX(documentPath)
    }
    if err != nil {
        return p.createErrorResult(err, documentPath, startTime)
    }
    
    // 4. Validate extracted content
    contentValidation := p.validationPipeline.DocumentValidator.ValidateDocumentContent(
        extractionResult.Text, docType)
    if !contentValidation.IsValid && len(contentValidation.Issues) > 0 {
        p.logger.Warn("Content validation issues", "issues", contentValidation.Issues)
    }
    
    // 5. Create extraction prompt
    jsonSchema := p.getProjectJSONSchema()
    prompt, err := p.promptManager.CreateExtractionPrompt(extractionResult.Text, jsonSchema)
    if err != nil {
        return p.createErrorResult(err, documentPath, startTime)
    }
    
    // 6. Generate response from LLM
    llmResponse, err := p.llmManager.GenerateWithFallback(context.Background(), prompt, GenerationOptions{
        Temperature: 0.1,
        MaxTokens:   4096,
    })
    if err != nil {
        return p.createErrorResult(err, documentPath, startTime)
    }
    
    // 7. Transform LLM response to structured data
    transformationResult := p.dataTransformer.Transform(llmResponse.Content)
    
    if transformationResult.Status == TransformationStatusSuccess || 
       transformationResult.Status == TransformationStatusPartial {
        
        // 8. Enrich the data if requested
        if enrich && transformationResult.TransformedData != nil {
            transformationResult.TransformedData = p.dataEnricher.EnrichData(transformationResult.TransformedData)
        }
        
        // 9. Validate the result if requested
        if validate && transformationResult.TransformedData != nil {
            validationResults := p.validationPipeline.ValidateComplete(map[string]interface{}{
                "project_structure": transformationResult.TransformedData,
                "extracted_content": extractionResult.Text,
                "document_type":     docType,
            }, documentPath)
            
            // Adjust confidence based on validation
            if validationResults != nil && validationResults.ValidationStages != nil {
                adjustment := validationResults.ValidationStages["confidence_adjustment"].(float64)
                transformationResult.ConfidenceScore += adjustment
                transformationResult.ConfidenceScore = math.Max(0.0, math.Min(1.0, transformationResult.ConfidenceScore))
            }
        }
    }
    
    // 10. Prepare final result
    processingTime := time.Since(startTime).Seconds()
    
    result := &ParseResult{
        Success:          transformationResult.Status == TransformationStatusSuccess || 
                         transformationResult.Status == TransformationStatusPartial,
        ProjectStructure: transformationResult.TransformedData,
        ExtractionMetadata: ExtractionMetadata{
            Confidence:     transformationResult.ConfidenceScore,
            Status:         string(transformationResult.Status),
            ProcessingTime: processingTime,
        },
    }
    
    if len(transformationResult.ValidationErrors) > 0 {
        result.ValidationError = transformationResult.ValidationErrors
    }
    
    if len(transformationResult.ProcessingNotes) > 0 {
        result.ProcessingNotes = transformationResult.ProcessingNotes
    }
    
    return result, nil
}
```

## Supporting Methods

### Document Type Detection
```go
func (p *ZhcpParser) getDocumentType(documentPath string) (string, error)
```

### PDF Parsing
```go
func (p *ZhcpParser) parsePDF(pdfPath string) (*ExtractionResult, error)
```

### DOCX Parsing
```go
func (p *ZhcpParser) parseDOCX(docxPath string) (*ExtractionResult, error)
```

### Project JSON Schema
```go
func (p *ZhcpParser) getProjectJSONSchema() map[string]interface{}
```

### Error Result Creation
```go
func (p *ZhcpParser) createErrorResult(err error, documentPath string, startTime time.Time) (*ParseResult, error)
```

## Thread Safety

### Concurrency Considerations
```go
// Use read-write mutex for concurrent access
func (p *ZhcpParser) ParseDocumentConcurrent(documentPath string, validate, enrich bool) (*ParseResult, error)

// Connection pooling for LLM providers
// Resource management for document parsing
```

## Logging Integration

### Structured Logging
```go
func (p *ZhcpParser) logParsingStart(documentPath string, validate, enrich bool)
func (p *ZhcpParser) logParsingProgress(documentPath string, stage string)
func (p *ZhcpParser) logParsingComplete(result *ParseResult, processingTime time.Duration)
```

## Performance Monitoring

### Metrics Collection
```go
type ParserMetrics struct {
    TotalParsed    int64
    Successful     int64
    Failed         int64
    AvgProcessing  float64
    TotalTokens    int64
    ErrorsByType   map[string]int64
    mu             sync.RWMutex
}

func (p *ZhcpParser) collectMetrics(result *ParseResult, processingTime time.Duration)
```

## Error Recovery

### Graceful Degradation
```go
func (p *ZhcpParser) attemptRecovery(err error, documentPath string, validate, enrich bool) (*ParseResult, error)
```

## Resource Management

### Cleanup and Resource Handling
```go
func (p *ZhcpParser) cleanupTempResources()
func (p *ZhcpParser) Close() error  // Close all resources
```

## Command-Line Interface Integration

### CLI-Friendly Interface
```go
type CLIOptions struct {
    Validate     bool
    Enrich       bool
    OutputFormat string
    Verbose      bool
}

func (p *ZhcpParser) ParseDocumentWithCLIOptions(documentPath string, opts CLIOptions) (*ParseResult, error)
```

## Configuration Integration

### Runtime Configuration Updates
```go
func (p *ZhcpParser) UpdateConfig(newConfig *Config) error
func (p *ZhcpParser) GetConfig() *Config
```

## Health Check

### System Health Monitoring
```go
type HealthStatus struct {
    IsHealthy     bool                    `json:"is_healthy"`
    Components    map[string]ComponentStatus `json:"components"`
    Message       string                  `json:"message"`
}

type ComponentStatus struct {
    IsHealthy bool   `json:"is_healthy"`
    Message   string `json:"message"`
}

func (p *ZhcpParser) HealthCheck() *HealthStatus
```

## Main Function (cmd/zhcp-parser/main.go)

### Application Entry Point
```go
func main() {
    // Initialize logger
    logger := initLogger()
    
    // Load configuration
    configManager, err := NewConfigManager("configs/llm_config.yaml")
    if err != nil {
        logger.Fatal("Failed to initialize config manager", "error", err)
    }
    
    config, err := configManager.LoadConfig()
    if err != nil {
        logger.Fatal("Failed to load config", "error", err)
    }
    
    // Initialize parser
    parser, err := NewZhcpParser(config)
    if err != nil {
        logger.Fatal("Failed to initialize parser", "error", err)
    }
    
    // Set up CLI commands using Cobra
    rootCmd := NewRootCommand(parser, logger)
    if err := rootCmd.Execute(); err != nil {
        logger.Error("Command execution failed", "error", err)
        os.Exit(1)
    }
}
```

## Dependency Injection

### Component Initialization
```go
func NewZhcpParser(config *Config) (*ZhcpParser, error) {
    parser := &ZhcpParser{
        config: config,
        logger: initLogger(),  // or inject logger
    }
    
    // Initialize all components with proper error handling
    if err := parser.initializeComponents(); err != nil {
        return nil, fmt.Errorf("failed to initialize components: %w", err)
    }
    
    return parser, nil
}

func (p *ZhcpParser) initializeComponents() error {
    var err error
    
    // Initialize document parsers
    p.pdfExtractor = NewPDFExtractor(p.logger)
    p.pdfValidator = NewPDFValidator()
    p.docxExtractor = NewDOCXExtractor(p.logger)
    p.docxValidator = NewDOCXValidator()
    p.textPreprocessor = NewTextPreprocessor()
    
    // Initialize LLM components
    p.llmManager, err = NewLLMManager(p.config)
    if err != nil {
        return fmt.Errorf("failed to initialize LLM manager: %w", err)
    }
    
    // Initialize AI components
    p.promptManager = NewPromptManager("prompts/")
    
    // Initialize transformers
    p.dataTransformer = NewDataTransformer()
    p.dataEnricher = NewDataEnricher()
    
    // Initialize validators
    p.validationPipeline = NewValidationPipeline()
    
    // Initialize error handler
    p.errorHandler = NewErrorHandler("logs/errors.log", 1000)
    
    return nil
}
```

This main application structure plan provides a comprehensive blueprint for implementing the core ЖЦП Parser in Go, with proper error handling, logging, configuration management, and all the necessary components working together in a cohesive system.