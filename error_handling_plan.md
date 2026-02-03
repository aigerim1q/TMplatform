# Error Handling and Validation Plan

## Overview
This document outlines the comprehensive error handling and validation strategy for the ЖЦП parsing module. The system will implement multiple layers of validation and error handling to ensure robust processing of documents and reliable extraction of project structure information.

## Error Handling Architecture

### Hierarchical Error Classification
```python
from enum import Enum
from typing import Dict, Any, List, Optional
import logging
from dataclasses import dataclass
from datetime import datetime

class ErrorCategory(Enum):
    """Categories of errors in the system"""
    PARSING_ERROR = "parsing_error"
    VALIDATION_ERROR = "validation_error"
    LLM_ERROR = "llm_error"
    TRANSFORMATION_ERROR = "transformation_error"
    CONFIGURATION_ERROR = "configuration_error"
    NETWORK_ERROR = "network_error"
    FILE_ERROR = "file_error"
    BUSINESS_LOGIC_ERROR = "business_logic_error"

class ErrorSeverity(Enum):
    """Severity levels for errors"""
    INFO = "info"
    WARNING = "warning"
    ERROR = "error"
    CRITICAL = "critical"

@dataclass
class ErrorInfo:
    """Comprehensive error information"""
    error_id: str
    category: ErrorCategory
    severity: ErrorSeverity
    message: str
    details: Dict[str, Any]
    timestamp: datetime
    document_path: Optional[str] = None
    component: Optional[str] = None
    traceback: Optional[str] = None
    suggested_action: Optional[str] = None

class ЖЦПError(Exception):
    """Base exception class for ЖЦП parser"""
    def __init__(self, message: str, category: ErrorCategory, 
                 details: Optional[Dict[str, Any]] = None):
        super().__init__(message)
        self.category = category
        self.details = details or {}
        self.timestamp = datetime.now()
        self.error_id = f"ERROR_{int(self.timestamp.timestamp())}_{hash(message) % 10000}"

class ParsingError(ЖЦПError):
    """Error during document parsing"""
    def __init__(self, message: str, document_path: str = None, details: Dict[str, Any] = None):
        super().__init__(message, ErrorCategory.PARSING_ERROR, details or {})
        self.document_path = document_path

class ValidationError(ЖЦПError):
    """Error during data validation"""
    def __init__(self, message: str, field: str = None, value: Any = None):
        details = {"field": field, "value": value}
        super().__init__(message, ErrorCategory.VALIDATION_ERROR, details)

class LLMError(ЖЦПError):
    """Error during LLM processing"""
    def __init__(self, message: str, provider: str = None, details: Dict[str, Any] = None):
        if details is None:
            details = {}
        if provider:
            details["provider"] = provider
        super().__init__(message, ErrorCategory.LLM_ERROR, details)
```

