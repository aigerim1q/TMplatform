# Chat-Bot Optimization Plan

## Current System Analysis

The current chat-bot implementation has several areas for improvement:

1. **Intent Detection**: Current system struggles to distinguish between editing file requests and question-answering requests
2. **Source ID System**: Currently uses UUIDs with 4-character prefixes, which are not user-friendly
3. **User Experience**: Complex command structure and limited error feedback
4. **Performance**: Basic retrieval algorithm without optimizations
5. **File Editing**: Limited ability to accurately modify file content based on user instructions

## Proposed Improvements

### 1. Enhanced Intent Detection System

#### Improved Intent Classification
- Add machine learning-based intent classification
- Better context awareness for distinguishing edit vs question requests
- Natural language understanding for complex requests

**Implementation in `src/cli/repl.ts`:**
```typescript
private detectIntent(input: string): 'edit' | 'question' {
  const inputLower = input.toLowerCase();
  
  // Weighted scoring for intent detection
  let editScore = 0;
  let questionScore = 0;
  
  // Strong indicators of edit requests with weights
  const editKeywords = [
    { keyword: 'edit', weight: 2 },
    { keyword: 'change', weight: 2 },
    { keyword: 'modify', weight: 2 },
    { keyword: 'update', weight: 2 },
    { keyword: 'rewrite', weight: 2 },
    { keyword: 'rephrase', weight: 2 },
    { keyword: 'fix', weight: 2 },
    { keyword: 'correct', weight: 2 },
    { keyword: 'adjust', weight: 2 },
    { keyword: 'alter', weight: 2 },
    { keyword: 'improve', weight: 2 },
    { keyword: 'replace', weight: 2 },
    { keyword: 'insert', weight: 2 },
    { keyword: 'remove', weight: 2 },
    { keyword: 'delete', weight: 2 },
    { keyword: 'append', weight: 2 },
    { keyword: 'add', weight: 3 }, // Adding has high weight as it's clearly an edit operation
    { keyword: 'create', weight: 3 },
    { keyword: 'generate', weight: 2 }
  ];
  
  // Strong indicators of questions with weights
  const questionKeywords = [
    { keyword: 'what', weight: 2 },
    { keyword: 'who', weight: 2 },
    { keyword: 'where', weight: 2 },
    { keyword: 'when', weight: 2 },
    { keyword: 'why', weight: 2 },
    { keyword: 'how', weight: 2 },
    { keyword: 'is', weight: 1 },
    { keyword: 'are', weight: 1 },
    { keyword: 'can', weight: 1 },
    { keyword: 'could', weight: 1 },
    { keyword: 'would', weight: 1 },
    { keyword: 'should', weight: 1 },
    { keyword: 'does', weight: 1 },
    { keyword: 'do', weight: 1 },
    { keyword: 'did', weight: 1 },
    { keyword: 'have', weight: 1 },
    { keyword: 'has', weight: 1 },
    { keyword: 'had', weight: 1 },
    { keyword: 'explain', weight: 2 },
    { keyword: 'describe', weight: 2 },
    { keyword: 'summarize', weight: 2 }
  ];
  
  // Check for edit keywords
  for (const { keyword, weight } of editKeywords) {
    if (inputLower.includes(keyword)) {
      editScore += weight;
    }
  }
  
  // Check for question keywords
  for (const { keyword, weight } of questionKeywords) {
    if (inputLower.includes(keyword)) {
      questionScore += weight;
    }
  }
  
  // Edit patterns with weights
  const editPatterns = [
    { pattern: /^edit\s+/i, weight: 4 },                              // "edit this"
    { pattern: /^change\s+/i, weight: 4 },                            // "change this"
    { pattern: /^update\s+/i, weight: 4 },                            // "update this"
    { pattern: /change.*to/i, weight: 3 },                           // "change X to Y"
    { pattern: /replace.*with/i, weight: 3 },                        // "replace X with Y"
    { pattern: /update.*to/i, weight: 3 },                           // "update X to Y"
    { pattern: /fix.*to/i, weight: 3 },                              // "fix X to Y"
    { pattern: /correct.*to/i, weight: 3 },                          // "correct X to Y"
    { pattern: /make.*be/i, weight: 3 },                             // "make this be X"
    { pattern: /set.*to/i, weight: 3 },                              // "set this to X"
    { pattern: /add\s+.*\s+and\s+.*\s+is\s+to/i, weight: 5 },       // "add John and his job is to" - high weight for this specific pattern
    { pattern: /add\s+.*\s+to\s+the\s+document/i, weight: 4 },       // "add John to the document"
    { pattern: /add\s+a?\s*person\s+named/i, weight: 4 },            // "add a person named"
    { pattern: /add\s+.*\s+with\s+task/i, weight: 4 },               // "add John with task"
    { pattern: /add\s+.*\s+and\s+assign/i, weight: 4 },              // "add John and assign"
    { pattern: /add\s+\w+\s+and\s+.*\s+is\s+/i, weight: 5 },        // "add X and his/their Y is Z" - pattern matching the user's example
  ];
  
  // Question patterns with weights
  const questionPatterns = [
    { pattern: /^(what|who|where|when|why|how)\s+/i, weight: 3 },     // "what is", "who is", etc.
    { pattern: /.*\?$/, weight: 3 },                                 // Ends with question mark
    { pattern: /^(is|are|can|could|would|should)\s+/i, weight: 2 },   // "is this", "can you", etc.
  ];
  
  // Apply pattern weights
  for (const { pattern, weight } of editPatterns) {
    if (pattern.test(input)) {
      editScore += weight;
    }
  }
  
  for (const { pattern, weight } of questionPatterns) {
    if (pattern.test(input)) {
      questionScore += weight;
    }
  }
  
  // Additional context analysis
  if (input.trim().endsWith('?')) {
    questionScore += 2;
  }
  
  if (this.containsEditIndicators(inputLower)) {
    editScore += 1;
  }
  
  // Special handling for the specific case from user feedback
  // "add Aidana and her task is similar to Beka aga" should be detected as edit
  if (inputLower.includes('add') && inputLower.includes('and') && inputLower.includes('is')) {
    // This pattern often indicates adding a person/task with details
    editScore += 2;
  }
  
  // Return the intent with higher score, but with some adjustments for safety
  if (editScore > questionScore) {
    return 'edit';
  } else if (questionScore > editScore) {
    return 'question';
  } else {
    // When scores are equal, if there's any indication of adding/modifying, lean toward edit
    // This helps with cases like "add X and Y is Z"
    if (editScore > 0) {
      return 'edit';
    } else {
      // Default to question for safety when truly ambiguous
      return 'question';
    }
  }
}

private containsEditIndicators(inputLower: string): boolean {
  // Additional checks for edit indicators that might not be caught by keywords
  const editIndicators = [
    'make it',
    'make this',
    'turn this',
    'convert this',
    'transform',
    'reformat',
    'restructure',
    'add to',
    'insert into',
    'put in',
    'place in',
    'update the',
    'modify the'
  ];
  
  return editIndicators.some(indicator => inputLower.includes(indicator));
}
```

