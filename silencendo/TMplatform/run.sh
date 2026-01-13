#!/bin/bash

# Script to run the Silencendo Go application with automatic module initialization

echo "Initializing Go modules and running Silencendo..."

# Navigate to the TMplatform directory
cd "$(dirname "$0")"

# Initialize and tidy Go modules
echo "Setting up Go modules..."
go mod tidy

if [ $? -eq 0 ]; then
    echo "Go modules initialized successfully!"
    echo "Running the application..."
    go run main.go
else
    echo "Error: Failed to initialize Go modules."
    echo "Please make sure Go is installed on your system."
    echo "You can download Go from https://golang.org/dl/"
fi