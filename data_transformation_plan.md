# Data Transformation Pipeline Plan

## Overview
This document outlines the implementation plan for the data transformation pipeline that converts LLM responses into standardized JSON structures for the ЖЦП parsing module. The pipeline will handle validation, normalization, and transformation of AI-generated data into the required project structure format.

## Requirements

### Functional Requirements
- Transform LLM JSON responses to standardized project structure schema
- Validate extracted data against defined JSON schema
- Normalize inconsistent data formats (dates, names, etc.)
- Handle missing or incomplete information gracefully
- Calculate confidence scores for transformed data
- Provide detailed error reporting for validation failures

### Technical Requirements
- Python 3.8+ compatibility
- Support for streaming transformation for large responses
- Memory-efficient processing
- Detailed logging and error reporting
- Extensible architecture for future schema changes

## Architecture Overview

### Core Transformation Pipeline
```python
from typing import Dict, Any, Optional, List, Tuple
import json
import logging
from datetime import datetime
from dataclasses import dataclass
import pydantic
from pydantic import BaseModel, Field, validator
from enum import Enum

class TransformationStatus(Enum):
    SUCCESS = "success"
    PARTIAL = "partial"
    FAILED = "failed"
    VALIDATION_ERROR = "validation_error"

@dataclass
class TransformationResult:
    """Result of data transformation"""
    transformed_data: Optional[Dict[str, Any]] = None
    status: TransformationStatus = TransformationStatus.FAILED
    confidence_score: float = 0.0
    validation_errors: List[str] = None
    processing_notes: List[str] = None
    tokens_used: Dict[str, int] = None

class ProjectStructure(BaseModel):
    """Pydantic model for validated project structure"""
    title: str
    description: str
    phases: List['Phase'] = Field(default_factory=list)
    metadata: Dict[str, Any] = Field(default_factory=dict)

class Phase(BaseModel):
    """Pydantic model for project phase"""
    id: str
    name: str
    description: str = ""
    start_date: Optional[str] = None
    end_date: Optional[str] = None
    tasks: List['Task'] = Field(default_factory=list)

class Task(BaseModel):
    """Pydantic model for project task"""
    id: str
    name: str
    description: str = ""
    start_date: Optional[str] = None
    end_date: Optional[str] = None
    responsible_persons: List['ResponsiblePerson'] = Field(default_factory=list)
    dependencies: List[str] = Field(default_factory=list)
    status: str = "planned"

class ResponsiblePerson(BaseModel):
    """Pydantic model for responsible person"""
    name: str
    role: str = ""
    contact: str = ""
```