### 2. Improved File Editing System

#### Enhanced Text Replacement Strategy
- More robust text replacement with fuzzy matching
- Better handling of partial matches and context preservation
- Improved file update mechanism

**Implementation in `src/cli/repl.ts`:**
```typescript
private async handleEditRequest(instruction: string): Promise<void> {
  try {
    // Get relevant chunks based on the instruction
    let relevantChunks: Chunk[] = [];
    if (this.groundingMode !== 'general') {
      const activeSourceIds = this.sourceManager.getActiveSources().map(s => s.id);
      // Use the instruction to find relevant chunks
      const retrievalResult = this.retriever.retrieve(instruction, activeSourceIds);
      relevantChunks = retrievalResult.chunks;
    }

    const context = this.contextManager.getContext();
    
    // Determine if this is an append/add operation vs a modification
    const instructionLower = instruction.toLowerCase();
    const isAddOperation = instructionLower.includes('add') || 
                          instructionLower.includes('create') || 
                          instructionLower.includes('insert');
    
    let textToEdit = '';
    let editInstruction = instruction;
    
    if (isAddOperation && relevantChunks.length === 0) {
      // For add operations with no relevant chunks, we'll work with the original source content
      const activeSources = this.sourceManager.getActiveSources();
      if (activeSources.length > 0) {
        // Get the content of the first active source to use as base
        const firstSourceId = activeSources[0].id;
        const originalContent = this.sourceManager.getSourceContent(firstSourceId);
        if (originalContent) {
          textToEdit = originalContent;
        }
      }
    } else if (relevantChunks.length > 0) {
      // Combine relevant chunks as potential text to edit
      textToEdit = relevantChunks.map(c => c.text).join('\n\n');
    } else {
      // If no relevant chunks and not an add operation, treat as a general edit
      textToEdit = instruction;
    }
    
    // Call the LLM to perform the edit
    const result = await this.llm.edit(textToEdit, editInstruction, relevantChunks, context);
    
    console.log('\n┌─────────────────────────────────────────────────────────┐');
    console.log('│                    📝 EDIT RESULT                     │');
    console.log('├─────────────────────────────────────────────────────────┤');
    
    // Split the result into lines and handle each line
    const allLines = result.split('\n');
    let firstLine = true;
    
    for (let i = 0; i < allLines.length; i++) {
      const line = allLines[i];
      
      if (firstLine && i === 0) {
        // Handle the first line specially - use substring approach for the first line only
        if (line.length <= 68) {
          console.log(`│ ${line.padEnd(68)} │`);
        } else {
          // Split first line if too long
          const firstPart = line.substring(0, 68);
          console.log(`│ ${firstPart.padEnd(68)} │`);
          
          // Handle remaining parts of the first line
          const remaining = line.substring(68);
          const wrappedLines = remaining.match(/.{1,68}/g) || [];
          for (const wrappedLine of wrappedLines) {
            console.log(`│ ${wrappedLine.padEnd(68)} │`);
          }
        }
        firstLine = false;
      } else {
        // Handle subsequent lines
        if (line.length <= 68) {
          console.log(`│ ${line.padEnd(68)} │`);
        } else {
          // Wrap long lines
          const wrappedLines = line.match(/.{1,68}/g) || [];
          for (const wrappedLine of wrappedLines) {
            console.log(`│ ${wrappedLine.padEnd(68)} │`);
          }
        }
      }
    }
    console.log('└─────────────────────────────────────────────────────────┘\n');
    
    // Attempt to update the source files with the edited content
    if (result !== textToEdit) {
      await this.attemptToUpdateSourceFiles(relevantChunks, result);
    }
    
  } catch (error) {
    console.error('Error during edit:', error);
  }
}

private async attemptToUpdateSourceFiles(relevantChunks: Chunk[], editedContent: string): Promise<void> {
  // Get unique source IDs from the relevant chunks
  const sourceIds = [...new Set(relevantChunks.map(chunk => chunk.sourceId))];
  
  // If no source IDs from chunks, try to get from active sources
  if (sourceIds.length === 0) {
    const activeSources = this.sourceManager.getActiveSources();
    if (activeSources.length > 0) {
      // Use the first active source for add operations
      const firstSourceId = activeSources[0].id;
      await this.updateSingleSourceFile(firstSourceId, editedContent, relevantChunks);
    } else {
      // If no active sources, show the content for manual update
      await this.showProposedContent(editedContent);
    }
  } else {
    // Update files from relevant chunks
    for (const sourceId of sourceIds) {
      await this.updateSingleSourceFile(sourceId, editedContent, relevantChunks);
    }
  }
}

private async updateSingleSourceFile(sourceId: string, editedContent: string, relevantChunks: Chunk[]): Promise<void> {
  // Get the original content of the source
  const originalContent = this.sourceManager.getSourceContent(sourceId);
  if (!originalContent) {
    return;
  }
  
  try {
    // Determine if this is likely a complete file rewrite by analyzing the edit request
    // If the edited content looks like a complete document (contains multiple lines/sections)
    // and is significantly different from individual chunks, treat it as a complete replacement
    
    const instruction = relevantChunks.length > 0 ? relevantChunks[0].text : '';
    const isLikelyCompleteRewrite = this.isCompleteRewriteInstruction(instruction) ||
                                   this.isCompleteDocumentContent(editedContent, originalContent);
    
    // Check if the edit request is an add operation
    const isAddOperation = instruction.toLowerCase().includes('add') || 
                          instruction.toLowerCase().includes('create') || 
                          instruction.toLowerCase().includes('insert');
    
    if (isLikelyCompleteRewrite) {
      // For complete rewrites (like deleting items), replace the entire file
      if (this.sourceManager.updateSourceContent(sourceId, editedContent)) {
        console.log(`\n📝 Source file updated successfully!`);
        console.log(`The changes have been saved to the original document.`);
        
        // After updating the file, we should re-ingest it to update the chunks
        const source = this.sourceManager.getSourceById(sourceId);
        if (source) {
          console.log(`\n💡 Tip: Run '/ingest' to update the knowledge base with the new content.`);
          console.log(`   Source: ${source.number} | ${source.id.substring(0, 4)} | ${source.title}`);
        }
      }
    } else if (isAddOperation) {
      // For add operations, append the new content to the original content
      const updatedContent = originalContent + '\n\n' + editedContent;
      if (this.sourceManager.updateSourceContent(sourceId, updatedContent)) {
        console.log(`\n📝 New content added to the source file!`);
        
        // After updating the file, we should re-ingest it to update the chunks
        const source = this.sourceManager.getSourceById(sourceId);
        if (source) {
          console.log(`\n💡 Tip: Run '/ingest' to update the knowledge base with the new content.`);
          console.log(`   Source: ${source.number} | ${source.id.substring(0, 4)} | ${source.title}`);
        }
      }
    } else {
      // For partial edits, try to find and replace specific content
      let updatedContent: string | null = null;
      
      // Strategy 1: Try to find exact matches of relevant chunks
      for (const chunk of relevantChunks) {
        if (originalContent.includes(chunk.text)) {
          updatedContent = originalContent.replace(chunk.text, editedContent);
          break;
        }
      }
      
      // Strategy 2: If no exact match, try fuzzy matching for partial content
      if (!updatedContent || updatedContent === originalContent) {
        updatedContent = this.attemptFuzzyReplacement(originalContent, relevantChunks, editedContent);
      }
      
      // If we found content to update, save it
      if (updatedContent && updatedContent !== originalContent) {
        if (this.sourceManager.updateSourceContent(sourceId, updatedContent)) {
          console.log(`\n📝 Source file updated successfully!`);
          console.log(`The changes have been saved to the original document.`);
          
          // After updating the file, we should re-ingest it to update the chunks
          const source = this.sourceManager.getSourceById(sourceId);
          if (source) {
            console.log(`\n💡 Tip: Run '/ingest' to update the knowledge base with the new content.`);
            console.log(`   Source: ${source.number} | ${source.id.substring(0, 4)} | ${source.title}`);
          }
        }
      } else {
        // As a fallback, if we can't do precise replacement, provide the content for manual update
        await this.showProposedContent(editedContent);
      }
    }
  } catch (error) {
    console.error(`Error updating source file ${sourceId}:`, error);
  }
}

private async showProposedContent(editedContent: string): Promise<void> {
  console.log(`\n⚠️  Could not automatically update the source file.`);
  console.log(`💡 The AI has generated the updated content. Consider manual update:`);
  console.log('┌─────────────────────────────────────────────────────────┐');
  console.log('│                PROPOSED FILE CONTENT                  │');
  console.log('├─────────────────────────────────────────────────────────┤');
  
  // Display the edited content with proper formatting
  const allLines = editedContent.split('\n');
  
  for (let i = 0; i < allLines.length; i++) {
    const line = allLines[i];
    
    if (line.length <= 68) {
      console.log(`│ ${line.padEnd(68)} │`);
    } else {
      // Wrap long lines
      const wrappedLines = line.match(/.{1,68}/g) || [];
      for (const wrappedLine of wrappedLines) {
        console.log(`│ ${wrappedLine.padEnd(68)} │`);
      }
    }
  }
  console.log('└─────────────────────────────────────────────────────────┘');
  
  // Show source information
  const activeSources = this.sourceManager.getActiveSources();
  if (activeSources.length > 0) {
    console.log(`\n💡 Tip: Run '/ingest' to update the knowledge base after manual changes.`);
    console.log(`   Active sources: ${activeSources.map(s => \`${s.number} | ${s.id.substring(0, 4)} | ${s.title}\`).join(', ')}`);
  }
}
```

### 3. Enhanced Retrieval System

#### Improved Similarity Algorithm
- Implement proper TF-IDF or BM25 ranking
- Add caching for frequently accessed content
- Optimize chunk storage and retrieval

**Implementation in `src/retrieval/retriever.ts`:**
```typescript
private calculateSimilarity(query: string, text: string): number {
  // Enhanced similarity calculation with multiple strategies
  const queryWords = this.tokenize(query.toLowerCase());
  const textWords = this.tokenize(text.toLowerCase());
  
  if (queryWords.length === 0 || textWords.length === 0) {
    return 0;
  }
  
  // Calculate term frequency-inverse document frequency (TF-IDF) like scoring
  const queryWordSet = new Set(queryWords);
  const textWordSet = new Set(textWords);
  
  // Calculate Jaccard similarity (intersection over union)
  let overlap = 0;
  for (const word of queryWordSet) {
    if (textWordSet.has(word)) {
      overlap++;
    }
  }
  
  const jaccard = overlap / Math.max(queryWordSet.size, textWordSet.size, 1);
  
  // Calculate overlap ratio relative to query size
  const queryOverlapRatio = overlap / Math.max(queryWordSet.size, 1);
  
  // Calculate cosine similarity between term vectors
  const commonWords = new Set([...queryWordSet].filter(x => textWordSet.has(x)));
  const cosineSimilarity = commonWords.size / (Math.sqrt(queryWordSet.size) * Math.sqrt(textWordSet.size));
  
  // Use a weighted combination of different similarity measures
  return (jaccard * 0.3 + queryOverlapRatio * 0.4 + cosineSimilarity * 0.3);
}
```

### 4. Caching Mechanism

#### Performance Optimization
- Add LRU cache for retrieval results
- Cache frequently accessed chunks
- Implement lazy loading for large documents

**Implementation:**
```typescript
// In src/retrieval/retriever.ts
import { LRUCache } from 'lru-cache';

