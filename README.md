# TMplatform

## Project Structure
- `parsing AI/` - Existing parsing AI component
- `chatbot/` - Advanced knowledge planning chatbot with RAG capabilities

## Chatbot Component
The chatbot is a knowledge planning bot with RAG (Retrieval-Augmented Generation) functionality that supports:
- File and URL source integration
- Context management (project/stage/task hierarchy)
- Question answering based on knowledge sources
- Text editing capabilities
- Multiple grounding modes (strict/hybrid/general)

### Running the Chatbot
1. Navigate to the chatbot directory: `cd chatbot`
2. Install dependencies: `npm install`
3. Set up environment variables (if using DeepSeek API)
4. Run: `npm run dev`

For detailed usage instructions, see the README.md in the chatbot directory.
