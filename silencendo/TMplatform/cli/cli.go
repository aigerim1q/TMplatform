package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"silencendo/context"
	"silencendo/ingestion"
	"silencendo/llm"
	"silencendo/retrieval"
	"silencendo/sources"
)

type CLIInterface struct {
	ctxManager        *context.ContextManager
	ctxCommands       *context.ContextCommands
	sourceManager     *sources.SourceManager
	sourceCommands    *sources.SourceCommands
	ingestionPipeline *ingestion.IngestionPipeline
	chunkManager      *ingestion.ChunkManager
	retriever         *retrieval.Retriever
	llmClient         llm.LLMClient
	groundingMode     string // "strict", "hybrid", or "general"
}

func NewCLIInterface() *CLIInterface {
	ctxManager := context.NewContextManager()
	ctxCommands := context.NewContextCommands(ctxManager)
	
	sourceManager := sources.NewSourceManager()
	sourceCommands := sources.NewSourceCommands(sourceManager)
	
	chunkManager := ingestion.NewChunkManager()
	ingestionPipeline := ingestion.NewIngestionPipeline(nil)
	retriever := retrieval.NewRetriever(chunkManager, nil)
	llmClient := llm.CreateLLMClient()
	
	return &CLIInterface{
		ctxManager:        ctxManager,
		ctxCommands:       ctxCommands,
		sourceManager:     sourceManager,
		sourceCommands:    sourceCommands,
		ingestionPipeline: ingestionPipeline,
		chunkManager:      chunkManager,
		retriever:         retriever,
		llmClient:         llmClient,
		groundingMode:     "strict", // Default mode
	}
}

func (c *CLIInterface) Start() error {
	fmt.Println("🤖 Knowledge + Planning Bot - MVP")
	fmt.Println("Type commands starting with / or ask questions directly")
	fmt.Println("Available commands: /project, /stage, /task, /context, /source, /mode, /help")
	fmt.Println("Type /help for more information\n")

	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Print("> ")
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		
		if strings.ToLower(input) == "/exit" || strings.ToLower(input) == "/quit" {
			fmt.Println("Goodbye!")
			break
		}
		
		if strings.HasPrefix(input, "/") {
			if err := c.handleCommand(input); err != nil {
				fmt.Printf("Error handling command: %v\n", err)
			}
		} else {
			if err := c.handleInput(input); err != nil {
				fmt.Printf("Error handling input: %v\n", err)
			}
		}
		
		fmt.Print("> ")
	}
	
	return scanner.Err()
}

func (c *CLIInterface) handleCommand(input string) error {
	parts := strings.Split(input, " ")
	command := strings.ToLower(strings.TrimPrefix(parts[0], "/"))
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}
	
	switch command {
	case "help":
		c.showHelp()
	case "project":
		result := c.ctxCommands.HandleProjectCommand(args)
		fmt.Println(result)
	case "stage":
		result := c.ctxCommands.HandleStageCommand(args)
		fmt.Println(result)
	case "task":
		result := c.ctxCommands.HandleTaskCommand(args)
		fmt.Println(result)
	case "context":
		result := c.ctxCommands.HandleContextCommand(args)
		fmt.Println(result)
	case "source":
		result := c.sourceCommands.HandleSourceCommand(args)
		fmt.Println(result)
	case "mode":
		c.handleModeCommand(args)
	case "ingest":
		return c.handleIngestCommand()
	case "redact":
		return c.handleRedactCommand(args)
	default:
		fmt.Printf("Unknown command: %s. Type /help for available commands.\n", command)
	}
	
	return nil
}

func (c *CLIInterface) handleRedactCommand(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: /redact <source_id> <pattern_to_redact>")
		fmt.Println("Example: /redact 1 \"sensitive information\"")
		return nil
	}

	sourceId := args[0]
	redactPattern := strings.Join(args[1:], " ")

	// Get all chunks for this source
	chunks := c.chunkManager.GetChunksBySource(sourceId)
	if len(chunks) == 0 {
		// Try to parse as number to provide better feedback
		if num, err := strconv.Atoi(sourceId); err == nil {
			fmt.Printf("No chunks found for source Number: %d. Please ingest the source first.\n", num)
		} else {
			fmt.Printf("No chunks found for source ID: %s. Please ingest the source first.\n", sourceId)
		}
		return nil
	}

	// Perform redaction on each chunk
	redactedCount := 0
	re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(redactPattern))
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		return err
	}

	for i := range chunks {
		originalText := chunks[i].Text
		redactedText := re.ReplaceAllString(originalText, "[REDACTED]")
		if originalText != redactedText {
			chunks[i].Text = redactedText
			redactedCount++
		}
	}

	fmt.Printf("Redaction completed. %d chunks were modified.\n", redactedCount)
	fmt.Printf("Pattern \"%s\" has been replaced with [REDACTED] in source ID: %s\n", redactPattern, sourceId)
	
	return nil
}

