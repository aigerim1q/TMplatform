# Knowledge + Planning Bot - MVP

An MVP universal chatbot for Knowledge + Planning with RAG (Retrieval-Augmented Generation) functionality.

## DeepSeek Integration

This bot supports integration with DeepSeek API for enhanced AI capabilities. To use DeepSeek:

1. **Get your API key**: Sign up at DeepSeek to obtain your API key
2. **Set the environment variable**:
   ```bash
   export DEEPSEEK_API_KEY=your_api_key_here
   ```
   On Windows:
   ```cmd
   setx DEEPSEEK_API_KEY "your_api_key_here"
   ```
3. **Optional configurations**:
   - `DEEPSEEK_BASE_URL`: Custom base URL (default: https://api.deepseek.com)
   - `DEEPSEEK_MODEL`: Model to use (default: deepseek-chat, options: deepseek-chat, deepseek-reasoner)

4. **Using .env file** (for local development only):
   Create a `.env` file in the project root:
   ```
   DEEPSEEK_API_KEY=your_api_key_here
   DEEPSEEK_MODEL=deepseek-chat
   ```

⚠️ **Security Note**: Never commit API keys to version control. The `.env` file is automatically ignored by git.

## Features

- **Knowledge Management**: Add and index various sources (files and URLs)
- **Planning Context**: Track project → stage → task hierarchy with goals, constraints, decisions, and next steps
- **Question Answering**: Ask questions about your knowledge sources
- **Text Editing**: Edit text based on your knowledge sources and context
- **Grounding Modes**: Strict, hybrid, and general modes for controlling source reliance

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   npm install
   ```

## Usage

To run the smoke test to verify DeepSeek integration:
```bash
npx ts-node src/smoke/deepseek_smoke.ts
```

Start the CLI bot:
```bash
npm run dev
```

## Available Commands

### Context Management
- `/project set <name>` - Set the current project
- `/stage set <name>` - Set the current stage  
- `/task set <title>` - Set the current task
- `/context show` - Show current context
- `/context reset` - Reset current context

### Source Management
- `/source add file <path>` - Add a file source (txt, md, pdf)
- `/source add url <url>` - Add a URL source
- `/source list` - List all sources
- `/source use <id1,id2>` - Set active sources by ID
- `/source use all` - Set all sources as active
- `/source clear` - Clear active sources

### Processing
- `/ingest` - Process active sources into searchable chunks

### Modes
- `/mode strict|hybrid|general` - Set grounding mode (default: strict)

### Other
- `/edit "instruction" <<< text` - Edit text with instruction
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

4. Ask questions:
   ```
   What are the main findings in the research document?
   ```

5. Edit text based on sources:
   ```
   /edit "Summarize in 3 points" <<< This is a long text that should be summarized based on the research...
   ```

## Architecture

The bot follows a modular architecture:

- **Context**: Manages project/stage/task context and planning elements
- **Sources**: Handles knowledge source management
- **Ingestion**: Processes sources into searchable chunks
- **Retrieval**: Finds relevant information using TF-IDF similarity
- **LLM**: Abstract interface with DeepSeek implementation
- **CLI**: Interactive command-line interface

## Grounding Modes

- **Strict**: Only answers from information found in sources
- **Hybrid**: Combines source information with general knowledge
- **General**: Answers using general knowledge without source constraints

## LLM Integration

The bot uses an abstract LLM interface with DeepSeek implementation. The system automatically detects if `DEEPSEEK_API_KEY` is set and uses the DeepSeek API; otherwise, it falls back to a mock implementation.

## Data Storage

All state is stored in the `.bot/` directory:
- `context.json` - Current context state
- `sources.json` - Knowledge source metadata
- `chunks.json` - Indexed text chunks

## Development

To run in development mode:
```bash
npm run dev
```

To build:
```bash
npm run build
```

## Dependencies

- Node.js 18+
- TypeScript
- pdf-parse - PDF text extraction
- cheerio - HTML parsing for web content
- node-fetch - HTTP requests
- openai - OpenAI API client (compatible with DeepSeek)
- dotenv - Environment variable management
- uuid - Unique identifier generation

## Troubleshooting Common Issues

### Chatbot Not Answering Questions Correctly

If the chatbot is returning generic responses like "This document is about: Artificial Intelligence",
it's likely using the mock LLM instead of a real one. This happens when no `DEEPSEEK_API_KEY` is configured.

To resolve this:
1. Configure a real LLM (recommended): Set the `DEEPSEEK_API_KEY` environment variable
2. Or understand that the mock LLM provides sample responses for demonstration purposes

The mock LLM has been improved to be more transparent about its limitations and provides tips
on how to configure a real LLM for better results.

For setup instructions, run: `./setup_example.sh`
