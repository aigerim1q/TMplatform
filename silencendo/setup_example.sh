#!/bin/bash

echo "==========================================="
echo "Silencendo Chatbot Setup Guide"
echo "==========================================="

echo
echo "Current Status:"
if [ -z "$DEEPSEEK_API_KEY" ]; then
    echo "❌ DEEPSEEK_API_KEY is not set - using mock LLM"
else
    echo "✅ DEEPSEEK_API_KEY is set - using DeepSeek LLM"
fi

echo
echo "To use the chatbot with a real LLM (recommended):"
echo "1. Get a DeepSeek API key from https://www.deepseek.com/"
echo "2. Set the environment variable:"
echo "   export DEEPSEEK_API_KEY='your-api-key-here'"
echo "3. Optionally set other variables:"
echo "   export DEEPSEEK_MODEL='deepseek-chat'  # or 'deepseek-reasoner'"
echo "   export DEEPSEEK_BASE_URL='https://api.deepseek.com'"
echo
echo "Alternatively, create a .env file in the project root:"
echo "   DEEPSEEK_API_KEY=your-api-key-here"
echo "   DEEPSEEK_MODEL=deepseek-chat"

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