package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"silencendo/chatbot"
	"silencendo/cli"
	"silencendo/db"
	router "silencendo/intent_router"
	"silencendo/llm"
	"silencendo/models"
	"silencendo/services"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	cfg := db.LoadConfig()
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	if err := db.Migrate(dbConn); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	actor := models.User{
		ID:   envOrDefault("BOT_USER_ID", "cli-user"),
		Name: envOrDefault("BOT_USER_NAME", "CLI User"),
	}
	if email := os.Getenv("BOT_USER_EMAIL"); strings.TrimSpace(email) != "" {
		emailCopy := strings.TrimSpace(email)
		actor.Email = &emailCopy
	}

	projectService := services.NewProjectService(dbConn)
	projectHandler := chatbot.NewHandler(projectService, actor)

	// Initialize the LLM router and set it globally
	llmClient := llm.CreateLLMClient()
	if llmClient != nil {
		llmRouter := router.New(llmClient)
		chatbot.SetGlobalRouter(llmRouter)
	}

	state := cli.NewSessionState()
	cliInterface := cli.NewCLIInterface()
	cliRouter := cli.NewIntentRouter(cliInterface.SourceManager(), cliInterface.IntentDetector())
	dispatcher := cli.NewDispatcher(cliRouter, projectHandler, cliInterface, state)
	repl := cli.NewUnifiedRepl(dispatcher)

	if err := repl.Start(context.Background(), os.Stdin); err != nil {
		fmt.Printf("Error starting CLI: %v\n", err)
		os.Exit(1)
	}
}

func loadEnv() error {
	if _, err := os.Stat(".env"); err == nil {
		return godotenv.Load()
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