### Error Handler Interface
```python
import traceback
import json
from pathlib import Path

class ErrorHandler:
    """Centralized error handling system"""
    
    def __init__(self, log_file: str = "logs/errors.log", 
                 max_errors: int = 1000):
        self.log_file = Path(log_file)
        self.max_errors = max_errors
        self.error_log = []
        self.logger = logging.getLogger(__name__)
        
        # Create log directory if it doesn't exist
        self.log_file.parent.mkdir(parents=True, exist_ok=True)
    
    def handle_error(self, error: Exception, 
                    document_path: str = None,
                    component: str = None) -> ErrorInfo:
        """Handle an error and return error information"""
        error_info = self._create_error_info(error, document_path, component)
        
        # Log the error
        self._log_error(error_info)
        
        # Store in memory (with rotation)
        self.error_log.append(error_info)
        if len(self.error_log) > self.max_errors:
            self.error_log.pop(0)
        
        return error_info
    
    def _create_error_info(self, error: Exception, 
                          document_path: str, 
                          component: str) -> ErrorInfo:
        """Create error information from exception"""
        if isinstance(error, ЖЦПError):
            category = error.category
            message = str(error)
            details = error.details
            severity = self._determine_severity(error.category)
        else:
            category = ErrorCategory.ERROR
            message = f"Unexpected error: {str(error)}"
            details = {
                "original_error_type": type(error).__name__,
                "original_error_message": str(error)
            }
            severity = ErrorSeverity.ERROR
        
        return ErrorInfo(
            error_id=f"ERR_{int(datetime.now().timestamp())}_{hash(str(error)) % 10000}",
            category=category,
            severity=severity,
            message=message,
            details=details,
            timestamp=datetime.now(),
            document_path=document_path,
            component=component,
            traceback=traceback.format_exc() if str(error) != "Unexpected error: None" else None,
            suggested_action=self._get_suggested_action(category)
        )
    
    def _determine_severity(self, category: ErrorCategory) -> ErrorSeverity:
        """Determine severity based on error category"""
        severity_mapping = {
            ErrorCategory.PARSING_ERROR: ErrorSeverity.ERROR,
            ErrorCategory.VALIDATION_ERROR: ErrorSeverity.WARNING,
            ErrorCategory.LLM_ERROR: ErrorSeverity.ERROR,
            ErrorCategory.TRANSFORMATION_ERROR: ErrorSeverity.ERROR,
            ErrorCategory.CONFIGURATION_ERROR: ErrorSeverity.CRITICAL,
            ErrorCategory.NETWORK_ERROR: ErrorSeverity.ERROR,
            ErrorCategory.FILE_ERROR: ErrorSeverity.ERROR,
            ErrorCategory.BUSINESS_LOGIC_ERROR: ErrorSeverity.ERROR
        }
        return severity_mapping.get(category, ErrorSeverity.ERROR)
    
    def _get_suggested_action(self, category: ErrorCategory) -> str:
        """Get suggested action based on error category"""
        actions = {
            ErrorCategory.PARSING_ERROR: "Check document format and encoding",
            ErrorCategory.VALIDATION_ERROR: "Review extracted data structure",
            ErrorCategory.LLM_ERROR: "Verify API configuration and credentials",
            ErrorCategory.TRANSFORMATION_ERROR: "Check data transformation logic",
            ErrorCategory.CONFIGURATION_ERROR: "Review configuration files",
            ErrorCategory.NETWORK_ERROR: "Check network connectivity",
            ErrorCategory.FILE_ERROR: "Verify file permissions and existence",
            ErrorCategory.BUSINESS_LOGIC_ERROR: "Review business logic implementation"
        }
        return actions.get(category, "Review error details and logs")
    
    def _log_error(self, error_info: ErrorInfo):
        """Log error to file"""
        log_entry = {
            "error_id": error_info.error_id,
            "category": error_info.category.value,
            "severity": error_info.severity.value,
            "message": error_info.message,
            "timestamp": error_info.timestamp.isoformat(),
            "document_path": error_info.document_path,
            "component": error_info.component,
            "details": error_info.details
        }
        
        with open(self.log_file, 'a', encoding='utf-8') as f:
            f.write(json.dumps(log_entry, ensure_ascii=False) + '\n')
    
    def get_error_summary(self) -> Dict[str, Any]:
        """Get summary of recent errors"""
        if not self.error_log:
            return {"total_errors": 0, "by_category": {}, "by_severity": {}}
        
        by_category = {}
        by_severity = {}
        total = len(self.error_log)
        
        for error in self.error_log:
            # Count by category
            cat = error.category.value
            by_category[cat] = by_category.get(cat, 0) + 1
            
            # Count by severity
            sev = error.severity.value
            by_severity[sev] = by_severity.get(sev, 0) + 1
        
        return {
            "total_errors": total,
            "by_category": by_category,
            "by_severity": by_severity,
            "recent_errors": [
                {
                    "error_id": e.error_id,
                    "category": e.category.value,
                    "severity": e.severity.value,
                    "message": e.message,
                    "timestamp": e.timestamp.isoformat()
                }
                for e in self.error_log[-10:]  # Last 10 errors
            ]
        }
```

## Validation System