func (c *CLIInterface) handleInput(input string) error {
	intent := c.detectIntent(input)
	
	if intent == "edit" {
		return c.handleEditRequest(input)
	} else {
		return c.handleQuestion(input)
	}
}

func (c *CLIInterface) detectIntent(input string) string {
	inputLower := strings.ToLower(input)
	
	// Weighted scoring for intent detection
	editScore := 0
	questionScore := 0
	
	// Strong indicators of edit requests with weights
	editKeywords := map[string]int{
		"edit":     2,
		"change":   2,
		"modify":   2,
		"update":   2,
		"rewrite":  2,
		"rephrase": 2,
		"fix":      2,
		"correct":  2,
		"adjust":   2,
		"alter":    2,
		"improve":  2,
		"replace":  2,
		"insert":   2,
		"remove":   2,
		"delete":   2,
		"append":   2,
		"add":      3, // Adding has high weight as it's clearly an edit operation
		"create":   3,
		"generate": 2,
	}
	
	// Strong indicators of questions with weights
	questionKeywords := map[string]int{
		"what":  2,
		"who":   2,
		"where": 2,
		"when":  2,
		"why":   2,
		"how":   2,
		"is":    1,
		"are":   1,
		"can":   1,
		"could": 1,
		"would": 1,
		"should": 1,
		"does":  1,
		"do":    1,
		"did":   1,
		"have":  1,
		"has":   1,
		"had":   1,
		"explain":    2,
		"describe":   2,
		"summarize":  2,
	}
	
	// Check for edit keywords
	for keyword, weight := range editKeywords {
		if strings.Contains(inputLower, keyword) {
			editScore += weight
		}
	}
	
	// Check for question keywords
	for keyword, weight := range questionKeywords {
		if strings.Contains(inputLower, keyword) {
			questionScore += weight
		}
	}
	
	// Edit patterns with weights
	editPatterns := []struct {
		pattern *regexp.Regexp
		weight  int
	}{
		{regexp.MustCompile(`(?i)^edit\s+`), 4},                              // "edit this"
		{regexp.MustCompile(`(?i)^change\s+`), 4},                            // "change this"
		{regexp.MustCompile(`(?i)^update\s+`), 4},                            // "update this"
		{regexp.MustCompile(`(?i)change.*to`), 3},                           // "change X to Y"
		{regexp.MustCompile(`(?i)replace.*with`), 3},                        // "replace X with Y"
		{regexp.MustCompile(`(?i)update.*to`), 3},                           // "update X to Y"
		{regexp.MustCompile(`(?i)fix.*to`), 3},                              // "fix X to Y"
		{regexp.MustCompile(`(?i)correct.*to`), 3},                          // "correct X to Y"
		{regexp.MustCompile(`(?i)make.*be`), 3},                             // "make this be X"
		{regexp.MustCompile(`(?i)set.*to`), 3},                              // "set this to X"
		{regexp.MustCompile(`(?i)add\s+.*\s+and\s+.*\s+is\s+to`), 5},       // "add John and his job is to" - high weight for this specific pattern
		{regexp.MustCompile(`(?i)add\s+.*\s+to\s+the\s+document`), 4},       // "add John to the document"
		{regexp.MustCompile(`(?i)add\s+a?\s*person\s+named`), 4},            // "add a person named"
		{regexp.MustCompile(`(?i)add\s+.*\s+with\s+task`), 4},               // "add John with task"
		{regexp.MustCompile(`(?i)add\s+.*\s+and\s+assign`), 4},              // "add John and assign"
		{regexp.MustCompile(`(?i)add\s+\w+\s+and\s+.*\s+is\s+`), 5},        // "add X and his/their Y is Z" - pattern matching the user's example
	}
	
	// Question patterns with weights
	questionPatterns := []struct {
		pattern *regexp.Regexp
		weight  int
	}{
		{regexp.MustCompile(`(?i)^(what|who|where|when|why|how)\s+`), 3},     // "what is", "who is", etc.
		{regexp.MustCompile(`.*\?$`), 3},                                 // Ends with question mark
		{regexp.MustCompile(`(?i)^(is|are|can|could|would|should)\s+`), 2},   // "is this", "can you", etc.
	}
	
	// Apply pattern weights
	for _, ep := range editPatterns {
		if ep.pattern.MatchString(input) {
			editScore += ep.weight
		}
	}
	
	for _, qp := range questionPatterns {
		if qp.pattern.MatchString(input) {
			questionScore += qp.weight
		}
	}
	
	// Additional context analysis
	if strings.HasSuffix(input, "?") {
		questionScore += 2
	}
	
	if c.containsEditIndicators(inputLower) {
		editScore += 1
	}
	
	// Special handling for the specific case from user feedback
	// "add Aidana and her task is similar to Beka aga" should be detected as edit
	if strings.Contains(inputLower, "add") && strings.Contains(inputLower, "and") && strings.Contains(inputLower, "is") {
		// This pattern often indicates adding a person/task with details
		editScore += 2
	}
	
	// Return the intent with higher score, but with some adjustments for safety
	if editScore > questionScore {
		return "edit"
	} else if questionScore > editScore {
		return "question"
	} else {
		// When scores are equal, if there's any indication of adding/modifying, lean toward edit
		// This helps with cases like "add X and Y is Z"
		if editScore > 0 {
			return "edit"
		} else {
			// Default to question for safety when truly ambiguous
			return "question"
		}
	}
}

