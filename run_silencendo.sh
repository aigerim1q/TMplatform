#!/bin/bash

# Master script to run the Silencendo Go application
# This script will guide you through the setup and execution process

echo "==========================================="
echo "    Silencendo - Knowledge + Planning Bot  "
echo "==========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed on your system."
    echo ""
    echo "Please install Go first:"
    echo "  Option 1: Using Homebrew (recommended)"
    echo "    brew install go"
    echo ""
    echo "  Option 2: Download from official site"
    echo "    Visit https://golang.org/dl/ and follow installation instructions"
    echo ""
    echo "After installing Go, please restart your terminal and run this script again."
    exit 1
fi

echo "✅ Go is installed!"

echo ""
echo "🚀 Starting unified TMplatform bot..."
cd silencendo/TMplatform
./run.sh