export class Retriever {
  private chunkManager: ChunkManager;
  private options: RetrievalOptions;
  private cache: LRUCache<string, RetrievalResult>;

  constructor(chunkManager: ChunkManager, options?: Partial<RetrievalOptions>) {
    this.chunkManager = chunkManager;
    this.options = {
      topK: options?.topK || 5,
      minScore: options?.minScore || 0.01,
    };
    
    // Initialize cache with 100 entries max
    this.cache = new LRUCache<string, RetrievalResult>({ max: 100 });
  }

  public retrieve(query: string, sourceIds?: string[]): RetrievalResult {
    // Create cache key based on query and source IDs
    const cacheKey = `${query}_${sourceIds ? sourceIds.sort().join(',') : 'all'}`;
    
    // Check if result is already cached
    const cachedResult = this.cache.get(cacheKey);
    if (cachedResult) {
      return cachedResult;
    }
    
    // Perform the actual retrieval
    const result = this.performRetrieval(query, sourceIds);
    
    // Cache the result
    this.cache.set(cacheKey, result);
    
    return result;
  }
  
  private performRetrieval(query: string, sourceIds?: string[]): RetrievalResult {
    // Existing retrieval logic here
    // ... (same as current implementation)
  }
}
```

### 5. User Experience Improvements

#### Command Aliases and Shortcuts
- Add shortcuts for common commands
- Natural language support
- Better error messages

**Implementation in `src/cli/repl.ts`:**
```typescript
private async handleCommand(input: string): Promise<void> {
  const [command, ...args] = input.split(' ');
  let commandName = command.substring(1).toLowerCase(); // Remove the '/' and lowercase

  // Map aliases to actual commands
  const commandAliases: { [key: string]: string } = {
    'add': 'source',
    'list': 'source',
    'use': 'source',
    'ls': 'source',
    'rm': 'source',
    'q': 'context',
    'p': 'project',
    'st': 'stage',
    't': 'task',
    'ing': 'ingest',
    'ingest': 'ingest'
  };

  if (commandAliases[commandName]) {
    // For source aliases, we need to adjust the args
    if (commandAliases[commandName] === 'source') {
      args.unshift(commandName); // Put the alias back as first arg
      commandName = 'source';
    } else {
      commandName = commandAliases[commandName];
    }
  }

  switch (commandName) {
    case 'help':
      this.showHelp();
      break;
    case 'project':
      console.log(this.contextCommands.handleProjectCommand(args));
      break;
    case 'stage':
      console.log(this.contextCommands.handleStageCommand(args));
      break;
    case 'task':
      console.log(this.contextCommands.handleTaskCommand(args));
      break;
    case 'context':
      console.log(this.contextCommands.handleContextCommand(args));
      break;
    case 'source':
      console.log(this.sourceCommands.handleSourceCommand(args));
      break;
    case 'mode':
      this.handleModeCommand(args);
      break;
    case 'ingest':
      await this.handleIngestCommand();
      break;
    case 'redact':
      await this.handleRedactCommand(args);
      break;
    default:
      console.log(`Unknown command: ${commandName}. Type /help for available commands.`);
  }
}
```

### 6. Fuzzy Matching for Source IDs

#### Enhanced ID Matching
- Implement fuzzy matching for source IDs to handle typos
- Support partial matches and similarity scoring
- Provide suggestions when exact matches aren't found

**Implementation in `src/sources/manager.ts`:**
```typescript
public getSourceById(id: string): KnowledgeSource | undefined {
  // First try exact UUID match
  const exactMatch = this.state.allSources.find(source => source.id === id);
  if (exactMatch) {
    return exactMatch;
  }
  
  // Then try partial UUID match (for backward compatibility with 4-char prefixes)
  const partialMatches = this.state.allSources.filter(source => source.id.startsWith(id));
  if (partialMatches.length === 1) {
    return partialMatches[0];
  }
  
  // Then try sequential number match
  const number = parseInt(id, 10);
  if (!isNaN(number)) {
    const numberMatch = this.state.allSources.find(source => source.number === number);
    if (numberMatch) {
      return numberMatch;
    }
  }
  
  // Try fuzzy matching for typos
  const fuzzyMatch = this.findFuzzyMatch(id);
  if (fuzzyMatch) {
    return fuzzyMatch;
  }
  
  return undefined;
}

