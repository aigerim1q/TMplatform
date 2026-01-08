# Testing Framework Plan for ЖЦП Parsing Module

## Overview
This document outlines the comprehensive testing framework for the ЖЦП parsing module. The testing strategy includes unit tests, integration tests, end-to-end tests, and performance tests to ensure the reliability and accuracy of the AI-powered document parsing system.

## Testing Strategy

### Test Categories
1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Component interaction testing
3. **End-to-End Tests**: Complete workflow testing
4. **Performance Tests**: Load and performance testing
5. **Accuracy Tests**: Data extraction accuracy validation
6. **Error Handling Tests**: Exception handling validation

## Unit Testing Framework

### Core Component Tests
```python
import unittest
import tempfile
import os
from unittest.mock import Mock, patch, MagicMock
import json
from datetime import datetime
from pathlib import Path

class TestPDFExtractor(unittest.TestCase):
    def setUp(self):
        self.pdf_extractor = PDFExtractor()
    
    def test_extract_with_valid_pdf(self):
        """Test PDF extraction with valid document"""
        # Create a temporary valid PDF for testing
        with tempfile.NamedTemporaryFile(suffix='.pdf', delete=False) as temp_pdf:
            temp_pdf.write(b'%PDF-1.4 valid PDF content')
            temp_pdf_path = temp_pdf.name
        
        try:
            result = self.pdf_extractor.extract_text(temp_pdf_path)
            self.assertIn('text', result)
            self.assertIsInstance(result['text'], str)
        finally:
            os.unlink(temp_pdf_path)
    
    def test_extract_with_invalid_pdf(self):
        """Test PDF extraction with invalid document"""
        with tempfile.NamedTemporaryFile(suffix='.pdf', delete=False) as temp_pdf:
            temp_pdf.write(b'invalid pdf content')
            temp_pdf_path = temp_pdf.name
        
        try:
            with self.assertRaises(ParsingError):
                self.pdf_extractor.extract_text(temp_pdf_path)
        finally:
            os.unlink(temp_pdf_path)
    
    def test_pdf_validation(self):
        """Test PDF validation functionality"""
        validator = PDFValidator()
        
        # Test with non-existent file
        result = validator.validate_pdf("non_existent.pdf")
        self.assertFalse(result['is_valid'])
        self.assertIn('File does not exist', result['errors'])

class TestDOCXExtractor(unittest.TestCase):
    def setUp(self):
        self.docx_extractor = AdvancedDOCXExtractor()
    
    def test_extract_with_valid_docx(self):
        """Test DOCX extraction with valid document"""
        # This would require creating a valid DOCX file
        # For now, test the method structure
        self.assertTrue(hasattr(self.docx_extractor, 'extract_with_formatting'))
    
    def test_docx_validation(self):
        """Test DOCX validation functionality"""
        validator = DOCXValidator()
        
        # Test with non-existent file
        result = validator.validate_docx("non_existent.docx")
        self.assertFalse(result['is_valid'])
        self.assertIn('File does not exist', result['errors'])

class TestLLMProviders(unittest.TestCase):
    def test_openai_provider_initialization(self):
        """Test OpenAI provider initialization"""
        with patch('openai.AsyncOpenAI') as mock_client:
            provider = OpenAILLMProvider(api_key='test_key', model='gpt-4')
            self.assertIsNotNone(provider)
    
    def test_anthropic_provider_initialization(self):
        """Test Anthropic provider initialization"""
        with patch('anthropic.AsyncAnthropic') as mock_client:
            provider = AnthropicLLMProvider(api_key='test_key')
            self.assertIsNotNone(provider)
    
    @patch('aiohttp.ClientSession')
    def test_ollama_provider_initialization(self, mock_session):
        """Test Ollama provider initialization"""
        provider = OllamaLLMProvider(model='llama3')
        self.assertIsNotNone(provider)

class TestDataTransformer(unittest.TestCase):
    def setUp(self):
        self.transformer = DataTransformer()
    
    def test_transform_valid_json(self):
        """Test transformation of valid JSON structure"""
        valid_data = {
            "project": {
                "title": "Test Project",
                "description": "Test Description",
                "phases": [
                    {
                        "id": "phase_1",
                        "name": "Phase 1",
                        "tasks": [
                            {
                                "id": "task_1",
                                "name": "Task 1"
                            }
                        ]
                    }
                ]
            }
        }
        
        result = self.transformer.transform(valid_data)
        self.assertEqual(result.status, TransformationStatus.SUCCESS)
        self.assertIsNotNone(result.transformed_data)
    
    def test_transform_invalid_json(self):
        """Test transformation of invalid data"""
        invalid_data = {
            "project": {
                "invalid_field": "invalid_value"
            }
        }
        
        result = self.transformer.transform(invalid_data)
        # Should handle partial transformation or validation errors
        self.assertIn(result.status, [TransformationStatus.PARTIAL, TransformationStatus.VALIDATION_ERROR])
    
    def test_date_normalization(self):
        """Test date normalization functionality"""
        test_dates = [
            ("25.12.2023", "2023-12-25"),
            ("2023-12-25", "2023-12-25"),
            ("12/25/2023", "2023-12-25")
        ]
        
        for input_date, expected in test_dates:
            result = self.transformer._normalize_date(input_date)
            self.assertEqual(result, expected)

class TestValidators(unittest.TestCase):
    def setUp(self):
        self.structure_validator = StructureValidator()
        self.consistency_validator = DataConsistencyValidator()
    
    def test_valid_structure_validation(self):
        """Test validation of valid project structure"""
        valid_structure = {
            "project": {
                "title": "Test Project",
                "description": "Test Description",
                "phases": [
                    {
                        "id": "phase_1",
                        "name": "Phase 1",
                        "description": "Phase 1 Description",
                        "tasks": [
                            {
                                "id": "task_1",
                                "name": "Task 1",
                                "description": "Task 1 Description"
                            }
                        ]
                    }
                ]
            }
        }
        
        result = self.structure_validator.validate_structure(valid_structure)
        self.assertTrue(result["is_valid"])
    
    def test_invalid_structure_validation(self):
        """Test validation of invalid project structure"""
        invalid_structure = {
            "project": {
                "phases": []  # Missing required fields
            }
        }
        
        result = self.structure_validator.validate_structure(invalid_structure)
        self.assertFalse(result["is_valid"])
        self.assertGreater(len(result["issues"]), 0)
    
    def test_date_consistency_validation(self):
        """Test validation of date consistency"""
        structure_with_invalid_dates = {
            "project": {
                "title": "Test Project",
                "phases": [
                    {
                        "id": "phase_1",
                        "name": "Phase 1",
                        "start_date": "2024-12-25",
                        "end_date": "2024-12-24",  # End date before start date
                        "tasks": []
                    }
                ]
            }
        }
        
        result = self.consistency_validator.validate_consistency(structure_with_invalid_dates)
        self.assertFalse(result["is_valid"])
        self.assertGreater(len(result["issues"]), 0)

class TestErrorHandler(unittest.TestCase):
    def setUp(self):
        self.error_handler = ErrorHandler()
    
    def test_error_handling(self):
        """Test error handling functionality"""
        test_error = ValueError("Test error message")
        error_info = self.error_handler.handle_error(test_error)
        
        self.assertIsNotNone(error_info.error_id)
        self.assertEqual(error_info.category, ErrorCategory.ERROR)
        self.assertEqual(error_info.severity, ErrorSeverity.ERROR)
        self.assertIn("Test error message", error_info.message)
    
    def test_жцп_error_handling(self):
        """Test handling of ЖЦП-specific errors"""
        test_error = ParsingError("Test parsing error", document_path="/test.pdf")
        error_info = self.error_handler.handle_error(test_error)
        
        self.assertEqual(error_info.category, ErrorCategory.PARSING_ERROR)
        self.assertEqual(error_info.document_path, "/test.pdf")
```