### Multi-Level Validation
```python
from abc import ABC, abstractmethod
from typing import Protocol, runtime_checkable

@runtime_checkable
class Validator(Protocol):
    """Protocol for validation classes"""
    
    def validate(self, data: Any) -> Dict[str, Any]:
        """Validate data and return results"""
        ...

class DocumentValidator:
    """Validation for document parsing results"""
    
    def validate_document_content(self, content: str, doc_type: str) -> Dict[str, Any]:
        """Validate document content quality"""
        results = {
            "is_valid": True,
            "issues": [],
            "quality_score": 1.0,
            "suggestions": []
        }
        
        if not content or len(content.strip()) < 100:
            results["is_valid"] = False
            results["issues"].append("Document content is too short")
            results["quality_score"] = 0.1
        
        if content.count('\x00') > 0:
            results["issues"].append("Document contains null bytes (possible corruption)")
            results["quality_score"] *= 0.5
        
        # Check for Cyrillic content (expected for ЖЦП)
        cyrillic_chars = sum(1 for c in content if '\u0400' <= c <= '\u04FF')
        total_chars = len([c for c in content if c.isalpha()])
        
        if total_chars > 0:
            cyrillic_ratio = cyrillic_chars / total_chars
            if cyrillic_ratio < 0.3:  # Less than 30% Cyrillic
                results["suggestions"].append(
                    "Document may not contain expected Russian content for ЖЦП"
                )
        
        results["quality_score"] = min(1.0, results["quality_score"])
        return results

class StructureValidator:
    """Validation for extracted project structure"""
    
    def __init__(self):
        self.required_fields = {
            'project': ['title', 'phases'],
            'phase': ['id', 'name'],
            'task': ['id', 'name']
        }
    
    def validate_structure(self, structure: Dict[str, Any]) -> Dict[str, Any]:
        """Validate project structure completeness"""
        results = {
            "is_valid": True,
            "issues": [],
            "quality_score": 1.0,
            "warnings": [],
            "suggestions": []
        }
        
        # Validate project level
        project = structure.get('project', {})
        
        if not project:
            results["is_valid"] = False
            results["issues"].append("No project structure found")
            return results
        
        # Check required project fields
        for field in self.required_fields['project']:
            if not project.get(field):
                results["issues"].append(f"Missing required project field: {field}")
                results["is_valid"] = False
        
        # Validate phases
        phases = project.get('phases', [])
        if not phases:
            results["issues"].append("No phases found in project")
            results["is_valid"] = False
        
        for i, phase in enumerate(phases):
            if not isinstance(phase, dict):
                results["issues"].append(f"Phase {i} is not a dictionary")
                continue
            
            # Check required phase fields
            for field in self.required_fields['phase']:
                if not phase.get(field):
                    results["issues"].append(f"Phase {i} missing required field: {field}")
                    results["is_valid"] = False
            
            # Validate tasks within phase
            tasks = phase.get('tasks', [])
            for j, task in enumerate(tasks):
                if not isinstance(task, dict):
                    results["issues"].append(f"Task {j} in phase {i} is not a dictionary")
                    continue
                
                for field in self.required_fields['task']:
                    if not task.get(field):
                        results["issues"].append(f"Task {j} in phase {i} missing required field: {field}")
                        results["is_valid"] = False
        
        # Calculate quality score based on completeness
        results["quality_score"] = self._calculate_completeness_score(structure)
        
        return results
    
    def _calculate_completeness_score(self, structure: Dict[str, Any]) -> float:
        """Calculate completeness score for structure"""
        total_checks = 0
        passed_checks = 0
        
        project = structure.get('project', {})
        
        # Check project fields
        for field in self.required_fields['project']:
            total_checks += 1
            if project.get(field):
                passed_checks += 1
        
        # Check phases and tasks
        phases = project.get('phases', [])
        for phase in phases:
            for field in self.required_fields['phase']:
                total_checks += 1
                if phase.get(field):
                    passed_checks += 1
            
            tasks = phase.get('tasks', [])
            for task in tasks:
                for field in self.required_fields['task']:
                    total_checks += 1
                    if task.get(field):
                        passed_checks += 1
        
        if total_checks == 0:
            return 0.0
        
        return passed_checks / total_checks

class DataConsistencyValidator:
    """Validation for data consistency and business rules"""
    
    def validate_consistency(self, structure: Dict[str, Any]) -> Dict[str, Any]:
        """Validate data consistency and business rules"""
        results = {
            "is_valid": True,
            "issues": [],
            "warnings": [],
            "suggestions": []
        }
        
        project = structure.get('project', {})
        phases = project.get('phases', [])
        
        # Check for duplicate IDs
        phase_ids = set()
        task_ids = set()
        
        for i, phase in enumerate(phases):
            phase_id = phase.get('id')
            if phase_id:
                if phase_id in phase_ids:
                    results["issues"].append(f"Duplicate phase ID: {phase_id}")
                    results["is_valid"] = False
                else:
                    phase_ids.add(phase_id)
            
            tasks = phase.get('tasks', [])
            for j, task in enumerate(tasks):
                task_id = task.get('id')
                if task_id:
                    if task_id in task_ids:
                        results["issues"].append(f"Duplicate task ID: {task_id}")
                        results["is_valid"] = False
                    else:
                        task_ids.add(task_id)
        
        # Validate date relationships
        for i, phase in enumerate(phases):
            start_date = phase.get('start_date')
            end_date = phase.get('end_date')
            
            if start_date and end_date:
                try:
                    start = datetime.strptime(start_date, '%Y-%m-%d')
                    end = datetime.strptime(end_date, '%Y-%m-%d')
                    
                    if start > end:
                        results["issues"].append(
                            f"Phase '{phase.get('name', f'Phase {i}')}': "
                            f"start date ({start_date}) is after end date ({end_date})"
                        )
                        results["is_valid"] = False
                except ValueError:
                    results["issues"].append(
                        f"Phase '{phase.get('name', f'Phase {i}')}': "
                        f"invalid date format"
                    )
                    results["is_valid"] = False
            
            # Validate task dates within phase context
            tasks = phase.get('tasks', [])
            for j, task in enumerate(tasks):
                task_start = task.get('start_date')
                task_end = task.get('end_date')
                
                if task_start and start_date:
                    try:
                        task_start_dt = datetime.strptime(task_start, '%Y-%m-%d')
                        phase_start_dt = datetime.strptime(start_date, '%Y-%m-%d')
                        
                        if task_start_dt < phase_start_dt:
                            results["warnings"].append(
                                f"Task '{task.get('name', f'Task {j}')}': "
                                f"starts before phase begins"
                            )
                    except ValueError:
                        pass
                
                if task_end and end_date:
                    try:
                        task_end_dt = datetime.strptime(task_end, '%Y-%m-%d')
                        phase_end_dt = datetime.strptime(end_date, '%Y-%m-%d')
                        
                        if task_end_dt > phase_end_dt:
                            results["warnings"].append(
                                f"Task '{task.get('name', f'Task {j}')}': "
                                f"ends after phase ends"
                            )
                    except ValueError:
                        pass
        
        # Validate dependency references
        all_task_ids = set()
        for phase in phases:
            for task in phase.get('tasks', []):
                all_task_ids.add(task.get('id'))
        
        for phase in phases:
            for task in phase.get('tasks', []):
                dependencies = task.get('dependencies', [])
                invalid_deps = [dep for dep in dependencies if dep not in all_task_ids]
                
                if invalid_deps:
                    results["issues"].append(
                        f"Task '{task.get('name', task.get('id', 'Unknown'))}' "
                        f"has invalid dependencies: {invalid_deps}"
                    )
                    results["is_valid"] = False
        
        return results
```

