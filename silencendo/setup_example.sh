#!/bin/bash

echo "==========================================="
echo "Silencendo Chatbot Setup Guide"
echo "==========================================="

echo
echo "Current Status:"
if [ -z "$DEEPSEEK_API_KEY" ]; then
    echo "❌ DEEPSEEK_API_KEY is not set"
    if [ -z "$OPENAI_API_KEY" ]; then
        echo "❌ OPENAI_API_KEY is not set - using mock LLM"
    else
        echo "✅ OPENAI_API_KEY is set - using OpenAI LLM as fallback"
    fi
else
    echo "✅ DEEPSEEK_API_KEY is set - using DeepSeek LLM"
fi

echo
echo "To use the chatbot with a real LLM (recommended):"
echo "Option 1 - DeepSeek (primary):"
echo "1. Get a DeepSeek API key from https://www.deepseek.com/"
echo "2. Set the environment variable:"
echo "   export DEEPSEEK_API_KEY='your-api-key-here'"
echo "3. Optionally set other variables:"
echo "   export DEEPSEEK_MODEL='deepseek-chat'  # or 'deepseek-reasoner'"
echo "   export DEEPSEEK_BASE_URL='https://api.deepseek.com'"
echo
echo "Option 2 - OpenAI (fallback):"
echo "1. Get an OpenAI API key from https://platform.openai.com/"
echo "2. Set the environment variable:"
echo "   export OPENAI_API_KEY='your-api-key-here'"
echo "3. Optionally set other variables:"
echo "   export OPENAI_MODEL='gpt-4'  # or 'gpt-3.5-turbo'"
echo
echo "Alternatively, create a .env file in the project root:"
echo "   DEEPSEEK_API_KEY=your-api-key-here"
echo "   DEEPSEEK_MODEL=deepseek-chat"
echo "   OPENAI_API_KEY=your-openai-api-key-here"
echo "   OPENAI_MODEL=gpt-4"

echo
echo "To run the chatbot:"
echo "   cd silencendo"
echo "   go run main.go"

echo
echo "Basic workflow after setup:"
echo "1. Add sources: /source add file ./document.txt"
echo "2. Activate sources: /source use all"
echo "3. Process sources: /ingest"
echo "4. Ask questions: What is this document about?"
echo "==========================================="