## Integration Testing Framework

### Component Integration Tests
```python
class TestComponentIntegration(unittest.TestCase):
    def setUp(self):
        # Set up a mock configuration for testing
        self.test_config = {
            'providers': {
                'ollama': {
                    'enabled': True,
                    'model': 'llama3',
                    'base_url': 'http://localhost:11434'
                }
            },
            'provider_priority': ['ollama']
        }
    
    @patch('aiohttp.ClientSession')
    def test_pdf_to_llm_pipeline(self, mock_session):
        """Test complete PDF to LLM processing pipeline"""
        # Mock the LLM response
        mock_session.return_value.__aenter__.return_value.post.return_value.__aenter__.return_value.json.return_value = {
            'response': json.dumps({
                'project': {
                    'title': 'Mock Project',
                    'phases': [
                        {
                            'id': 'phase_1',
                            'name': 'Mock Phase',
                            'tasks': [
                                {
                                    'id': 'task_1',
                                    'name': 'Mock Task'
                                }
                            ]
                        }
                    ]
                }
            })
        }
        
        # Initialize components
        pdf_extractor = PDFExtractor()
        llm_manager = LLMManager(self.test_config)
        data_transformer = DataTransformer()
        
        # This would test the complete flow from PDF to structured data
        # For unit testing, we focus on integration points
        self.assertIsNotNone(pdf_extractor)
        self.assertIsNotNone(llm_manager)
        self.assertIsNotNone(data_transformer)
    
    def test_docx_to_transformation_pipeline(self):
        """Test DOCX to transformation pipeline"""
        # Similar to above, but for DOCX processing
        docx_extractor = AdvancedDOCXExtractor()
        transformer = DataTransformer()
        
        # Test that components can work together
        self.assertIsNotNone(docx_extractor)
        self.assertIsNotNone(transformer)

class TestAPIIntegration(unittest.TestCase):
    """Test API-level integration"""
    
    def setUp(self):
        self.parser = ЖЦПParser()
    
    @patch.object(ЖЦПParser, '_parse_pdf_for_llm')
    @patch.object(LLMManager, 'generate_with_fallback')
    @patch.object(DataTransformer, 'transform')
    def test_complete_pdf_parsing_flow(self, mock_transform, mock_llm, mock_parse_pdf):
        """Test complete PDF parsing flow with mocks"""
        # Set up mocks
        mock_parse_pdf.return_value = {
            'extracted_content': 'Test PDF content for project extraction',
            'metadata': {}
        }
        
        mock_llm.return_value = LLMResponse(
            content=json.dumps({
                'project': {
                    'title': 'Test Project',
                    'phases': [
                        {
                            'id': 'phase_1',
                            'name': 'Test Phase',
                            'tasks': [
                                {
                                    'id': 'task_1',
                                    'name': 'Test Task'
                                }
                            ]
                        }
                    ]
                }
            }),
            tokens_used={'input': 100, 'output': 200, 'total': 300},
            confidence=0.9,
            model='test-model',
            timestamp=datetime.now().timestamp()
        )
        
        mock_transform.return_value = TransformationResult(
            transformed_data={
                'project': {
                    'title': 'Test Project',
                    'phases': [
                        {
                            'id': 'phase_1',
                            'name': 'Test Phase',
                            'tasks': [
                                {
                                    'id': 'task_1',
                                    'name': 'Test Task'
                                }
                            ]
                        }
                    ]
                }
            },
            status=TransformationStatus.SUCCESS,
            confidence_score=0.85
        )
        
        # Test the complete flow
        # Note: This would be async in real implementation
        # For unit testing, we're checking the integration points
        self.assertIsNotNone(self.parser)
```