func (c *CLIInterface) containsEditIndicators(inputLower string) bool {
	// Additional checks for edit indicators that might not be caught by keywords
	editIndicators := []string{
		"make it",
		"make this",
		"turn this",
		"convert this",
		"transform",
		"reformat",
		"restructure",
		"add to",
		"insert into",
		"put in",
		"place in",
	}
	
	for _, indicator := range editIndicators {
		if strings.Contains(inputLower, indicator) {
			return true
		}
	}
	return false
}

func (c *CLIInterface) handleEditRequest(instruction string) error {
	// Get relevant chunks based on the instruction
	var relevantChunks []ingestion.Chunk
	if c.groundingMode != "general" {
		activeSourceIds := c.sourceManager.GetActiveSourceIds()
		// Use the instruction to find relevant chunks
		retrievalResult := c.retriever.Retrieve(instruction, activeSourceIds)
		relevantChunks = retrievalResult.Chunks
	}

	ctx := c.ctxManager.GetContext()
	
	// Determine if this is an append/add operation vs a modification
	instructionLower := strings.ToLower(instruction)
	isAddOperation := strings.Contains(instructionLower, "add") ||
		strings.Contains(instructionLower, "create") ||
		strings.Contains(instructionLower, "insert")
	
	textToEdit := ""
	editInstruction := instruction
	
	if isAddOperation && len(relevantChunks) == 0 {
		// For add operations with no relevant chunks, we'll work with the original source content
		activeSources := c.sourceManager.GetActiveSources()
		if len(activeSources) > 0 {
			// Get the content of the first active source to use as base
			firstSourceId := activeSources[0].ID
			originalContent := c.sourceManager.GetSourceContent(firstSourceId)
			if originalContent != "" {
				textToEdit = originalContent
			}
		}
	} else if len(relevantChunks) > 0 {
		// Combine relevant chunks as potential text to edit
		for i, chunk := range relevantChunks {
			if i > 0 {
				textToEdit += "\n\n"
			}
			textToEdit += chunk.Text
		}
	} else {
		// If no relevant chunks and not an add operation, treat as a general edit
		textToEdit = instruction
	}
	
	// Call the LLM to perform the edit
	result, err := c.llmClient.Edit(context.Background(), textToEdit, editInstruction, relevantChunks, ctx)
	if err != nil {
		return fmt.Errorf("error during edit: %w", err)
	}
	
	fmt.Println("\n┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│                    📝 EDIT RESULT                     │")
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	
	// Split the result into lines and handle each line
	allLines := strings.Split(result, "\n")
	
	for i, line := range allLines {
		if i == 0 {
			// Handle the first line specially
			if len(line) <= 68 {
				fmt.Printf("│ %-68s │\n", line)
			} else {
				// Split first line if too long
				firstPart := line[:68]
				fmt.Printf("│ %-68s │\n", firstPart)
				
				// Handle remaining parts of the first line
				remaining := line[68:]
				for len(remaining) > 0 {
					var wrappedLine string
					if len(remaining) <= 68 {
						wrappedLine = remaining
						remaining = ""
					} else {
						wrappedLine = remaining[:68]
						remaining = remaining[68:]
					}
					fmt.Printf("│ %-68s │\n", wrappedLine)
				}
			}
		} else {
			// Handle subsequent lines
			if len(line) <= 68 {
				fmt.Printf("│ %-68s │\n", line)
			} else {
				// Wrap long lines
				for len(line) > 0 {
					var wrappedLine string
					if len(line) <= 68 {
						wrappedLine = line
						line = ""
					} else {
						wrappedLine = line[:68]
						line = line[68:]
					}
					fmt.Printf("│ %-68s │\n", wrappedLine)
				}
			}
		}
	}
	fmt.Println("└─────────────────────────────────────────────────────────┘\n")
	
	// Attempt to update the source files with the edited content
	if result != textToEdit {
		return c.attemptToUpdateSourceFiles(relevantChunks, result)
	}
	
	return nil
}

