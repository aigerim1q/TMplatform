# Silencendo - Knowledge + Planning Bot (Go Version)

Optimized Go-only implementation of the Knowledge + Planning Bot with RAG (Retrieval-Augmented Generation) functionality.

## Features

- **Knowledge Management**: Add and index various sources (files and URLs)
- **Planning Context**: Track project → stage → task hierarchy with goals, constraints, decisions, and next steps
- **Question Answering**: Ask questions about your knowledge sources
- **Text Editing**: Edit text based on your knowledge sources and context
- **Grounding Modes**: Strict, hybrid, and general modes for controlling source reliance

## Prerequisites

- Go 1.21+ installed on your system

## Quick Setup & Run

1. **Install Go** (if not already installed):
   - Using Homebrew: `brew install go`
   - Or download from: https://golang.org/dl/

2. **Run the application with a single command**:
   After first setup, simply use:
   ```bash
   ./run
   ```
   
   Or alternatively (legacy):
   ```bash
   run-silencendo
   ```

## Manual Setup

1. **Navigate to the project directory**:
   ```bash
   cd silencendo
   ```

2. **Initialize Go modules**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

## Configuration

The DeepSeek API key is pre-configured in the `.env` file. If you need to update it, edit the `.env` file in either the main `silencendo` or `TMplatform` directory.

## Available Commands

Once the application is running, you can use these commands:
- `/project set <name>` - Set the current project
- `/stage set <name>` - Set the current stage
- `/task set <title>` - Set the current task
- `/context show` - Show current context
- `/source add file <path>` - Add a file source
- `/source add url <url>` - Add a URL source
- `/source list` - List all sources
- `/source use <id1,id2>` - Set active sources
- `/ingest` - Process active sources into searchable chunks
- `/mode strict|hybrid|general` - Set grounding mode
- `/help` - Show help message
- `/exit` or `/quit` - Exit the bot

## Example Workflow

1. Set up your project context:
   ```
   /project set Demo
   /stage set Research
   /task set Analyze market trends
   ```

2. Add knowledge sources:
   ```
   /source add file ./documents/research.pdf
   /source add url https://example.com/article
   /source list
   /source use all
   ```

3. Process sources:
   ```
   /ingest
   ```

4. Ask questions or edit text based on your sources.