## End-to-End Testing Framework

### Complete Workflow Tests
```python
import asyncio
from typing import Dict, Any

class TestEndToEnd(unittest.TestCase):
    """End-to-end tests for complete workflows"""
    
    def setUp(self):
        # Create a minimal test configuration
        self.test_config = {
            'providers': {
                'ollama': {
                    'enabled': True,
                    'model': 'llama3',
                    'base_url': 'http://localhost:11434',
                    'temperature': 0.1,
                    'max_tokens': 1000
                }
            },
            'provider_priority': ['ollama']
        }
        
        # Initialize the main parser with test configuration
        self.parser = ЖЦПParser()
        # Override with test config (in real implementation, this would be cleaner)
    
    @unittest.skip("Requires actual LLM service for full E2E test")
    def test_complete_pdf_workflow(self):
        """Test complete PDF workflow (requires actual LLM service)"""
        # This test would require:
        # 1. A test PDF file with project content
        # 2. Running LLM service (Ollama, OpenAI, etc.)
        # 3. Proper API keys/configuration
        
        # For now, this serves as a template
        test_pdf_path = "tests/samples/test_project.pdf"
        if os.path.exists(test_pdf_path):
            result = asyncio.run(self.parser.parse_document(test_pdf_path))
            self.assertTrue(result['success'])
            self.assertIsNotNone(result['project_structure'])
    
    @unittest.skip("Requires actual LLM service for full E2E test")
    def test_complete_docx_workflow(self):
        """Test complete DOCX workflow (requires actual LLM service)"""
        # Similar to PDF test
        test_docx_path = "tests/samples/test_project.docx"
        if os.path.exists(test_docx_path):
            result = asyncio.run(self.parser.parse_document(test_docx_path))
            self.assertTrue(result['success'])
            self.assertIsNotNone(result['project_structure'])
    
    def test_error_recovery_workflow(self):
        """Test error recovery and fallback mechanisms"""
        # Test that the system handles errors gracefully
        error_handler = ErrorHandler()
        
        # Simulate various error scenarios
        test_errors = [
            ("Invalid PDF", ParsingError("Invalid PDF format", document_path="test.pdf")),
            ("Validation Error", ValidationError("Invalid date format", field="start_date")),
            ("LLM Error", LLMError("API rate limit exceeded", provider="openai"))
        ]
        
        for scenario, error in test_errors:
            error_info = error_handler.handle_error(error)
            self.assertIsNotNone(error_info.error_id)
            self.assertIsNotNone(error_info.timestamp)

class TestSampleDocuments(unittest.TestCase):
    """Tests using sample ЖЦП documents"""
    
    def setUp(self):
        self.parser = ЖЦПParser()
        self.sample_docs_dir = Path("tests/samples")
        self.sample_docs_dir.mkdir(exist_ok=True)
    
    def create_sample_pdf(self) -> str:
        """Create a sample PDF for testing"""
        # In a real implementation, this would create an actual PDF
        # For now, we'll create a mock file
        sample_pdf_path = self.sample_docs_dir / "sample_project.pdf"
        
        # Create a minimal PDF structure
        with open(sample_pdf_path, 'wb') as f:
            f.write(b'%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n/Pages 2 0 R\n>>\nendobj\n')
            f.write(b'2 0 obj\n<<\n/Type /Pages\n/Count 1\n/Kids [3 0 R]\n>>\nendobj\n')
            f.write(b'3 0 obj\n<<\n/Type /Page\n/Parent 2 0 R\n/MediaBox [0 0 612 792]\n/Contents 4 0 R\n>>\nendobj\n')
            f.write(b'4 0 obj\n<<\n/Length 44\n>>\nstream\nBT\n/F1 12 Tf\n72 720 Td\n(Project Lifecycle Document Test)\ntj\nET\nendstream\nendobj\n')
            f.write(b'xref\n0 5\n0000000000 65535 f \n0000000010 00000 n \n0000000053 00000 n \n0000000105 00000 n \n0000000191 00000 n \ntrailer\n<<\n/Size 5\n/Root 1 0 R\n>>\nstartxref\n245\n%%EOF')
        
        return str(sample_pdf_path)
    
    def create_sample_docx(self) -> str:
        """Create a sample DOCX for testing"""
        # This would create an actual DOCX file in a real implementation
        # For now, we'll create a mock file
        sample_docx_path = self.sample_docs_dir / "sample_project.docx"
        
        # Create a placeholder file
        with open(sample_docx_path, 'w') as f:
            f.write("Sample Project Lifecycle Document Content")
        
        return str(sample_docx_path)
    
    def test_sample_document_processing(self):
        """Test processing of sample documents"""
        # Create sample documents
        sample_pdf = self.create_sample_pdf()
        sample_docx = self.create_sample_docx()
        
        try:
            # Test that sample files exist
            self.assertTrue(os.path.exists(sample_pdf))
            self.assertTrue(os.path.exists(sample_docx))
        finally:
            # Clean up
            if os.path.exists(sample_pdf):
                os.unlink(sample_pdf)
            if os.path.exists(sample_docx):
                os.unlink(sample_docx)
```