func (c *CLIInterface) attemptToUpdateSourceFiles(relevantChunks []ingestion.Chunk, editedContent string) error {
	// Get unique source IDs from the relevant chunks
	sourceIds := make(map[string]bool)
	for _, chunk := range relevantChunks {
		sourceIds[chunk.SourceID] = true
	}
	
	var uniqueSourceIds []string
	for id := range sourceIds {
		uniqueSourceIds = append(uniqueSourceIds, id)
	}
	
	// If no source IDs from chunks, try to get from active sources
	if len(uniqueSourceIds) == 0 {
		activeSources := c.sourceManager.GetActiveSources()
		if len(activeSources) > 0 {
			// Use the first active source for add operations
			firstSourceId := activeSources[0].ID
			return c.updateSingleSourceFile(firstSourceId, editedContent, relevantChunks)
		} else {
			// If no active sources, show the content for manual update
			return c.showProposedContent(editedContent)
		}
	} else {
		// Update files from relevant chunks
		for _, sourceId := range uniqueSourceIds {
			if err := c.updateSingleSourceFile(sourceId, editedContent, relevantChunks); err != nil {
				return err
			}
		}
	}
	
	return nil
}

func (c *CLIInterface) updateSingleSourceFile(sourceId string, editedContent string, relevantChunks []ingestion.Chunk) error {
	// Get the original content of the source
	originalContent := c.sourceManager.GetSourceContent(sourceId)
	if originalContent == "" {
		return nil
	}
	
	// Determine if this is likely a complete file rewrite by analyzing the edit request
	// If the edited content looks like a complete document (contains multiple lines/sections)
	// and is significantly different from individual chunks, treat it as a complete replacement
	var instruction string
	if len(relevantChunks) > 0 {
		instruction = relevantChunks[0].Text
	}
	
	isLikelyCompleteRewrite := c.isCompleteRewriteInstruction(instruction) ||
		c.isCompleteDocumentContent(editedContent, originalContent)
	
	// Check if the edit request is an add operation
	isAddOperation := strings.Contains(strings.ToLower(instruction), "add") ||
		strings.Contains(strings.ToLower(instruction), "create") ||
		strings.Contains(strings.ToLower(instruction), "insert")
	
	if isLikelyCompleteRewrite {
		// For complete rewrites (like deleting items), replace the entire file
		if c.sourceManager.UpdateSourceContent(sourceId, editedContent) {
			fmt.Println("\n📝 Source file updated successfully!")
			fmt.Println("The changes have been saved to the original document.")
			
			// After updating the file, we should re-ingest it to update the chunks
			source := c.sourceManager.GetSourceByID(sourceId)
			if source != nil {
				fmt.Println("\n💡 Tip: Run '/ingest' to update the knowledge base with the new content.")
				fmt.Printf("   Source: %d | %s | %s\n", source.Number, source.ID[:4], source.Title)
			}
		}
	} else if isAddOperation {
		// For add operations, append the new content to the original content
		updatedContent := originalContent + "\n\n" + editedContent
		if c.sourceManager.UpdateSourceContent(sourceId, updatedContent) {
			fmt.Println("\n📝 New content added to the source file!")
			
			// After updating the file, we should re-ingest it to update the chunks
			source := c.sourceManager.GetSourceByID(sourceId)
			if source != nil {
				fmt.Println("\n💡 Tip: Run '/ingest' to update the knowledge base with the new content.")
				fmt.Printf("   Source: %d | %s | %s\n", source.Number, source.ID[:4], source.Title)
			}
		}
	} else {
		// For partial edits, try to find and replace specific content
		var updatedContent string
		
		// Strategy 1: Try to find exact matches of relevant chunks
		for _, chunk := range relevantChunks {
			if strings.Contains(originalContent, chunk.Text) {
				updatedContent = strings.Replace(originalContent, chunk.Text, editedContent, 1)
				break
			}
		}
		
		// Strategy 2: If no exact match, try fuzzy matching for partial content
		if updatedContent == "" || updatedContent == originalContent {
			updatedContent = c.attemptFuzzyReplacement(originalContent, relevantChunks, editedContent)
		}
		
		// If we found content to update, save it
		if updatedContent != "" && updatedContent != originalContent {
			if c.sourceManager.UpdateSourceContent(sourceId, updatedContent) {
				fmt.Println("\n📝 Source file updated successfully!")
				fmt.Println("The changes have been saved to the original document.")
				
				// After updating the file, we should re-ingest it to update the chunks
				source := c.sourceManager.GetSourceByID(sourceId)
				if source != nil {
					fmt.Println("\n💡 Tip: Run '/ingest' to update the knowledge base with the new content.")
					fmt.Printf("   Source: %d | %s | %s\n", source.Number, source.ID[:4], source.Title)
				}
			}
		} else {
			// As a fallback, if we can't do precise replacement, provide the content for manual update
			return c.showProposedContent(editedContent)
		}
	}
	
	return nil
}