## Validation Pipeline

### Comprehensive Validation Pipeline
```python
class ValidationPipeline:
    """Comprehensive validation pipeline with multiple stages"""
    
    def __init__(self):
        self.document_validator = DocumentValidator()
        self.structure_validator = StructureValidator()
        self.consistency_validator = DataConsistencyValidator()
        self.error_handler = ErrorHandler()
        self.logger = logging.getLogger(__name__)
    
    def validate_complete(self, data: Dict[str, Any], 
                         document_path: str = None) -> Dict[str, Any]:
        """Run complete validation pipeline"""
        results = {
            "overall_valid": True,
            "validation_stages": {},
            "issues": [],
            "warnings": [],
            "quality_score": 1.0,
            "confidence_adjustment": 0.0
        }
        
        try:
            # Stage 1: Document content validation
            doc_results = self.document_validator.validate_document_content(
                data.get('extracted_content', ''), 
                data.get('document_type', 'unknown')
            )
            results["validation_stages"]["document_content"] = doc_results
            
            if not doc_results["is_valid"]:
                results["overall_valid"] = False
                results["issues"].extend(doc_results["issues"])
            
            results["warnings"].extend(doc_results.get("suggestions", []))
            
            # Stage 2: Structure validation
            structure_results = self.structure_validator.validate_structure(
                data.get('project_structure', {})
            )
            results["validation_stages"]["structure"] = structure_results
            
            if not structure_results["is_valid"]:
                results["overall_valid"] = False
                results["issues"].extend(structure_results["issues"])
            
            results["warnings"].extend(structure_results["warnings"])
            
            # Stage 3: Consistency validation
            consistency_results = self.consistency_validator.validate_consistency(
                data.get('project_structure', {})
            )
            results["validation_stages"]["consistency"] = consistency_results
            
            if not consistency_results["is_valid"]:
                results["overall_valid"] = False
                results["issues"].extend(consistency_results["issues"])
            
            results["warnings"].extend(consistency_results["warnings"])
            
            # Calculate overall quality score
            scores = [
                doc_results.get("quality_score", 1.0),
                structure_results.get("quality_score", 1.0)
            ]
            results["quality_score"] = sum(scores) / len(scores) if scores else 1.0
            
            # Adjust confidence based on validation results
            results["confidence_adjustment"] = self._calculate_confidence_adjustment(
                results["issues"], results["warnings"]
            )
            
        except Exception as e:
            error_info = self.error_handler.handle_error(
                e, document_path, "ValidationPipeline"
            )
            results["overall_valid"] = False
            results["issues"].append(f"Validation pipeline error: {str(e)}")
            self.logger.error(f"Validation pipeline failed: {error_info.message}")
        
        return results
    
    def _calculate_confidence_adjustment(self, issues: List[str], 
                                       warnings: List[str]) -> float:
        """Calculate confidence adjustment based on validation results"""
        adjustment = 0.0
        
        # Issues significantly reduce confidence
        adjustment -= len(issues) * 0.2
        
        # Warnings moderately reduce confidence
        adjustment -= len(warnings) * 0.05
        
        # Cap adjustment between -0.5 and 0 (don't increase confidence)
        return max(-0.5, min(0.0, adjustment))
```

