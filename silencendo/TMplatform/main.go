package main

import (
	"fmt"
	"log"
	"os"
	"silencendo/cli"
)

func main() {
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	// Create and start the CLI interface
	cliInterface := cli.NewCLIInterface()
	if err := cliInterface.Start(); err != nil {
		fmt.Printf("Error starting CLI: %v\n", err)
		os.Exit(1)
	}
}

func loadEnv() error {
	// Try to load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		// Import godotenv only if .env file exists
		// For now, we'll skip this to avoid dependency issues
		// In a real implementation, you would use github.com/joho/godotenv
		return nil
	}
	return nil
}