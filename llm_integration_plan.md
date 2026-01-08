# LLM API Integration Plan

## Overview
This document outlines the implementation plan for integrating Large Language Model (LLM) APIs to extract structured project information from parsed document text. The integration will support multiple providers to ensure flexibility and reliability.

## Requirements

### Functional Requirements
- Support for multiple LLM providers (OpenAI, Anthropic, Ollama/local models)
- Proper API key management and security
- Configurable prompt templates for different document types
- Rate limiting and retry mechanisms
- Response parsing and validation
- Confidence scoring for extracted data

### Technical Requirements
- Python 3.8+ compatibility
- Async support for better performance
- Proper error handling and logging
- Configurable timeout settings
- Cost tracking and optimization

## Supported LLM Providers

### 1. OpenAI GPT Models
- **Models**: gpt-4, gpt-4-turbo, gpt-3.5-turbo
- **Features**: High accuracy, good for complex document parsing
- **Considerations**: Cost per token, API rate limits

### 2. Anthropic Claude
- **Models**: claude-3-opus, claude-3-sonnet, claude-3-haiku
- **Features**: Strong reasoning capabilities, good for structured extraction
- **Considerations**: Cost, availability in region

### 3. Local Models (Ollama, vLLM)
- **Models**: Llama 3, Mistral, Mixtral
- **Features**: Cost-effective, privacy control
- **Considerations**: Hardware requirements, setup complexity

## Implementation Architecture

### Core LLM Interface
```python
from abc import ABC, abstractmethod
from typing import Dict, Any, Optional, List
import asyncio
import logging

class LLMProvider(ABC):
    """Abstract base class for LLM providers"""
    
    @abstractmethod
    async def generate(self, prompt: str, **kwargs) -> Dict[str, Any]:
        """Generate response from LLM"""
        pass
    
    @abstractmethod
    def get_cost_estimate(self, input_tokens: int, output_tokens: int) -> float:
        """Estimate cost for API call"""
        pass

class LLMResponse:
    """Structured response from LLM"""
    def __init__(self, content: str, tokens_used: Dict[str, int], 
                 confidence: float, model: str, timestamp: str):
        self.content = content
        self.tokens_used = tokens_used
        self.confidence = confidence
        self.model = model
        self.timestamp = timestamp
        self.parsed_data = None  # Will be set after JSON parsing
```

### OpenAI Implementation
```python
import openai
import tiktoken
from typing import Dict, Any, Optional
import json
import time

class OpenAILLMProvider(LLMProvider):
    def __init__(self, api_key: str, model: str = "gpt-4-turbo"):
        self.client = openai.AsyncOpenAI(api_key=api_key)
        self.model = model
        self.encoding = tiktoken.encoding_for_model(model)
        self.logger = logging.getLogger(__name__)
    
    async def generate(self, prompt: str, **kwargs) -> LLMResponse:
        """Generate response from OpenAI API"""
        try:
            start_time = time.time()
            
            response = await self.client.chat.completions.create(
                model=self.model,
                messages=[
                    {"role": "system", "content": "You are an expert in extracting structured project information from documents. Return only valid JSON without additional text."},
                    {"role": "user", "content": prompt}
                ],
                temperature=kwargs.get('temperature', 0.1),
                max_tokens=kwargs.get('max_tokens', 4096),
                response_format={"type": "json_object"}
            )
            
            content = response.choices[0].message.content
            usage = response.usage
            
            # Calculate confidence based on response quality
            confidence = self._calculate_confidence(content, usage)
            
            llm_response = LLMResponse(
                content=content,
                tokens_used={
                    'input': usage.prompt_tokens,
                    'output': usage.completion_tokens,
                    'total': usage.total_tokens
                },
                confidence=confidence,
                model=self.model,
                timestamp=time.time()
            )
            
            return llm_response
            
        except Exception as e:
            self.logger.error(f"OpenAI API error: {str(e)}")
            raise
    
    def get_cost_estimate(self, input_tokens: int, output_tokens: int) -> float:
        """Calculate cost based on OpenAI pricing"""
        # Example pricing (gpt-4-turbo): $10/1M input tokens, $30/1M output tokens
        input_cost = (input_tokens / 1_000_000) * 10
        output_cost = (output_tokens / 1_000_000) * 30
        return input_cost + output_cost
    
    def _calculate_confidence(self, content: str, usage) -> float:
        """Calculate confidence score based on response quality"""
        if not content or 'error' in content.lower():
            return 0.1
        
        # Check if response is valid JSON
        try:
            json.loads(content)
            json_validity = 1.0
        except:
            json_validity = 0.3
        
        # Consider token usage for complexity
        total_tokens = usage.total_tokens if usage else 0
        length_factor = min(total_tokens / 1000, 1.0)  # Cap at 1.0
        
        return (json_validity + length_factor) / 2
```