private findFuzzyMatch(input: string): KnowledgeSource | undefined {
  // Calculate similarity scores for all sources and return the best match
  const candidates = this.state.allSources;
  let bestMatch: KnowledgeSource | undefined = undefined;
  let bestScore = 0;
  
  for (const source of candidates) {
    // Calculate similarity with both ID and title
    const idScore = this.calculateStringSimilarity(input, source.id);
    const titleScore = this.calculateStringSimilarity(input, source.title);
    const pathScore = this.calculateStringSimilarity(input, source.path);
    
    // Use the highest score among all identifiers
    const maxScore = Math.max(idScore, titleScore, pathScore);
    
    if (maxScore > bestScore && maxScore > 0.6) { // Only return matches above threshold
      bestScore = maxScore;
      bestMatch = source;
    }
  }
  
  return bestMatch;
}

private calculateStringSimilarity(str1: string, str2: string): number {
  // Calculate similarity using a simple algorithm (Levenshtein distance normalized)
  const longer = str1.length > str2.length ? str1 : str2;
  const shorter = str1.length > str2.length ? str2 : str1;
  
  if (longer.length === 0) {
    return 1.0;
  }
  
  const editDistance = this.levenshteinDistance(longer.toLowerCase(), shorter.toLowerCase());
  return (longer.length - editDistance) / longer.length;
}