### Data Transformer Class
```python
import re
from typing import Union
import dateutil.parser

class DataTransformer:
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.date_patterns = [
            r'\d{2}\.\d{2}\.\d{4}',  # DD.MM.YYYY
            r'\d{4}-\d{2}-\d{2}',    # YYYY-MM-DD
            r'\d{2}/\d{2}/\d{4}',    # MM/DD/YYYY
        ]
    
    def transform(self, llm_response: Union[str, Dict]) -> TransformationResult:
        """
        Transform LLM response to standardized project structure
        
        Args:
            llm_response: Raw LLM response (string or dict)
            
        Returns:
            TransformationResult with transformed data and status
        """
        result = TransformationResult()
        
        try:
            # Parse response if it's a string
            if isinstance(llm_response, str):
                try:
                    response_dict = json.loads(llm_response)
                except json.JSONDecodeError as e:
                    self.logger.error(f"Invalid JSON in LLM response: {str(e)}")
                    result.status = TransformationStatus.FAILED
                    result.validation_errors = [f"Invalid JSON: {str(e)}"]
                    return result
            else:
                response_dict = llm_response
            
            # Extract project structure
            project_data = response_dict.get('project', response_dict)
            
            # Normalize and validate data
            normalized_data = self._normalize_data(project_data)
            
            # Validate against schema
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
            
        except Exception as e:
            self.logger.error(f"Transformation error: {str(e)}")
            result.status = TransformationStatus.FAILED
            result.validation_errors = [f"Transformation error: {str(e)}"]
        
        return result
    
    def _normalize_data(self, raw_data: Dict[str, Any]) -> Dict[str, Any]:
        """Normalize raw data to standard format"""
        normalized = {}
        
        # Normalize top-level fields
        normalized['title'] = self._normalize_text(raw_data.get('title', ''))
        normalized['description'] = self._normalize_text(raw_data.get('description', ''))
        
        # Normalize phases
        raw_phases = raw_data.get('phases', raw_data.get('phases', []))
        normalized['phases'] = self._normalize_phases(raw_phases)
        
        # Add metadata
        normalized['metadata'] = raw_data.get('metadata', {})
        
        return normalized
    
    def _normalize_phases(self, raw_phases: List[Dict]) -> List[Dict]:
        """Normalize phases data"""
        normalized_phases = []
        
        for i, raw_phase in enumerate(raw_phases):
            if not isinstance(raw_phase, dict):
                continue
            
            normalized_phase = {
                'id': raw_phase.get('id', f'phase_{i+1}'),
                'name': self._normalize_text(raw_phase.get('name', f'Phase {i+1}')),
                'description': self._normalize_text(raw_phase.get('description', '')),
                'start_date': self._normalize_date(raw_phase.get('start_date')),
                'end_date': self._normalize_date(raw_phase.get('end_date')),
                'tasks': self._normalize_tasks(raw_phase.get('tasks', []), f'phase_{i+1}')
            }
            
            normalized_phases.append(normalized_phase)
        
        return normalized_phases
    
    def _normalize_tasks(self, raw_tasks: List[Dict], phase_id: str) -> List[Dict]:
        """Normalize tasks data"""
        normalized_tasks = []
        
        for i, raw_task in enumerate(raw_tasks):
            if not isinstance(raw_task, dict):
                continue
            
            task_id = raw_task.get('id', f'{phase_id}_task_{i+1}')
            
            normalized_task = {
                'id': task_id,
                'name': self._normalize_text(raw_task.get('name', f'Task {i+1}')),
                'description': self._normalize_text(raw_task.get('description', '')),
                'start_date': self._normalize_date(raw_task.get('start_date')),
                'end_date': self._normalize_date(raw_task.get('end_date')),
                'responsible_persons': self._normalize_responsibles(
                    raw_task.get('responsible_persons', [])
                ),
                'dependencies': raw_task.get('dependencies', []),
                'status': self._normalize_status(raw_task.get('status', 'planned'))
            }
            
            normalized_tasks.append(normalized_task)
        
        return normalized_tasks
    
    def _normalize_responsibles(self, raw_responsibles: List[Dict]) -> List[Dict]:
        """Normalize responsible persons data"""
        normalized_responsibles = []
        
        for raw_resp in raw_responsibles:
            if not isinstance(raw_resp, dict):
                continue
            
            normalized_resp = {
                'name': self._normalize_text(raw_resp.get('name', '')),
                'role': self._normalize_text(raw_resp.get('role', '')),
                'contact': self._normalize_text(raw_resp.get('contact', ''))
            }
            
            if normalized_resp['name']:  # Only include if name exists
                normalized_responsibles.append(normalized_resp)
        
        return normalized_responsibles
    
    def _normalize_date(self, date_value: Any) -> Optional[str]:
        """Normalize date to YYYY-MM-DD format"""
        if not date_value:
            return None
        
        if isinstance(date_value, str):
            # Handle various date formats
            date_str = date_value.strip()
            
            # Try to parse the date
            try:
                parsed_date = dateutil.parser.parse(date_str, dayfirst=True)
                return parsed_date.strftime('%Y-%m-%d')
            except:
                # If parsing fails, try regex patterns
                for pattern in self.date_patterns:
                    match = re.search(pattern, date_str)
                    if match:
                        date_part = match.group(0)
                        try:
                            parsed_date = dateutil.parser.parse(date_part, dayfirst=True)
                            return parsed_date.strftime('%Y-%m-%d')
                        except:
                            continue
        
        elif isinstance(date_value, (int, float)):
            # Handle timestamp
            try:
                parsed_date = datetime.fromtimestamp(date_value)
                return parsed_date.strftime('%Y-%m-%d')
            except:
                pass
        
        return None
    
    def _normalize_text(self, text: Any) -> str:
        """Normalize text content"""
        if text is None:
            return ""
        
        text = str(text).strip()
        # Remove excessive whitespace and normalize line breaks
        text = re.sub(r'\s+', ' ', text)
        return text
    
    def _normalize_status(self, status: Any) -> str:
        """Normalize task status"""
        if not status:
            return "planned"
        
        status_str = str(status).lower().strip()
        
        # Map various status representations to standard values
        status_mapping = {
            'planned': ['planned', 'planning', 'not started', 'в плане', 'запланировано'],
            'in_progress': ['in_progress', 'in progress', 'progress', 'в работе', 'выполняется'],
            'completed': ['completed', 'done', 'finished', 'завершено', 'выполнено', 'complete']
        }
        
        for standard_status, variations in status_mapping.items():
            if status_str in variations or standard_status in status_str:
                return standard_status
        
        return "planned"  # Default to planned
    
    def _validate_data(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Validate data against schema"""
        validation_result = {
            'is_valid': True,
            'errors': [],
            'warnings': [],
            'notes': []
        }
        
        try:
            # Validate using Pydantic model
            project_structure = ProjectStructure(**data)
            
            # Additional business logic validation
            business_validation = self._business_validation(data)
            validation_result['errors'].extend(business_validation['errors'])
            validation_result['warnings'].extend(business_validation['warnings'])
            
        except pydantic.ValidationError as e:
            validation_result['is_valid'] = False
            validation_result['errors'].extend([str(error) for error in e.errors()])
        except Exception as e:
            validation_result['is_valid'] = False
            validation_result['errors'].append(f"Validation error: {str(e)}")
        
        if validation_result['errors']:
            validation_result['is_valid'] = False
        
        return validation_result
    
    def _business_validation(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Perform business logic validation"""
        result = {'errors': [], 'warnings': []}
        
        # Validate date relationships
        phases = data.get('phases', [])
        for phase in phases:
            phase_start = phase.get('start_date')
            phase_end = phase.get('end_date')
            
            if phase_start and phase_end:
                try:
                    start_date = datetime.strptime(phase_start, '%Y-%m-%d')
                    end_date = datetime.strptime(phase_end, '%Y-%m-%d')
                    
                    if start_date > end_date:
                        result['errors'].append(
                            f"Phase '{phase['name']}': Start date is after end date"
                        )
                except ValueError:
                    pass  # Date validation handled by schema
            
            # Validate tasks within phase
            tasks = phase.get('tasks', [])
            for task in tasks:
                task_start = task.get('start_date')
                task_end = task.get('end_date')
                
                if task_start and task_end:
                    try:
                        start_date = datetime.strptime(task_start, '%Y-%m-%d')
                        end_date = datetime.strptime(task_end, '%Y-%m-%d')
                        
                        if start_date > end_date:
                            result['errors'].append(
                                f"Task '{task['name']}': Start date is after end date"
                            )
                    except ValueError:
                        pass
        
        # Validate dependency references
        for phase in phases:
            tasks = phase.get('tasks', [])
            task_ids = {task['id'] for task in tasks}
            
            for task in tasks:
                dependencies = task.get('dependencies', [])
                invalid_deps = [dep for dep in dependencies if dep not in task_ids]
                if invalid_deps:
                    result['errors'].append(
                        f"Task '{task['name']}' has invalid dependencies: {invalid_deps}"
                    )
        
        return result
    
    def _create_partial_transformation(self, data: Dict[str, Any], 
                                     validation_result: Dict[str, Any]) -> Optional[Dict[str, Any]]:
        """Create partial transformation when full validation fails"""
        # For now, return None to indicate complete failure
        # In a more sophisticated implementation, this would attempt
        # to fix validation errors and create a partially valid structure
        return None
    
    def _calculate_confidence_score(self, data: Dict[str, Any], 
                                  validation_result: Dict[str, Any]) -> float:
        """Calculate confidence score based on data quality"""
        score = 1.0
        
        # Reduce score based on validation issues
        error_count = len(validation_result.get('errors', []))
        warning_count = len(validation_result.get('warnings', []))
        
        # Significant penalties for validation errors
        score -= error_count * 0.2
        score -= warning_count * 0.05
        
        # Consider data completeness
        if not data.get('title'):
            score -= 0.1
        if not data.get('description'):
            score -= 0.05
        
        phases = data.get('phases', [])
        if not phases:
            score -= 0.3
        
        # Consider task completeness
        total_tasks = sum(len(phase.get('tasks', [])) for phase in phases)
        if total_tasks == 0:
            score -= 0.2
        
        # Ensure score stays within bounds
        return max(0.0, min(1.0, score))
```