### Anthropic Implementation
```python
import anthropic
import json
import time

class AnthropicLLMProvider(LLMProvider):
    def __init__(self, api_key: str, model: str = "claude-3-sonnet-20240229"):
        self.client = anthropic.AsyncAnthropic(api_key=api_key)
        self.model = model
        self.logger = logging.getLogger(__name__)
    
    async def generate(self, prompt: str, **kwargs) -> LLMResponse:
        """Generate response from Anthropic API"""
        try:
            start_time = time.time()
            
            response = await self.client.messages.create(
                model=self.model,
                max_tokens=kwargs.get('max_tokens', 4096),
                temperature=kwargs.get('temperature', 0.1),
                system="You are an expert in extracting structured project information from documents. Return only valid JSON without additional text.",
                messages=[
                    {"role": "user", "content": prompt}
                ]
            )
            
            content = response.content[0].text if response.content else ""
            usage = response.usage
            
            confidence = self._calculate_confidence(content, usage)
            
            llm_response = LLMResponse(
                content=content,
                tokens_used={
                    'input': usage.input_tokens,
                    'output': usage.output_tokens,
                    'total': usage.input_tokens + usage.output_tokens
                },
                confidence=confidence,
                model=self.model,
                timestamp=time.time()
            )
            
            return llm_response
            
        except Exception as e:
            self.logger.error(f"Anthropic API error: {str(e)}")
            raise
    
    def get_cost_estimate(self, input_tokens: int, output_tokens: int) -> float:
        """Calculate cost based on Anthropic pricing"""
        # Example pricing (Claude 3 Sonnet): $3/1M input tokens, $15/1M output tokens
        input_cost = (input_tokens / 1_000_000) * 3
        output_cost = (output_tokens / 1_000_000) * 15
        return input_cost + output_cost
    
    def _calculate_confidence(self, content: str, usage) -> float:
        """Calculate confidence score for Anthropic response"""
        if not content or 'error' in content.lower():
            return 0.1
        
        try:
            json.loads(content)
            json_validity = 1.0
        except:
            json_validity = 0.3
        
        return json_validity
```

### Local Model Implementation (Ollama)
```python
import aiohttp
import json
import time

class OllamaLLMProvider(LLMProvider):
    def __init__(self, model: str = "llama3", base_url: str = "http://localhost:11434"):
        self.model = model
        self.base_url = base_url.rstrip('/')
        self.logger = logging.getLogger(__name__)
    
    async def generate(self, prompt: str, **kwargs) -> LLMResponse:
        """Generate response from local Ollama instance"""
        try:
            start_time = time.time()
            
            async with aiohttp.ClientSession() as session:
                payload = {
                    "model": self.model,
                    "prompt": prompt,
                    "stream": False,
                    "options": {
                        "temperature": kwargs.get('temperature', 0.1),
                        "num_predict": kwargs.get('max_tokens', 4096)
                    }
                }
                
                async with session.post(
                    f"{self.base_url}/api/generate",
                    json=payload
                ) as response:
                    result = await response.json()
                    
                    content = result.get('response', '')
                    tokens_used = {
                        'input': len(prompt.split()),
                        'output': len(content.split()),
                        'total': len(prompt.split()) + len(content.split())
                    }
                    
                    confidence = self._calculate_confidence(content)
                    
                    llm_response = LLMResponse(
                        content=content,
                        tokens_used=tokens_used,
                        confidence=confidence,
                        model=self.model,
                        timestamp=time.time()
                    )
                    
                    return llm_response
        
        except Exception as e:
            self.logger.error(f"Ollama API error: {str(e)}")
            raise
    
    def get_cost_estimate(self, input_tokens: int, output_tokens: int) -> float:
        """Local models have zero cost"""
        return 0.0
    
    def _calculate_confidence(self, content: str) -> float:
        """Calculate confidence for local model response"""
        if not content or 'error' in content.lower():
            return 0.1
        
        try:
            json.loads(content)
            json_validity = 1.0
        except:
            json_validity = 0.3
        
        return json_validity
```