## Integration with Main Components

### Error Handling in Document Parsing
```python
class PDFExtractor:
    def __init__(self):
        self.error_handler = ErrorHandler()
        self.logger = logging.getLogger(__name__)
    
    def extract_text(self, pdf_path: str) -> Dict[str, Any]:
        """Extract text with comprehensive error handling"""
        try:
            # Validate input
            if not Path(pdf_path).exists():
                raise ParsingError(f"PDF file does not exist: {pdf_path}")
            
            # Perform extraction
            result = self._extract_with_pymupdf(pdf_path)
            
            # Validate result
            validation = self._validate_extraction_result(result)
            if not validation["is_valid"]:
                self.logger.warning(f"PDF extraction validation issues: {validation['issues']}")
            
            return result
            
        except Exception as e:
            error_info = self.error_handler.handle_error(
                e, pdf_path, "PDFExtractor"
            )
            raise ParsingError(
                f"PDF extraction failed: {error_info.message}",
                document_path=pdf_path,
                details=error_info.details
            )

class DOCXExtractor:
    def __init__(self):
        self.error_handler = ErrorHandler()
        self.logger = logging.getLogger(__name__)
    
    def extract_text(self, docx_path: str) -> Dict[str, Any]:
        """Extract text with comprehensive error handling"""
        try:
            # Validate input
            if not Path(docx_path).exists():
                raise ParsingError(f"DOCX file does not exist: {docx_path}")
            
            # Perform extraction
            result = self._extract_with_python_docx(docx_path)
            
            # Validate result
            validation = self._validate_extraction_result(result)
            if not validation["is_valid"]:
                self.logger.warning(f"DOCX extraction validation issues: {validation['issues']}")
            
            return result
            
        except Exception as e:
            error_info = self.error_handler.handle_error(
                e, docx_path, "DOCXExtractor"
            )
            raise ParsingError(
                f"DOCX extraction failed: {error_info.message}",
                document_path=docx_path,
                details=error_info.details
            )
```

### Error Handling in LLM Processing
```python
class LLMManager:
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.providers = {}
        self.error_handler = ErrorHandler()
        self.logger = logging.getLogger(__name__)
        self._initialize_providers()
    
    async def generate_with_fallback(self, prompt: str, **kwargs) -> LLMResponse:
        """Generate response with fallback and error handling"""
        provider_priority = self.config.get('provider_priority', [
            ProviderType.OPENAI.value,
            ProviderType.ANTHROPIC.value,
            ProviderType.OLLAMA.value
        ])
        
        errors = []
        
        for provider_name in provider_priority:
            try:
                provider_type = ProviderType(provider_name)
                provider = self.providers.get(provider_type)
                
                if not provider:
                    continue
                
                self.logger.info(f"Attempting to use {provider_type.value} provider")
                response = await provider.generate(prompt, **kwargs)
                return response
                
            except Exception as e:
                error_info = self.error_handler.handle_error(
                    e, component=f"LLMProvider_{provider_name}"
                )
                errors.append(error_info)
                
                self.logger.warning(
                    f"Provider {provider_name} failed: {error_info.message}"
                )
                
                # Check if this is a configuration error (shouldn't retry)
                if (isinstance(e, ЖЦПError) and 
                    e.category == ErrorCategory.CONFIGURATION_ERROR):
                    break
        
        # If all providers failed, raise a comprehensive error
        error_messages = [err.message for err in errors]
        raise LLMError(
            f"All LLM providers failed: {'; '.join(error_messages[:3])}",  # Limit message length
            details={"provider_errors": [err.message for err in errors]}
        )
```