private levenshteinDistance(str1: string, str2: string): number {
  const matrix = Array(str2.length + 1).fill(0).map(() => Array(str1.length + 1).fill(0));

  for (let i = 0; i <= str1.length; i++) {
    matrix[0][i] = i;
  }

  for (let j = 0; j <= str2.length; j++) {
    matrix[j][0] = j;
  }

  for (let j = 1; j <= str2.length; j++) {
    for (let i = 1; i <= str1.length; i++) {
      const cost = str1[i - 1] === str2[j - 1] ? 0 : 1;
      matrix[j][i] = Math.min(
        matrix[j][i - 1] + 1,     // deletion
        matrix[j - 1][i] + 1,     // insertion
        matrix[j - 1][i - 1] + cost // substitution
      );
    }
  }

  return matrix[str2.length][str1.length];
}
```

## Implementation Roadmap

### Phase 1: Core Intent Detection and Editing Improvements
1. Update intent detection algorithm in `src/cli/repl.ts`
2. Enhance file editing functionality with better replacement strategies
3. Improve error handling and feedback
4. Test intent detection with various inputs

### Phase 2: Performance Optimizations
1. Implement caching mechanisms for retrieval
2. Optimize similarity calculations
3. Add performance monitoring
4. Benchmark before and after improvements

### Phase 3: User Experience Enhancements
1. Add command aliases and shortcuts
2. Implement fuzzy matching for source IDs
3. Improve error messages and user feedback
4. Add command suggestions

### Phase 4: Advanced Features
1. Context-aware features
2. Advanced error recovery
3. Natural language processing improvements
4. Additional optimization based on user feedback

## Expected Benefits

1. **Better Intent Recognition**: Improved accuracy in distinguishing between edit and question requests
2. **Enhanced File Editing**: More robust and accurate file modification capabilities
3. **Improved Performance**: Faster query processing with caching mechanisms
4. **Better UX**: Natural language processing and error handling improve user experience
5. **Robustness**: Fuzzy matching and error recovery make the system more resilient

## Success Metrics

- **Intent Detection Accuracy**: Higher success rate in correctly identifying edit vs question requests
- **File Update Success Rate**: More successful automatic file updates
- **Response Time**: Faster query processing with optimized algorithms
- **User Satisfaction**: Easier to use with better error messages and feedback
- **Error Rate**: Reduced user errors with better feedback and suggestions