func (c *CLIInterface) showProposedContent(editedContent string) error {
	fmt.Println("\n⚠️  Could not automatically update the source file.")
	fmt.Println("💡 The AI has generated the updated content. Consider manual update:")
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│                PROPOSED FILE CONTENT                  │")
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	
	// Display the edited content with proper formatting
	allLines := strings.Split(editedContent, "\n")
	
	for _, line := range allLines {
		if len(line) <= 68 {
			fmt.Printf("│ %-68s │\n", line)
		} else {
			// Wrap long lines
			for len(line) > 0 {
				var wrappedLine string
				if len(line) <= 68 {
					wrappedLine = line
					line = ""
				} else {
					wrappedLine = line[:68]
					line = line[68:]
				}
				fmt.Printf("│ %-68s │\n", wrappedLine)
			}
		}
	}
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	
	// Show source information
	activeSources := c.sourceManager.GetActiveSources()
	if len(activeSources) > 0 {
		fmt.Println("\n💡 Tip: Run '/ingest' to update the knowledge base after manual changes.")
		var sourceInfo []string
		for _, s := range activeSources {
			sourceInfo = append(sourceInfo, fmt.Sprintf("%d | %s | %s", s.Number, s.ID[:4], s.Title))
		}
		fmt.Printf("   Active sources: %s\n", strings.Join(sourceInfo, ", "))
	}
	
	return nil
}

func (c *CLIInterface) isCompleteRewriteInstruction(instruction string) bool {
	deleteKeywords := []string{"delete", "remove", "erase", "eliminate", "get rid of", "clear out"}
	rewriteKeywords := []string{"rewrite", "reformat", "restructure", "update completely"}
	
	lowerInstruction := strings.ToLower(instruction)
	
	for _, keyword := range deleteKeywords {
		if strings.Contains(lowerInstruction, keyword) {
			return true
		}
	}
	
	for _, keyword := range rewriteKeywords {
		if strings.Contains(lowerInstruction, keyword) {
			return true
		}
	}
	
	return false
}