## LLM Manager with Provider Selection

```python
import asyncio
import logging
from enum import Enum
from typing import Dict, Any, List

class ProviderType(Enum):
    OPENAI = "openai"
    ANTHROPIC = "anthropic"
    OLLAMA = "ollama"

class LLMManager:
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.providers = {}
        self.logger = logging.getLogger(__name__)
        self._initialize_providers()
    
    def _initialize_providers(self):
        """Initialize configured LLM providers"""
        provider_configs = self.config.get('providers', {})
        
        for provider_name, provider_config in provider_configs.items():
            if provider_config.get('enabled', False):
                try:
                    if provider_name == 'openai':
                        from .providers import OpenAILLMProvider
                        self.providers[ProviderType.OPENAI] = OpenAILLMProvider(
                            api_key=provider_config['api_key'],
                            model=provider_config.get('model', 'gpt-4-turbo')
                        )
                    elif provider_name == 'anthropic':
                        from .providers import AnthropicLLMProvider
                        self.providers[ProviderType.ANTHROPIC] = AnthropicLLMProvider(
                            api_key=provider_config['api_key'],
                            model=provider_config.get('model', 'claude-3-sonnet-20240229')
                        )
                    elif provider_name == 'ollama':
                        from .providers import OllamaLLMProvider
                        self.providers[ProviderType.OLLAMA] = OllamaLLMProvider(
                            model=provider_config.get('model', 'llama3'),
                            base_url=provider_config.get('base_url', 'http://localhost:11434')
                        )
                except Exception as e:
                    self.logger.error(f"Failed to initialize {provider_name} provider: {str(e)}")
    
    async def generate_with_fallback(self, prompt: str, **kwargs) -> LLMResponse:
        """Generate response with fallback to alternative providers"""
        provider_priority = self.config.get('provider_priority', [
            ProviderType.OPENAI.value,
            ProviderType.ANTHROPIC.value,
            ProviderType.OLLAMA.value
        ])
        
        last_error = None
        
        for provider_name in provider_priority:
            try:
                provider_type = ProviderType(provider_name)
                provider = self.providers.get(provider_type)
                
                if provider:
                    self.logger.info(f"Using {provider_type.value} provider")
                    response = await provider.generate(prompt, **kwargs)
                    return response
            except Exception as e:
                self.logger.warning(f"Provider {provider_name} failed: {str(e)}")
                last_error = e
                continue
        
        raise Exception(f"All providers failed. Last error: {str(last_error)}")
    
    async def generate_with_comparison(self, prompt: str, **kwargs) -> Dict[str, LLMResponse]:
        """Generate responses from multiple providers for comparison"""
        responses = {}
        
        for provider_type, provider in self.providers.items():
            try:
                response = await provider.generate(prompt, **kwargs)
                responses[provider_type.value] = response
            except Exception as e:
                self.logger.error(f"Provider {provider_type.value} comparison failed: {str(e)}")
        
        return responses
```

## Prompt Management System