### Error Handling in Data Transformation
```python
class DataTransformer:
    def __init__(self):
        self.error_handler = ErrorHandler()
        self.logger = logging.getLogger(__name__)
    
    def transform(self, llm_response: Union[str, Dict]) -> TransformationResult:
        """Transform with comprehensive error handling"""
        result = TransformationResult()
        
        try:
            # Parse response
            if isinstance(llm_response, str):
                try:
                    response_dict = json.loads(llm_response)
                except json.JSONDecodeError as e:
                    raise ValidationError(
                        f"Invalid JSON in LLM response: {str(e)}",
                        details={"raw_response_preview": llm_response[:200]}
                    )
            else:
                response_dict = llm_response
            
            # Perform transformation
            normalized_data = self._normalize_data(response_dict)
            
            # Validate result
            validation_result = self._validate_data(normalized_data)
            
            if validation_result['is_valid']:
                result.transformed_data = normalized_data
                result.status = TransformationStatus.SUCCESS
                result.confidence_score = self._calculate_confidence_score(
                    normalized_data, validation_result
                )
            else:
                # Attempt partial transformation
                partial_data = self._create_partial_transformation(
                    normalized_data, validation_result
                )
                
                if partial_data:
                    result.transformed_data = partial_data
                    result.status = TransformationStatus.PARTIAL
                    result.confidence_score = self._calculate_confidence_score(
                        partial_data, validation_result
                    )
                else:
                    result.status = TransformationStatus.VALIDATION_ERROR
                    result.validation_errors = validation_result['errors']
            
            result.processing_notes = validation_result.get('notes', [])
            
        except ValidationError as e:
            error_info = self.error_handler.handle_error(
                e, component="DataTransformer"
            )
            result.status = TransformationStatus.VALIDATION_ERROR
            result.validation_errors = [e.message]
            
        except Exception as e:
            error_info = self.error_handler.handle_error(
                e, component="DataTransformer"
            )
            result.status = TransformationStatus.FAILED
            result.validation_errors = [f"Transformation error: {error_info.message}"]
        
        return result
```

## Main Parser Integration

