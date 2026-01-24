package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"zhcp-parser-go/internal/ai"
	"zhcp-parser-go/internal/ai/llm_providers/anthropic"
	"zhcp-parser-go/internal/ai/llm_providers/deepseek"
	"zhcp-parser-go/internal/ai/llm_providers/ollama"
	"zhcp-parser-go/internal/ai/llm_providers/openai"
	"zhcp-parser-go/internal/common"
	"zhcp-parser-go/internal/config"
	"zhcp-parser-go/internal/parser"
	"zhcp-parser-go/internal/storage/sqlite"
	"zhcp-parser-go/internal/validators"

	"github.com/spf13/cobra"
)

var (
	configPath string
	validate   bool
	enrich     bool
	outputFile string
	dbPath     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zhcp-parser",
	Short: "AI-powered parser for Project Lifecycle Documents (ЖЦП)",
	Long: `ZhCP Parser is an AI-powered module that automatically extracts project structure 
information from PDF and DOCX documents (ЖЦП - Жизненный Цикл Проекта / Project Lifecycle Documents).

Features:
- Multi-format Support: PDF and DOCX document parsing
- AI-Powered Extraction: Uses LLMs to extract structured data
- Fallback Mechanisms: Supports multiple LLM providers (OpenAI, Anthropic, Ollama)
- Data Validation: Comprehensive validation and quality assurance
- Error Handling: Robust error handling and recovery`,
}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse [document_path]",
	Short: "Parse a document and extract project structure",
	Long: `Parse a PDF or DOCX document to extract structured project information including:
- Project phases
- Tasks within each phase
- Timeline information (start/end dates)
- Responsible persons and their roles
- Task dependencies and relationships`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parseDocument(args[0])
	},
}

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch [directory_path]",
	Short: "Parse all supported documents in a directory",
	Long:  `Parse all PDF and DOCX documents in a directory and extract structured project information.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parseBatch(args[0])
	},
}

func init() {
	// Add flags to root command
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "configs/llm_config.yaml", "Configuration file path")

	// Add flags to parse command
	parseCmd.Flags().BoolVarP(&validate, "validate", "v", true, "Whether to perform validation")
	parseCmd.Flags().BoolVarP(&enrich, "enrich", "e", true, "Whether to enrich data with computed fields")
	parseCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for results (JSON format)")
	parseCmd.Flags().StringVarP(&dbPath, "db", "d", "zhcp.db", "Path to SQLite database")

	// Add subcommands
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(batchCmd)
}

func main() {
	// Register all LLM providers to avoid circular imports
	registerProviders()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// registerProviders registers all LLM providers
func registerProviders() {
	ai.RegisterProvider("openai", func(config common.ProviderConfig) (ai.LLMProvider, error) {
		return openai.NewOpenAIProvider(config.APIKey, config.Model)
	})

	ai.RegisterProvider("anthropic", func(config common.ProviderConfig) (ai.LLMProvider, error) {
		return anthropic.NewAnthropicProvider(config.APIKey, config.Model)
	})

	ai.RegisterProvider("ollama", func(config common.ProviderConfig) (ai.LLMProvider, error) {
		return ollama.NewOllamaProvider(config.Model, config.BaseURL)
	})

	ai.RegisterProvider("deepseek", func(config common.ProviderConfig) (ai.LLMProvider, error) {
		return deepseek.NewDeepSeekProvider(config.APIKey, config.Model)
	})
}

func parseDocument(documentPath string) {
	// Initialize configuration
	configManager := config.NewConfigManager(configPath)
	cfg, err := configManager.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize the parser
	zhcpParser, err := parser.NewZhcpParser(cfg)
	if err != nil {
		fmt.Printf("Error initializing parser: %v\n", err)
		os.Exit(1)
	}
	defer zhcpParser.Close()

	// Parse the document
	fmt.Printf("[DOC] Processing: %s\n", documentPath)
	fmt.Println("--------------------------------------------------")
	result, err := zhcpParser.ParseDocument(documentPath, validate, enrich)
	if err != nil {
		fmt.Printf("Error parsing document: %v\n", err)
		os.Exit(1)
	}

	// Display results
	if result.Success {
		fmt.Printf("[SUCCESS] Successfully processed: %s\n", filepath.Base(documentPath))
		fmt.Printf("[INFO] Project Title: %s\n", result.ProjectStructure.Project.Title)
		fmt.Printf("[INFO] Description: %s\n", truncateString(result.ProjectStructure.Project.Description, 100))
		fmt.Printf("[INFO] Phases: %d\n", len(result.ProjectStructure.Project.Phases))

		totalTasks := 0
		for _, phase := range result.ProjectStructure.Project.Phases {
			totalTasks += len(phase.Tasks)
		}
		fmt.Printf("[INFO] Total Tasks: %d\n", totalTasks)
		fmt.Printf("[INFO] Confidence Score: %.2f\n", result.ExtractionMetadata.Confidence)
		fmt.Println()
		fmt.Println("[INFO] Project Structure:")
		fmt.Println()

		// Print detailed project structure
		for i, phase := range result.ProjectStructure.Project.Phases {
			fmt.Printf("  Phase %d: %s\n", i+1, phase.Name)
			fmt.Printf("    Description: %s\n", truncateString(phase.Description, 100))
			if phase.StartDate != "" && phase.EndDate != "" {
				fmt.Printf("    Timeline: %s to %s\n", phase.StartDate, phase.EndDate)
			}

			for j, task := range phase.Tasks {
				fmt.Printf("      Task %d: %s\n", j+1, task.Name)
				fmt.Printf("        Description: %s\n", truncateString(task.Description, 100))
				if task.StartDate != "" && task.EndDate != "" {
					fmt.Printf("        Timeline: %s to %s\n", task.StartDate, task.EndDate)
				}

				for _, resp := range task.ResponsiblePersons {
					fmt.Printf("        Responsible: %s (%s)\n", resp.Name, resp.Role)
				}
			}
			fmt.Println()
		}

		// Database Storage Logic
		fmt.Println("--------------------------------------------------")
		fmt.Println("[DB] Validating for Database persistence...")
		dbVal := validators.NewDBValidator()
		missingFields := dbVal.ValidateForDB(result.ProjectStructure)

		if len(missingFields) > 0 {
			fmt.Println("[DB] ⚠️  Cannot save to Database. The following information is missing:")
			for _, m := range missingFields {
				fmt.Printf("  - %s\n", m)
			}
			fmt.Println("[DB] Please provide a document with more complete information.")
		} else {
			fmt.Println("[DB] Validation passed! Saving to database...")
			store := sqlite.New(dbPath)
			if err := store.Init(context.Background()); err != nil {
				fmt.Printf("[DB] ❌ Failed to initialize database: %v\n", err)
			} else {
				defer store.Close()
				dbRes, err := store.PersistProjectStructure(context.Background(), result.ProjectStructure)
				if err != nil {
					fmt.Printf("[DB] ❌ Failed to persist project: %v\n", err)
				} else {
					fmt.Printf("[DB] ✅ Successfully saved to database (%s)!\n", dbPath)
					fmt.Printf("[DB] Created Project ID: %d\n", dbRes.ProjectID)
					fmt.Printf("[DB] Created Phases: %d, Tasks: %d\n", len(dbRes.PhaseIDs), len(dbRes.TaskIDs))
				}
			}
		}
		fmt.Println("--------------------------------------------------")

		// Optionally save to output file
		if outputFile != "" {
			// Save the result to the specified file
			if err := saveResultToFile(result, outputFile); err != nil {
				fmt.Printf("Error saving results to file: %v\n", err)
			} else {
				fmt.Printf("Results saved to: %s\n", outputFile)
			}
		} else {
			// Automatically save to Desktop or Downloads with a generated filename
			outputPath, err := saveResultToDesktopOrDownloads(result, documentPath)
			if err != nil {
				fmt.Printf("Error saving results to Desktop/Downloads: %v\n", err)
			} else {
				fmt.Printf("Results automatically saved to: %s\n", outputPath)
			}
		}
		fmt.Printf("[SUCCESS] Processing completed successfully!\n")
		fmt.Println("============================================================")
	} else {
		fmt.Printf("[ERROR] Failed to process: %s\n", filepath.Base(documentPath))
		if result.Error != nil {
			fmt.Printf("Error: %s\n", result.Error.Message)
			fmt.Printf("Category: %s\n", result.Error.Category)
		}
		if result.ValidationError != nil {
			fmt.Printf("Validation errors: %v\n", result.ValidationError)
		}
	}
}

func parseBatch(directoryPath string) {
	fmt.Printf("Parsing all documents in directory: %s\n", directoryPath)

	// Read directory contents
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Filter for supported file types
	var supportedFiles []string
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".pdf" || ext == ".docx" {
				supportedFiles = append(supportedFiles, filepath.Join(directoryPath, file.Name()))
			}
		}
	}

	if len(supportedFiles) == 0 {
		fmt.Println("No supported files found in directory (PDF or DOCX)")
		return
	}

	fmt.Printf("Found %d supported files to process\n", len(supportedFiles))

	// Initialize configuration once
	configManager := config.NewConfigManager(configPath)
	cfg, err := configManager.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}

	// Initialize the parser once
	zhcpParser, err := parser.NewZhcpParser(cfg)
	if err != nil {
		fmt.Printf("Error initializing parser: %v\n", err)
		return
	}
	defer zhcpParser.Close()

	// Process each file
	successful := 0
	failed := 0

	for _, filePath := range supportedFiles {
		fmt.Printf("Processing: %s\n", filepath.Base(filePath))

		result, err := zhcpParser.ParseDocument(filePath, validate, enrich)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			failed++
			continue
		}

		if result.Success {
			fmt.Printf("  ✅ Success - Confidence: %.2f\n", result.ExtractionMetadata.Confidence)
			successful++
		} else {
			fmt.Printf("  ❌ Failed - Error: %v\n", result.Error.Message)
			failed++
		}
	}

	fmt.Printf("\nBatch processing complete!\n")
	fmt.Printf("Successful: %d\n", successful)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total: %d\n", len(supportedFiles))
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// saveResultToFile saves the parse result to a JSON file
func saveResultToFile(result *parser.ParseResult, filename string) error {
	// Convert the result to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write result to file: %w", err)
	}

	return nil
}

// getDesktopOrDownloadsPath gets the path to either Desktop or Downloads folder
func getDesktopOrDownloadsPath() (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Check for Desktop directory first
	desktopPath := filepath.Join(homeDir, "Desktop")
	if _, err := os.Stat(desktopPath); err == nil {
		return desktopPath, nil
	}

	// If Desktop doesn't exist, try Downloads
	downloadsPath := filepath.Join(homeDir, "Downloads")
	if _, err := os.Stat(downloadsPath); err == nil {
		return downloadsPath, nil
	}

	// If neither exists, return home directory
	return homeDir, nil
}

// saveResultToDesktopOrDownloads saves the result to the Desktop or Downloads folder
func saveResultToDesktopOrDownloads(result *parser.ParseResult, documentPath string) (string, error) {
	desktopOrDownloadsPath, err := getDesktopOrDownloadsPath()
	if err != nil {
		return "", fmt.Errorf("failed to get Desktop or Downloads path: %w", err)
	}

	// Generate a filename based on the input document name and timestamp
	docBaseName := filepath.Base(documentPath)
	docName := strings.TrimSuffix(docBaseName, filepath.Ext(docBaseName))
	timestamp := time.Now().Format("20060102_150405")
	outputFilename := fmt.Sprintf("%s_extraction_result_%s.json", docName, timestamp)
	outputPath := filepath.Join(desktopOrDownloadsPath, outputFilename)

	return outputPath, saveResultToFile(result, outputPath)
}