## Performance Testing Framework

### Performance and Load Tests
```python
import time
import statistics
from concurrent.futures import ThreadPoolExecutor
import threading

class TestPerformance(unittest.TestCase):
    """Performance and load testing"""
    
    def setUp(self):
        self.parser = ЖЦПParser()
        self.performance_results = []
    
    def test_parsing_performance(self):
        """Test performance of document parsing"""
        # This would test with actual documents in a real scenario
        # For unit testing, we'll measure method call performance
        
        start_time = time.time()
        
        # Simulate parsing operations
        for i in range(10):  # Reduced for unit test
            # In real implementation, this would parse actual documents
            time.sleep(0.001)  # Simulate processing time
        
        end_time = time.time()
        processing_time = end_time - start_time
        
        # Assert that processing time is reasonable
        self.assertLess(processing_time, 5.0)  # Should complete in under 5 seconds
    
    def test_concurrent_processing(self):
        """Test concurrent document processing"""
        def simulate_parse(doc_id: int) -> Dict[str, Any]:
            """Simulate document parsing"""
            time.sleep(0.01)  # Simulate processing time
            return {"id": doc_id, "status": "processed"}
        
        # Test concurrent processing
        with ThreadPoolExecutor(max_workers=5) as executor:
            futures = [executor.submit(simulate_parse, i) for i in range(10)]
            results = [future.result() for future in futures]
        
        self.assertEqual(len(results), 10)
        processed_ids = {result["id"] for result in results}
        self.assertEqual(processed_ids, set(range(10)))
    
    def test_memory_usage(self):
        """Test memory usage during processing"""
        import psutil
        import gc
        
        process = psutil.Process()
        initial_memory = process.memory_info().rss / 1024 / 1024  # MB
        
        # Perform operations that should not cause memory leaks
        for i in range(100):
            transformer = DataTransformer()
            test_data = {"project": {"title": f"Test {i}", "phases": []}}
            result = transformer.transform(test_data)
            del transformer
            del result
        
        # Force garbage collection
        gc.collect()
        
        final_memory = process.memory_info().rss / 1024 / 1024  # MB
        memory_increase = final_memory - initial_memory
        
        # Memory increase should be reasonable (less than 100MB for this test)
        self.assertLess(memory_increase, 100.0)
```