```python
class ЖЦПParser:
    def __init__(self):
        # Initialize all components
        self.pdf_extractor = PDFExtractor()
        self.docx_extractor = AdvancedDOCXExtractor()
        self.pdf_validator = PDFValidator()
        self.docx_validator = DOCXValidator()
        self.text_preprocessor = TextPreprocessor()
        self.docx_preprocessor = DOCXTextPreprocessor()
        
        self.llm_config = LLMConfig()
        self.llm_manager = LLMManager(self.llm_config.config)
        self.prompt_manager = PromptManager()
        
        self.data_transformer = DataTransformer()
        self.data_enricher = DataEnricher()
        self.validation_pipeline = ValidationPipeline()
        
        # Centralized error handler
        self.error_handler = ErrorHandler()
        self.logger = logging.getLogger(__name__)
    
    async def parse_document(self, document_path: str) -> Dict[str, Any]:
        """Complete document parsing with comprehensive error handling"""
        start_time = datetime.now()
        
        try:
            # Validate document path
            doc_path = Path(document_path)
            if not doc_path.exists():
                raise ParsingError(f"Document does not exist: {document_path}")
            
            # Determine document type and extract content
            if document_path.lower().endswith('.pdf'):
                extraction_result = await self._parse_pdf_for_llm(document_path)
                doc_type = 'pdf'
            elif document_path.lower().endswith('.docx'):
                extraction_result = await self._parse_docx_for_llm(document_path)
                doc_type = 'docx'
            else:
                raise ParsingError(f"Unsupported document format: {document_path}")
            
            # Validate extracted content
            content_validation = self.validation_pipeline.document_validator.validate_document_content(
                extraction_result.get('extracted_content', ''),
                doc_type
            )
            
            if not content_validation["is_valid"]:
                self.logger.warning(f"Content validation issues: {content_validation['issues']}")
            
            # Create extraction prompt
            json_schema = self._get_project_json_schema()
            prompt = self.prompt_manager.create_extraction_prompt(
                document_content=extraction_result['extracted_content'],
                json_schema=json_schema
            )
            
            # Generate response from LLM
            llm_response = await self.llm_manager.generate_with_fallback(
                prompt,
                temperature=0.1,
                max_tokens=4096
            )
            
            # Transform LLM response
            transformation_result = self.data_transformer.transform(
                llm_response.content
            )
            
            if transformation_result.status in [TransformationStatus.SUCCESS, TransformationStatus.PARTIAL]:
                # Enrich the data
                enriched_data = self.data_enricher.enrich_data(
                    transformation_result.transformed_data
                )
                transformation_result.transformed_data = enriched_data
                
                # Run comprehensive validation
                validation_results = self.validation_pipeline.validate_complete(
                    {
                        'project_structure': enriched_data,
                        'extracted_content': extraction_result['extracted_content'],
                        'document_type': doc_type
                    },
                    document_path
                )
                
                # Adjust confidence based on validation
                if validation_results["confidence_adjustment"] != 0:
                    transformation_result.confidence_score += validation_results["confidence_adjustment"]
                    transformation_result.confidence_score = max(
                        0.0, min(1.0, transformation_result.confidence_score)
                    )
            
            # Prepare final result
            final_result = {
                'success': transformation_result.status in [TransformationStatus.SUCCESS, TransformationStatus.PARTIAL],
                'project_structure': transformation_result.transformed_data,
                'extraction_metadata': {
                    'confidence': transformation_result.confidence_score,
                    'status': transformation_result.status.value,
                    'processing_time': (datetime.now() - start_time).total_seconds(),
                    'validation_results': validation_results if 'validation_results' in locals() else None
                }
            }
            
            if transformation_result.validation_errors:
                final_result['validation_errors'] = transformation_result.validation_errors
            
            if transformation_result.processing_notes:
                final_result['processing_notes'] = transformation_result.processing_notes
            
            return final_result
            
        except ЖЦПError as e:
            error_info = self.error_handler.handle_error(
                e, document_path, "ЖЦПParser"
            )
            return {
                'success': False,
                'error': {
                    'message': e.message,
                    'category': e.category.value,
                    'details': e.details,
                    'error_id': error_info.error_id
                }
            }
        
        except Exception as e:
            error_info = self.error_handler.handle_error(
                e, document_path, "ЖЦПParser"
            )
            return {
                'success': False,
                'error': {
                    'message': f"Unexpected error during parsing: {str(e)}",
                    'category': ErrorCategory.ERROR.value,
                    'details': {'traceback': error_info.traceback},
                    'error_id': error_info.error_id
                }
            }
    
    def get_error_summary(self) -> Dict[str, Any]:
        """Get error summary from the error handler"""
        return self.error_handler.get_error_summary()
```

## Error Recovery and Fallback Strategies