func (c *CLIInterface) isCompleteDocumentContent(editedContent string, originalContent string) bool {
	// Check if the edited content contains structural elements that suggest a complete document
	hasMultipleSections := strings.Count(editedContent, "\n") > 3
	hasListMarkers := regexp.MustCompile(`[-*•]\s|^\d+\.`).MatchString(editedContent)
	hasHeaders := regexp.MustCompile(`^\s*#+\s|\b[A-Z][A-Z\s]*:\s*$`).MatchString(editedContent)
	
	// If it has document-like structure and is reasonably complete compared to original
	return (hasMultipleSections || hasListMarkers || hasHeaders) &&
		float64(len(editedContent)) > float64(len(originalContent))*0.5
}

func (c *CLIInterface) attemptFuzzyReplacement(originalContent string, relevantChunks []ingestion.Chunk, editedContent string) string {
	// Try to find the best matching section to replace
	for _, chunk := range relevantChunks {
		// Try different levels of matching
		matches := c.findBestMatchingSection(originalContent, chunk.Text)
		if len(matches) > 0 {
			// Replace the best matching section
			bestMatch := matches[0]
			return originalContent[:bestMatch.start] +
				editedContent +
				originalContent[bestMatch.end:]
		}
	}
	return ""
}

type match struct {
	start int
	end   int
	score float64
}

func (c *CLIInterface) findBestMatchingSection(content string, target string) []match {
	targetLines := strings.FieldsFunc(target, func(c rune) bool { return c == '\n' })
	
	var matches []match
	
	for i := 0; i < len(content); i++ {
		for _, targetLine := range targetLines {
			if strings.Contains(content[i:], targetLine) {
				// Find surrounding context to get a larger section
				// This is a simplified version - in a full implementation you'd want more sophisticated matching
				matchStart := i
				matchEnd := i + len(targetLine)
				score := c.calculateSimilarity(content[matchStart:matchEnd], target)
				matches = append(matches, match{start: matchStart, end: matchEnd, score: score})
			}
		}
	}
	
	// Sort by score (highest first)
	for i := 0; i < len(matches)-1; i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].score < matches[j].score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}
	
	return matches
}