## Advanced Transformation Features

### Data Enrichment
```python
class DataEnricher:
    """Enriches extracted data with additional context"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
    
    def enrich_data(self, transformed_data: Dict[str, Any]) -> Dict[str, Any]:
        """Add computed fields and derived information"""
        enriched = transformed_data.copy()
        
        # Add computed metadata
        enriched['metadata'] = enriched.get('metadata', {})
        enriched['metadata'].update({
            'calculated_fields': self._calculate_derived_fields(transformed_data),
            'data_quality_score': self._calculate_data_quality(transformed_data),
            'complexity_metrics': self._calculate_complexity_metrics(transformed_data)
        })
        
        return enriched
    
    def _calculate_derived_fields(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Calculate derived fields from the data"""
        phases = data.get('phases', [])
        
        # Calculate project timeline
        all_dates = []
        for phase in phases:
            if phase.get('start_date'):
                all_dates.append(phase['start_date'])
            if phase.get('end_date'):
                all_dates.append(phase['end_date'])
            
            for task in phase.get('tasks', []):
                if task.get('start_date'):
                    all_dates.append(task['start_date'])
                if task.get('end_date'):
                    all_dates.append(task['end_date'])
        
        if all_dates:
            all_dates.sort()
            return {
                'project_start': all_dates[0] if all_dates else None,
                'project_end': all_dates[-1] if all_dates else None,
                'total_duration_days': self._calculate_duration(all_dates[0], all_dates[-1]) if len(all_dates) >= 2 else 0
            }
        
        return {}
    
    def _calculate_duration(self, start_date: str, end_date: str) -> int:
        """Calculate duration in days between two dates"""
        try:
            start = datetime.strptime(start_date, '%Y-%m-%d')
            end = datetime.strptime(end_date, '%Y-%m-%d')
            return (end - start).days
        except:
            return 0
    
    def _calculate_data_quality(self, data: Dict[str, Any]) -> float:
        """Calculate data quality score"""
        total_fields = 0
        filled_fields = 0
        
        # Check top-level fields
        for field in ['title', 'description']:
            total_fields += 1
            if data.get(field):
                filled_fields += 1
        
        # Check phases
        phases = data.get('phases', [])
        for phase in phases:
            for field in ['name', 'description']:
                total_fields += 1
                if phase.get(field):
                    filled_fields += 1
            
            # Check tasks
            tasks = phase.get('tasks', [])
            for task in tasks:
                for field in ['name', 'description']:
                    total_fields += 1
                    if task.get(field):
                        filled_fields += 1
        
        if total_fields == 0:
            return 0.0
        
        return filled_fields / total_fields
    
    def _calculate_complexity_metrics(self, data: Dict[str, Any]) -> Dict[str, int]:
        """Calculate project complexity metrics"""
        phases = data.get('phases', [])
        
        return {
            'total_phases': len(phases),
            'total_tasks': sum(len(phase.get('tasks', [])) for phase in phases),
            'total_responsibles': sum(
                len(task.get('responsible_persons', [])) 
                for phase in phases 
                for task in phase.get('tasks', [])
            )
        }
```