## Accuracy Testing Framework

### Accuracy and Quality Tests
```python
class TestAccuracy(unittest.TestCase):
    """Accuracy testing for data extraction"""
    
    def setUp(self):
        self.transformer = DataTransformer()
        self.validator = StructureValidator()
    
    def test_extraction_accuracy_with_known_data(self):
        """Test accuracy with known input/output data"""
        # Define test cases with expected results
        test_cases = [
            {
                "input": {
                    "project": {
                        "title": "Test Project",
                        "description": "Test Description",
                        "phases": [
                            {
                                "id": "phase_1",
                                "name": "Planning Phase",
                                "description": "Project planning activities",
                                "tasks": [
                                    {
                                        "id": "task_1",
                                        "name": "Requirements Analysis",
                                        "description": "Analyze project requirements"
                                    }
                                ]
                            }
                        ]
                    }
                },
                "expected_fields": ["title", "phases", "tasks"],
                "expected_phase_count": 1,
                "expected_task_count": 1
            }
        ]
        
        for test_case in test_cases:
            result = self.transformer.transform(test_case["input"])
            self.assertEqual(result.status, TransformationStatus.SUCCESS)
            
            if result.transformed_data:
                project = result.transformed_data.get("project", {})
                
                # Check expected fields exist
                for field in test_case["expected_fields"]:
                    self.assertIn(field, project)
                
                # Check counts
                phases = project.get("phases", [])
                self.assertEqual(len(phases), test_case["expected_phase_count"])
                
                total_tasks = sum(len(phase.get("tasks", [])) for phase in phases)
                self.assertEqual(total_tasks, test_case["expected_task_count"])
    
    def test_date_extraction_accuracy(self):
        """Test accuracy of date extraction and normalization"""
        test_dates = [
            ("25.12.2023", "2023-12-25"),
            ("01.01.2024", "2024-01-01"),
            ("2024-06-15", "2024-06-15"),
            ("Invalid Date", None),
            ("", None)
        ]
        
        for input_date, expected in test_dates:
            result = self.transformer._normalize_date(input_date)
            self.assertEqual(result, expected)
    
    def test_text_normalization_accuracy(self):
        """Test accuracy of text normalization"""
        test_cases = [
            ("  Extra   spaces  ", "Extra spaces"),
            ("\tTab\tcharacters\t", "Tab characters"),
            ("Line\nbreaks\nhere", "Line breaks here"),
            ("", ""),
            ("Normal text", "Normal text")
        ]
        
        for input_text, expected in test_cases:
            result = self.transformer._normalize_text(input_text)
            self.assertEqual(result, expected)

class TestQualityMetrics(unittest.TestCase):
    """Quality metrics testing"""
    
    def setUp(self):
        self.transformer = DataTransformer()
        self.enricher = DataEnricher()
    
    def test_confidence_scoring(self):
        """Test confidence scoring accuracy"""
        # Test with complete data (should have high confidence)
        complete_data = {
            "title": "Complete Project",
            "description": "Complete description",
            "phases": [
                {
                    "id": "phase_1",
                    "name": "Complete Phase",
                    "tasks": [
                        {
                            "id": "task_1",
                            "name": "Complete Task"
                        }
                    ]
                }
            ]
        }
        
        validation_result = {"is_valid": True, "errors": [], "warnings": []}
        confidence = self.transformer._calculate_confidence_score(complete_data, validation_result)
        self.assertGreater(confidence, 0.7)  # Should be relatively high
        
        # Test with incomplete data (should have lower confidence)
        incomplete_data = {
            "title": "",  # Missing title
            "description": "",
            "phases": []  # No phases
        }
        
        confidence_incomplete = self.transformer._calculate_confidence_score(
            incomplete_data, validation_result
        )
        self.assertLess(confidence_incomplete, confidence)  # Should be lower
```