```python
import json
from pathlib import Path
from typing import Dict, Any

class PromptManager:
    def __init__(self, prompts_dir: str = "prompts"):
        self.prompts_dir = Path(prompts_dir)
        self.prompts = {}
        self._load_prompts()
    
    def _load_prompts(self):
        """Load all prompt templates from files"""
        if self.prompts_dir.exists():
            for file_path in self.prompts_dir.glob("*.json"):
                with open(file_path, 'r', encoding='utf-8') as f:
                    prompt_data = json.load(f)
                    self.prompts[file_path.stem] = prompt_data
    
    def get_prompt(self, prompt_name: str, **kwargs) -> str:
        """Get formatted prompt with provided arguments"""
        if prompt_name not in self.prompts:
            raise ValueError(f"Prompt '{prompt_name}' not found")
        
        prompt_data = self.prompts[prompt_name]
        template = prompt_data.get('template', '')
        
        # Replace placeholders with actual values
        try:
            formatted_prompt = template.format(**kwargs)
            return formatted_prompt
        except KeyError as e:
            raise ValueError(f"Missing required argument for prompt '{prompt_name}': {str(e)}")
    
    def create_extraction_prompt(self, document_content: str, 
                               json_schema: Dict[str, Any]) -> str:
        """Create specialized prompt for project structure extraction"""
        return self.get_prompt(
            'project_extraction',
            document_content=document_content,
            json_schema=json.dumps(json_schema, ensure_ascii=False, indent=2)
        )
```

## Configuration Management

```python
import yaml
from typing import Dict, Any

class LLMConfig:
    def __init__(self, config_path: str = "config/llm_config.yaml"):
        self.config_path = config_path
        self.config = self._load_config()
    
    def _load_config(self) -> Dict[str, Any]:
        """Load LLM configuration from YAML file"""
        default_config = {
            'providers': {
                'openai': {
                    'enabled': False,
                    'api_key': '',
                    'model': 'gpt-4-turbo',
                    'temperature': 0.1,
                    'max_tokens': 4096
                },
                'anthropic': {
                    'enabled': False,
                    'api_key': '',
                    'model': 'claude-3-sonnet-20240229',
                    'temperature': 0.1,
                    'max_tokens': 4096
                },
                'ollama': {
                    'enabled': True,
                    'model': 'llama3',
                    'base_url': 'http://localhost:11434',
                    'temperature': 0.1,
                    'max_tokens': 4096
                }
            },
            'provider_priority': ['openai', 'anthropic', 'ollama'],
            'retry_settings': {
                'max_retries': 3,
                'backoff_factor': 1.0,
                'status_codes': [429, 502, 503, 504]
            },
            'rate_limiting': {
                'requests_per_minute': 60,
                'tokens_per_minute': 100000
            }
        }
        
        try:
            with open(self.config_path, 'r', encoding='utf-8') as f:
                user_config = yaml.safe_load(f)
                # Merge user config with defaults
                return self._deep_merge(default_config, user_config)
        except FileNotFoundError:
            # Create default config file
            with open(self.config_path, 'w', encoding='utf-8') as f:
                yaml.dump(default_config, f, default_flow_style=False, allow_unicode=True)
            return default_config
    
    def _deep_merge(self, base: Dict, override: Dict) -> Dict:
        """Deep merge two dictionaries"""
        result = base.copy()
        for key, value in override.items():
            if key in result and isinstance(result[key], dict) and isinstance(value, dict):
                result[key] = self._deep_merge(result[key], value)
            else:
                result[key] = value
        return result
    
    def get_provider_config(self, provider_name: str) -> Dict[str, Any]:
        """Get configuration for specific provider"""
        return self.config['providers'].get(provider_name, {})
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
        
        # Initialize LLM components
        self.llm_config = LLMConfig()
        self.llm_manager = LLMManager(self.llm_config.config)
        self.prompt_manager = PromptManager()
        
        self.logger = logging.getLogger(__name__)
    
    async def extract_project_structure(self, document_path: str) -> Dict[str, Any]:
        """Main method to extract project structure from document"""
        
        # Parse document based on type
        if document_path.lower().endswith('.pdf'):
            doc_content = await self._parse_pdf_for_llm(document_path)
        elif document_path.lower().endswith('.docx'):
            doc_content = await self._parse_docx_for_llm(document_path)
        else:
            raise ValueError(f"Unsupported document format: {document_path}")
        
        # Create extraction prompt
        json_schema = self._get_project_json_schema()
        prompt = self.prompt_manager.create_extraction_prompt(
            document_content=doc_content,
            json_schema=json_schema
        )
        
        # Generate response from LLM
        llm_response = await self.llm_manager.generate_with_fallback(
            prompt,
            temperature=0.1,
            max_tokens=4096
        )
        
        # Parse and validate response
        try:
            extracted_data = json.loads(llm_response.content)
            validated_data = self._validate_extraction(extracted_data)
            
            return {
                'project_structure': validated_data,
                'extraction_metadata': {
                    'confidence': llm_response.confidence,
                    'model_used': llm_response.model,
                    'tokens_used': llm_response.tokens_used,
                    'processing_time': llm_response.timestamp
                }
            }
        except json.JSONDecodeError as e:
            self.logger.error(f"LLM response is not valid JSON: {str(e)}")
            raise ValueError(f"LLM response parsing failed: {str(e)}")
    
    def _get_project_json_schema(self) -> Dict[str, Any]:
        """Return the expected JSON schema for project structure"""
        return {
            "type": "object",
            "properties": {
                "project": {
                    "type": "object",
                    "properties": {
                        "title": {"type": "string"},
                        "description": {"type": "string"},
                        "phases": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "id": {"type": "string"},
                                    "name": {"type": "string"},
                                    "description": {"type": "string"},
                                    "start_date": {"type": ["string", "null"], "format": "date"},
                                    "end_date": {"type": ["string", "null"], "format": "date"},
                                    "tasks": {
                                        "type": "array",
                                        "items": {
                                            "type": "object",
                                            "properties": {
                                                "id": {"type": "string"},
                                                "name": {"type": "string"},
                                                "description": {"type": "string"},
                                                "start_date": {"type": ["string", "null"], "format": "date"},
                                                "end_date": {"type": ["string", "null"], "format": "date"},
                                                "responsible_persons": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "object",
                                                        "properties": {
                                                            "name": {"type": "string"},
                                                            "role": {"type": "string"},
                                                            "contact": {"type": "string"}
                                                        }
                                                    }
                                                },
                                                "dependencies": {"type": "array", "items": {"type": "string"}},
                                                "status": {"type": "string", "enum": ["planned", "in_progress", "completed"]}
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
```