## Error Handling and Validation

### Validation Pipeline
```python
class ValidationPipeline:
    """Comprehensive validation pipeline"""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.validators = [
            self._schema_validation,
            self._business_logic_validation,
            self._consistency_validation,
            self._completeness_validation
        ]
    
    def validate(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Run comprehensive validation"""
        results = {
            'overall_valid': True,
            'validation_results': {},
            'errors': [],
            'warnings': [],
            'suggestions': []
        }
        
        for validator in self.validators:
            validator_result = validator(data)
            results['validation_results'][validator.__name__] = validator_result
            
            if not validator_result['valid']:
                results['overall_valid'] = False
                results['errors'].extend(validator_result.get('errors', []))
            
            results['warnings'].extend(validator_result.get('warnings', []))
            results['suggestions'].extend(validator_result.get('suggestions', []))
        
        return results
    
    def _schema_validation(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Validate against JSON schema"""
        try:
            ProjectStructure(**data)
            return {'valid': True, 'errors': [], 'warnings': [], 'suggestions': []}
        except pydantic.ValidationError as e:
            errors = [str(error) for error in e.errors()]
            return {
                'valid': False,
                'errors': errors,
                'warnings': [],
                'suggestions': []
            }
    
    def _business_logic_validation(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Validate business logic constraints"""
        errors = []
        warnings = []
        suggestions = []
        
        # Check date logic
        phases = data.get('phases', [])
        for phase in phases:
            start = phase.get('start_date')
            end = phase.get('end_date')
            
            if start and end:
                try:
                    start_date = datetime.strptime(start, '%Y-%m-%d')
                    end_date = datetime.strptime(end, '%Y-%m-%d')
                    
                    if start_date > end_date:
                        errors.append(f"Phase '{phase['name']}' has end date before start date")
                except ValueError:
                    pass
        
        return {
            'valid': len(errors) == 0,
            'errors': errors,
            'warnings': warnings,
            'suggestions': suggestions
        }
    
    def _consistency_validation(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Validate data consistency"""
        errors = []
        warnings = []
        suggestions = []
        
        # Check for duplicate IDs
        phase_ids = set()
        task_ids = set()
        
        for phase in data.get('phases', []):
            phase_id = phase.get('id')
            if phase_id in phase_ids:
                errors.append(f"Duplicate phase ID: {phase_id}")
            else:
                phase_ids.add(phase_id)
            
            for task in phase.get('tasks', []):
                task_id = task.get('id')
                if task_id in task_ids:
                    errors.append(f"Duplicate task ID: {task_id}")
                else:
                    task_ids.add(task_id)
        
        return {
            'valid': len(errors) == 0,
            'errors': errors,
            'warnings': warnings,
            'suggestions': suggestions
        }
    
    def _completeness_validation(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Validate data completeness"""
        errors = []
        warnings = []
        suggestions = []
        
        # Check for essential information
        if not data.get('title'):
            warnings.append("Project title is missing")
        
        if not data.get('description'):
            warnings.append("Project description is missing")
        
        phases = data.get('phases', [])
        if not phases:
            errors.append("No phases found in project structure")
        else:
            for phase in phases:
                if not phase.get('name'):
                    warnings.append(f"Phase is missing a name")
                
                tasks = phase.get('tasks', [])
                if not tasks:
                    warnings.append(f"Phase '{phase.get('name')}' has no tasks")
        
        return {
            'valid': len(errors) == 0,
            'errors': errors,
            'warnings': warnings,
            'suggestions': suggestions
        }
```