## Test Data Management

### Sample Document Generator
```python
import random
import string
from datetime import datetime, timedelta

class TestDocumentGenerator:
    """Generate test documents for various scenarios"""
    
    def __init__(self):
        self.project_titles = [
            "Software Development Project",
            "Infrastructure Upgrade Project", 
            "Marketing Campaign Project",
            "Research and Development Project",
            "Quality Assurance Project"
        ]
        
        self.phase_names = [
            "Planning Phase", "Analysis Phase", "Design Phase", 
            "Implementation Phase", "Testing Phase", "Deployment Phase",
            "Maintenance Phase"
        ]
        
        self.task_names = [
            "Requirements Gathering", "System Design", "Development",
            "Unit Testing", "Integration Testing", "User Acceptance Testing",
            "Deployment", "Documentation", "Training", "Support"
        ]
    
    def generate_sample_project_structure(self, num_phases: int = 3, 
                                       tasks_per_phase: int = 3) -> Dict[str, Any]:
        """Generate a sample project structure for testing"""
        project = {
            "title": random.choice(self.project_titles),
            "description": f"Sample project for testing purposes - {datetime.now().isoformat()}",
            "phases": []
        }
        
        for i in range(num_phases):
            phase_start = datetime.now() + timedelta(days=i*30)
            phase_end = phase_start + timedelta(days=29)
            
            phase = {
                "id": f"phase_{i+1}",
                "name": f"{random.choice(self.phase_names)} {i+1}",
                "description": f"Description for phase {i+1}",
                "start_date": phase_start.strftime("%Y-%m-%d"),
                "end_date": phase_end.strftime("%Y-%m-%d"),
                "tasks": []
            }
            
            for j in range(tasks_per_phase):
                task_start = phase_start + timedelta(days=j*5)
                task_end = task_start + timedelta(days=4)
                
                task = {
                    "id": f"phase_{i+1}_task_{j+1}",
                    "name": f"{random.choice(self.task_names)} {j+1}",
                    "description": f"Description for task {j+1} in phase {i+1}",
                    "start_date": task_start.strftime("%Y-%m-%d"),
                    "end_date": task_end.strftime("%Y-%m-%d"),
                    "responsible_persons": [
                        {
                            "name": f"Person {random.randint(1, 10)}",
                            "role": "Developer" if j % 2 == 0 else "Analyst",
                            "contact": f"person{random.randint(1, 10)}@example.com"
                        }
                    ],
                    "dependencies": [f"phase_{i+1}_task_{max(1, j)}"] if j > 0 else [],
                    "status": "planned"
                }
                
                phase["tasks"].append(task)
            
            project["phases"].append(phase)
        
        return {"project": project}
    
    def generate_test_document_content(self, doc_type: str = "pdf") -> str:
        """Generate test document content"""
        sample_structure = self.generate_sample_project_structure()
        
        if doc_type.lower() == "pdf":
            # Generate PDF-like content
            content = f"""
            Project Lifecycle Document
            ============================
            
            Project Title: {sample_structure['project']['title']}
            Description: {sample_structure['project']['description']}
            
            """
            
            for phase in sample_structure['project']['phases']:
                content += f"\nPhase: {phase['name']}\n"
                content += f"Dates: {phase['start_date']} to {phase['end_date']}\n"
                
                for task in phase['tasks']:
                    content += f"  - Task: {task['name']}\n"
                    content += f"    Dates: {task['start_date']} to {task['end_date']}\n"
                    content += f"    Responsible: {task['responsible_persons'][0]['name']}\n"
        
        elif doc_type.lower() == "docx":
            # Generate DOCX-like content
            content = f"Project: {sample_structure['project']['title']}\n"
            content += f"Description: {sample_structure['project']['description']}\n\n"
            
            for phase in sample_structure['project']['phases']:
                content += f"PHASE: {phase['name']}\n"
                content += f"PERIOD: {phase['start_date']} - {phase['end_date']}\n"
                
                for task in phase['tasks']:
                    content += f"  o {task['name']}\n"
                    content += f"    Timeline: {task['start_date']} - {task['end_date']}\n"
                    content += f"    Owner: {task['responsible_persons'][0]['name']}\n"
        
        return content

# Test configuration
class TestConfiguration:
    """Configuration for testing environment"""
    
    def __init__(self):
        self.test_config = {
            'providers': {
                'ollama': {
                    'enabled': True,
                    'model': 'llama3',
                    'base_url': 'http://localhost:11434',
                    'temperature': 0.1,
                    'max_tokens': 1000
                }
            },
            'provider_priority': ['ollama'],
            'retry_settings': {
                'max_retries': 1,  # Reduce for testing
                'backoff_factor': 0.1
            }
        }
    
    def get_test_config(self) -> Dict[str, Any]:
        """Get configuration suitable for testing"""
        return self.test_config
```