## Error Handling and Retry Mechanisms

```python
import asyncio
import time
from functools import wraps

def retry_async(max_retries: int = 3, backoff_factor: float = 1.0, 
                status_codes: List[int] = None):
    """Decorator for retrying async functions"""
    if status_codes is None:
        status_codes = [429, 502, 503, 504]
    
    def decorator(func):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            last_exception = None
            
            for attempt in range(max_retries + 1):
                try:
                    return await func(*args, **kwargs)
                except Exception as e:
                    last_exception = e
                    
                    if attempt == max_retries:
                        break
                    
                    # Calculate delay with exponential backoff
                    delay = backoff_factor * (2 ** attempt)
                    await asyncio.sleep(delay)
            
            raise last_exception
        return wrapper
    return decorator

class LLMError(Exception):
    """Custom exception for LLM-related errors"""
    pass

class RateLimitError(LLMError):
    """Exception for rate limit exceeded errors"""
    pass

class APIError(LLMError):
    """Exception for general API errors"""
    pass
```

## Testing Strategy

### Unit Tests
- Test each LLM provider implementation
- Test prompt formatting and template loading
- Test error handling and retry mechanisms
- Test cost calculation accuracy

### Integration Tests
- Test end-to-end extraction pipeline
- Test provider fallback mechanisms
- Test with various document types and complexities
- Test performance under load

## Security Considerations

### API Key Management
- Store API keys in environment variables or secure vault
- Implement key rotation mechanisms
- Use configuration validation to prevent accidental exposure

### Rate Limiting
- Implement client-side rate limiting
- Handle API rate limit responses gracefully
- Monitor and log rate limit events

This LLM integration plan provides a flexible, robust foundation for connecting the document parsing system with AI capabilities while ensuring reliability, cost optimization, and security.