## Integration with Main Pipeline

```python
class ЖЦПParser:
    def __init__(self):
        # ... existing initialization code ...
        
        # Initialize transformation components
        self.data_transformer = DataTransformer()
        self.data_enricher = DataEnricher()
        self.validation_pipeline = ValidationPipeline()
        
        self.logger = logging.getLogger(__name__)
    
    async def extract_and_transform(self, document_path: str) -> TransformationResult:
        """Complete pipeline: extract, transform, validate"""
        
        # Step 1: Extract project structure from document
        extraction_result = await self.extract_project_structure(document_path)
        
        # Step 2: Transform LLM response to structured data
        llm_content = extraction_result['project_structure']
        transformation_result = self.data_transformer.transform(llm_content)
        
        if transformation_result.status in [TransformationStatus.SUCCESS, TransformationStatus.PARTIAL]:
            # Step 3: Enrich the data
            enriched_data = self.data_enricher.enrich_data(
                transformation_result.transformed_data
            )
            transformation_result.transformed_data = enriched_data
            
            # Step 4: Run comprehensive validation
            validation_results = self.validation_pipeline.validate(enriched_data)
            
            if not validation_results['overall_valid'] and transformation_result.status == TransformationStatus.SUCCESS:
                transformation_result.status = TransformationStatus.PARTIAL
                transformation_result.processing_notes.extend(
                    validation_results['warnings']
                )
        
        # Add extraction metadata to transformation result
        if 'extraction_metadata' in extraction_result:
            if transformation_result.processing_notes is None:
                transformation_result.processing_notes = []
            transformation_result.processing_notes.append(
                f"Original extraction confidence: {extraction_result['extraction_metadata'].get('confidence', 0)}"
            )
        
        return transformation_result
```

## Performance Optimization

### Streaming Transformation
```python
import asyncio
from typing import AsyncGenerator

class StreamingTransformer:
    """Handles streaming transformation for large documents"""
    
    def __init__(self, chunk_size: int = 1000):
        self.chunk_size = chunk_size
        self.transformer = DataTransformer()
    
    async def transform_streaming(self, 
                                llm_response_stream: AsyncGenerator[str, None]) -> TransformationResult:
        """Transform streaming LLM response"""
        full_response = ""
        
        async for chunk in llm_response_stream:
            full_response += chunk
            # Process in chunks if needed for very large responses
            if len(full_response) > self.chunk_size * 10:
                # Validate partial response
                pass
        
        return self.transformer.transform(full_response)
```

## Testing Strategy

### Unit Tests
- Test each transformation function individually
- Test validation logic with various error scenarios
- Test date normalization with different formats
- Test confidence scoring algorithms
- Test error handling and fallback mechanisms

### Integration Tests
- Test end-to-end transformation pipeline
- Test with various LLM response formats
- Test performance with large documents
- Test edge cases and malformed responses

This data transformation pipeline provides a robust, extensible system for converting LLM responses into standardized project structure data while ensuring data quality and consistency.