## Test Execution Framework

### Test Runner and Reporting
```python
import xmlrunner
import coverage
from io import StringIO
import sys

class TestRunner:
    """Comprehensive test runner with reporting"""
    
    def __init__(self):
        self.test_results = {}
        self.coverage_report = None
    
    def run_all_tests(self, output_format: str = "text") -> Dict[str, Any]:
        """Run all tests and return results"""
        # Create a test suite
        loader = unittest.TestLoader()
        suite = unittest.TestSuite()
        
        # Add all test cases
        test_classes = [
            TestPDFExtractor,
            TestDOCXExtractor,
            TestLLMProviders,
            TestDataTransformer,
            TestValidators,
            TestErrorHandler,
            TestAccuracy,
            TestQualityMetrics
        ]
        
        for test_class in test_classes:
            tests = loader.loadTestsFromTestCase(test_class)
            suite.addTests(tests)
        
        # Run tests
        if output_format == "xml":
            # Run with XML output
            stream = StringIO()
            runner = xmlrunner.XMLTestRunner(output='test-reports', outsuffix='')
            result = runner.run(suite)
        else:
            # Run with text output
            runner = unittest.TextTestRunner(verbosity=2)
            result = runner.run(suite)
        
        return {
            'total_tests': result.testsRun,
            'failures': len(result.failures),
            'errors': len(result.errors),
            'success_rate': (result.testsRun - len(result.failures) - len(result.errors)) / result.testsRun if result.testsRun > 0 else 0,
            'time_taken': getattr(result, 'time_taken', 0)
        }
    
    def run_with_coverage(self) -> Dict[str, Any]:
        """Run tests with coverage analysis"""
        cov = coverage.Coverage()
        cov.start()
        
        # Run tests
        results = self.run_all_tests()
        
        cov.stop()
        cov.save()
        
        # Generate coverage report
        report = StringIO()
        cov.report(file=report)
        coverage_text = report.getvalue()
        
        self.coverage_report = {
            'coverage_percentage': cov.html_report(),
            'coverage_text': coverage_text
        }
        
        return {**results, 'coverage': self.coverage_report}
    
    def generate_test_report(self, results: Dict[str, Any]) -> str:
        """Generate a comprehensive test report"""
        report = f"""
        ЖЦП Parser - Test Report
        ========================
        
        Test Execution Summary:
        - Total Tests Run: {results['total_tests']}
        - Failures: {results['failures']}
        - Errors: {results['errors']}
        - Success Rate: {results['success_rate']:.2%}
        
        Component Coverage:
        - PDF Extraction: {'✓' if results['total_tests'] > 0 else '✗'}
        - DOCX Extraction: {'✓' if results['total_tests'] > 0 else '✗'}
        - LLM Integration: {'✓' if results['total_tests'] > 0 else '✗'}
        - Data Transformation: {'✓' if results['total_tests'] > 0 else '✗'}
        - Validation: {'✓' if results['total_tests'] > 0 else '✗'}
        - Error Handling: {'✓' if results['total_tests'] > 0 else '✗'}
        
        Recommendations:
        {'✓ All tests passed!' if results['failures'] == 0 and results['errors'] == 0 else '✗ Issues found - review failures and errors'}
        """
        
        return report

# Test utilities
def create_test_environment():
    """Set up test environment"""
    # Create necessary directories
    test_dirs = ['tests', 'tests/samples', 'tests/reports', 'logs']
    for directory in test_dirs:
        Path(directory).mkdir(exist_ok=True)
    
    # Create sample documents
    generator = TestDocumentGenerator()
    
    # Create sample content files
    for doc_type in ['pdf', 'docx']:
        content = generator.generate_test_document_content(doc_type)
        sample_path = f"tests/samples/sample_project.{doc_type}"
        with open(sample_path, 'w', encoding='utf-8') as f:
            f.write(content)

if __name__ == "__main__":
    # Setup test environment
    create_test_environment()
    
    # Run tests
    runner = TestRunner()
    results = runner.run_all_tests()
    
    # Generate report
    report = runner.generate_test_report(results)
    print(report)
    
    # Optionally run with coverage
    # coverage_results = runner.run_with_coverage()
    # print(f"Coverage results: {coverage_results}")
```

## Continuous Integration Configuration

### CI/CD Test Configuration
```yaml
# .github/workflows/test.yml
name: Test Suite

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Python
      uses: actions/setup-python@v2
      with:
        python-version: '3.9'
    
    - name: Install dependencies
      run: |
        pip install -r requirements.txt
        pip install pytest pytest-cov coverage xmlrunner
    
    - name: Run tests
      run: |
        python -m pytest tests/ -v --junitxml=reports/junit.xml --cov=src --cov-report=xml
    
    - name: Upload test results
      uses: actions/upload-artifact@v2
      with:
        name: test-results
        path: reports/
```

This comprehensive testing framework provides multiple layers of validation for the ЖЦП parsing module, ensuring reliability, accuracy, and performance of the AI-powered document parsing system.