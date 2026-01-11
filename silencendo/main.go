package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	fmt.Println("Silencendo document handler is now library-only. Use the TMplatform unified CLI runner to interact.")
	os.Exit(0)
}

func loadEnv() error {
	// Look for .env file in the current working directory
	envFile := ".env"

	// Read the .env file if it exists
	if _, err := os.Stat(envFile); err == nil {
		// Read and parse the .env file
		data, err := os.ReadFile(envFile)
		if err != nil {
			return err
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Remove surrounding quotes if present
				if len(value) >= 2 {
					if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
						(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
						value = value[1 : len(value)-1]
					}
				}

				// Set the environment variable if not already set
				if os.Getenv(key) == "" {
					os.Setenv(key, value)
				}
			}
		}
	}

	return nil
}