func (c *CLIInterface) calculateSimilarity(text1 string, text2 string) float64 {
	words1 := make(map[string]bool)
	for _, word := range strings.Fields(text1) {
		words1[strings.ToLower(word)] = true
	}
	
	words2 := make(map[string]bool)
	for _, word := range strings.Fields(text2) {
		words2[strings.ToLower(word)] = true
	}
	
	intersection := 0
	for word := range words1 {
		if words2[word] {
			intersection++
		}
	}
	
	union := len(words1) + len(words2) - intersection
	
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

func (c *CLIInterface) handleQuestion(question string) error {
	var relevantChunks []ingestion.Chunk
	
	if c.groundingMode != "general" {
		// Get active sources and retrieve relevant chunks
		activeSourceIds := c.sourceManager.GetActiveSourceIds()
		retrievalResult := c.retriever.Retrieve(question, activeSourceIds)
		relevantChunks = retrievalResult.Chunks
		
		// In strict mode, only return if no chunks found
		if len(relevantChunks) == 0 && c.groundingMode == "strict" {
			fmt.Println("In strict mode: No relevant information found in the provided sources to answer this question.")
			return nil
		}
	}

	ctx := c.ctxManager.GetContext()
	answer, err := c.llmClient.Answer(context.Background(), question, relevantChunks, ctx)
	if err != nil {
		return fmt.Errorf("error processing question: %w", err)
	}
	
	fmt.Println("\n┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│                     🤖 AI RESPONSE                     │")
	fmt.Println("├─────────────────────────────────────────────────────────┤")
	
	// Split the answer into lines and handle each line
	allLines := strings.Split(answer, "\n")
	
	for i, line := range allLines {
		if i == 0 {
			// Handle the first line specially
			if len(line) <= 68 {
				fmt.Printf("│ %-68s │\n", line)
			} else {
				// Split first line if too long
				firstPart := line[:68]
				fmt.Printf("│ %-68s │\n", firstPart)
				
				// Handle remaining parts of the first line
				remaining := line[68:]
				for len(remaining) > 0 {
					var wrappedLine string
					if len(remaining) <= 68 {
						wrappedLine = remaining
						remaining = ""
					} else {
						wrappedLine = remaining[:68]
						remaining = remaining[68:]
					}
					fmt.Printf("│ %-68s │\n", wrappedLine)
				}
			}
		} else {
			// Handle subsequent lines
			if len(line) <= 68 {
				fmt.Printf("│ %-68s │\n", line)
			} else {
				// Wrap long lines
				for len(line) > 0 {
					var wrappedLine string
					if len(line) <= 68 {
						wrappedLine = line
						line = ""
					} else {
						wrappedLine = line[:68]
						line = line[68:]
					}
					fmt.Printf("│ %-68s │\n", wrappedLine)
				}
			}
		}
	}
	fmt.Println("└─────────────────────────────────────────────────────────┘\n")
	
	return nil
}

func (c *CLIInterface) handleModeCommand(args []string) {
	if len(args) == 0 {
		fmt.Printf("Current mode: %s\n", c.groundingMode)
		fmt.Println("Available modes: strict, hybrid, general")
		return
	}

	mode := strings.ToLower(args[0])
	if mode == "strict" || mode == "hybrid" || mode == "general" {
		c.groundingMode = mode
		fmt.Printf("Mode set to: %s\n", mode)
	} else {
		fmt.Printf("Invalid mode: %s. Use strict, hybrid, or general\n", mode)
	}
}

func (c *CLIInterface) handleIngestCommand() error {
	activeSources := c.sourceManager.GetActiveSources()
	
	if len(activeSources) == 0 {
		fmt.Println("No active sources to ingest. Use /source use <id> to activate sources first.")
		return nil
	}

	fmt.Printf("Processing %d active source(s)...\n", len(activeSources))
	// Add a list of the sources being processed
	if len(activeSources) > 0 {
		var sourceList []string
		for _, source := range activeSources {
			sourceList = append(sourceList, fmt.Sprintf("  - %d | %s", source.Number, source.Title))
		}
		fmt.Printf("Sources to be processed:\n%s\n\n", strings.Join(sourceList, "\n"))
	}
	
	for _, source := range activeSources {
		fmt.Printf("\nProcessing source: %d | %s\n", source.Number, source.Title)
		
		// Clear old chunks from this source before adding new ones
		c.chunkManager.RemoveChunksBySource(source.ID)
		
		result, err := c.ingestionPipeline.ProcessSource(source)
		if err != nil {
			fmt.Printf("✗ Error processing %d | %s: %v\n", source.Number, source.Title, err)
			continue
		}
		
		if result.Status == "success" {
			c.chunkManager.AddChunks(result.Chunks)
			fmt.Printf("✓ Successfully processed %d chunks from %d | %s\n", len(result.Chunks), source.Number, source.Title)
		} else {
			fmt.Printf("✗ Error processing %d | %s: %s\n", source.Number, source.Title, result.Error)
		}
	}
	
	fmt.Println("\nIngestion completed!")
	return nil
}

func (c *CLIInterface) showHelp() {
	fmt.Println("\n📚 Available Commands:")
	fmt.Println("/help - Show this help message")
	fmt.Println("/project set <name> - Set the current project")
	fmt.Println("/stage set <name> - Set the current stage")
	fmt.Println("/task set <title> - Set the current task")
	fmt.Println("/context show - Show current context")
	fmt.Println("/context reset - Reset current context")
	fmt.Println("/source add file <path> - Add a file source")
	fmt.Println("/source add url <url> - Add a URL source")
	fmt.Println("/source list - List all sources")
	fmt.Println("/source use <id1,id2> - Set active sources (supports numbers: 1,2 or IDs: abc1,def2)")
	fmt.Println("/source use all - Set all sources as active")
	fmt.Println("/source clear - Clear active sources")
	fmt.Println("/source remove <id> - Completely remove a source from the system")
	fmt.Println("/ingest - Process active sources into chunks")
	fmt.Println("/redact <source_id> <pattern> - Redact sensitive information from a source")
	fmt.Println("/mode strict|hybrid|general - Set grounding mode")
	fmt.Println("/exit or /quit - Exit the bot")
	fmt.Println("")
	fmt.Println("💡 Tips:")
	fmt.Println("- Add sources first with /source add")
	fmt.Println("- Use /source use to activate sources for retrieval")
	fmt.Println("- Run /ingest to process sources into searchable chunks")
	fmt.Println("- Use /source remove <id> to completely remove a source from the system")
	fmt.Println("- Sources can be referenced by number (1, 2, 3) or partial ID (abc1, def2)")
	fmt.Println("- Ask questions directly (without /) for Q&A mode")
	fmt.Println("")
}