### Fallback Mechanisms
```python
class FallbackManager:
    """Manages fallback strategies for error recovery"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def apply_fallback_strategy(self, error_category: ErrorCategory, 
                               original_error: Exception,
                               context: Dict[str, Any]) -> Dict[str, Any]:
        """Apply appropriate fallback strategy based on error type"""
        
        if error_category == ErrorCategory.LLM_ERROR:
            return self._handle_llm_fallback(original_error, context)
        elif error_category == ErrorCategory.PARSING_ERROR:
            return self._handle_parsing_fallback(original_error, context)
        elif error_category == ErrorCategory.TRANSFORMATION_ERROR:
            return self._handle_transformation_fallback(original_error, context)
        else:
            return self._handle_general_fallback(original_error, context)
    
    def _handle_llm_fallback(self, error: Exception, context: Dict[str, Any]) -> Dict[str, Any]:
        """Handle LLM-related fallbacks"""
        result = {
            'fallback_applied': True,
            'new_result': None,
            'confidence_adjustment': -0.3
        }
        
        # Try alternative LLM provider if available
        if context.get('alternative_providers'):
            try:
                alternative_provider = context['alternative_providers'][0]
                # Attempt to use alternative provider
                result['new_result'] = alternative_provider.generate(
                    context['prompt'],
                    **context.get('generation_params', {})
                )
            except:
                result['fallback_applied'] = False
        
        return result
    
    def _handle_parsing_fallback(self, error: Exception, context: Dict[str, Any]) -> Dict[str, Any]:
        """Handle parsing-related fallbacks"""
        result = {
            'fallback_applied': True,
            'new_result': None,
            'confidence_adjustment': -0.2
        }
        
        # Try alternative parsing method
        doc_path = context.get('document_path')
        if doc_path and doc_path.endswith('.pdf'):
            # Try different PDF parsing library
            try:
                # Implementation would use alternative PDF library
                pass
            except:
                result['fallback_applied'] = False
        
        return result
    
    def _handle_transformation_fallback(self, error: Exception, context: Dict[str, Any]) -> Dict[str, Any]:
        """Handle transformation-related fallbacks"""
        result = {
            'fallback_applied': True,
            'new_result': None,
            'confidence_adjustment': -0.4
        }
        
        # Try to extract basic information even if full transformation fails
        raw_content = context.get('raw_content', '')
        if raw_content:
            # Attempt basic extraction using regex or simple parsing
            basic_extraction = self._basic_extraction_fallback(raw_content)
            if basic_extraction:
                result['new_result'] = basic_extraction
        
        return result
    
    def _basic_extraction_fallback(self, content: str) -> Dict[str, Any]:
        """Basic extraction using regex patterns as fallback"""
        import re
        
        # Simple pattern matching for basic project information
        patterns = {
            'phases': r'(?:этап|фаза|phase)\s*[0-9]?\s*[:\-]?\s*([^\n\r]{10,100})',
            'tasks': r'(?:задача|task)\s*[0-9]?\s*[:\-]?\s*([^\n\r]{10,200})',
            'dates': r'(\d{2}\.\d{2}\.\d{4}|\d{4}-\d{2}-\d{2})',
            'responsible': r'(?:ответственный|responsible)\s*[:\-]?\s*([^\n\r]{5,100})'
        }
        
        result = {
            'project': {
                'title': 'Partial Extraction (Fallback)',
                'description': 'Basic extraction due to transformation error',
                'phases': [],
                'metadata': {'extraction_method': 'fallback_regex'}
            }
        }
        
        # Extract phases
        phase_matches = re.findall(patterns['phases'], content, re.IGNORECASE)
        for i, match in enumerate(phase_matches[:5]):  # Limit to 5 phases
            result['project']['phases'].append({
                'id': f'fallback_phase_{i+1}',
                'name': match.strip(),
                'tasks': []
            })
        
        return result if any(result['project']['phases']) else None
    
    def _handle_general_fallback(self, error: Exception, context: Dict[str, Any]) -> Dict[str, Any]:
        """Handle general fallbacks"""
        return {
            'fallback_applied': False,
            'new_result': None,
            'confidence_adjustment': 0.0
        }
```

## Testing Strategy

### Error Handling Tests
```python
import unittest
import tempfile
import os
from unittest.mock import Mock, patch

class TestErrorHandling(unittest.TestCase):
    def setUp(self):
        self.error_handler = ErrorHandler()
    
    def test_parsing_error_handling(self):
        """Test parsing error handling"""
        error = ParsingError("Test parsing error", document_path="/test.pdf")
        error_info = self.error_handler.handle_error(error)
        
        self.assertEqual(error_info.category, ErrorCategory.PARSING_ERROR)
        self.assertEqual(error_info.severity, ErrorSeverity.ERROR)
    
    def test_validation_error_handling(self):
        """Test validation error handling"""
        error = ValidationError("Invalid date format", field="start_date", value="invalid")
        error_info = self.error_handler.handle_error(error)
        
        self.assertEqual(error_info.category, ErrorCategory.VALIDATION_ERROR)
        self.assertEqual(error_info.details["field"], "start_date")
    
    def test_unexpected_error_handling(self):
        """Test handling of unexpected errors"""
        error = ValueError("Unexpected value error")
        error_info = self.error_handler.handle_error(error)
        
        self.assertEqual(error_info.category, ErrorCategory.ERROR)
        self.assertEqual(error_info.severity, ErrorSeverity.ERROR)

class TestValidation(unittest.TestCase):
    def setUp(self):
        self.validator = StructureValidator()
    
    def test_valid_structure(self):
        """Test validation of valid structure"""
        valid_structure = {
            'project': {
                'title': 'Test Project',
                'phases': [
                    {
                        'id': 'phase_1',
                        'name': 'Phase 1',
                        'tasks': [
                            {
                                'id': 'task_1',
                                'name': 'Task 1'
                            }
                        ]
                    }
                ]
            }
        }
        
        result = self.validator.validate_structure(valid_structure)
        self.assertTrue(result["is_valid"])
    
    def test_invalid_structure(self):
        """Test validation of invalid structure"""
        invalid_structure = {
            'project': {
                'phases': []  # Missing required title
            }
        }
        
        result = self.validator.validate_structure(invalid_structure)
        self.assertFalse(result["is_valid"])
        self.assertGreater(len(result["issues"]), 0)
```

This comprehensive error handling and validation plan provides multiple layers of protection and validation for the ЖЦП parsing module, ensuring robust processing of documents and reliable extraction of project structure information.