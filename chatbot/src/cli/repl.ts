import * as readline from 'readline';
import { ContextManager } from '../context/manager';
import { ContextCommands } from '../context/commands';
import { SourceManager } from '../sources/manager';
import { SourceCommands } from '../sources/commands';
import { IngestionPipeline } from '../ingestion/pipeline';
import { ChunkManager } from '../ingestion/chunkManager';
import { Retriever } from '../retrieval/retriever';
import { LLMClient } from '../llm/types';
import { Context } from '../context/types';
import { Chunk } from '../ingestion/types';
import { createLLMClient } from './llmFactory';
import { getChatbotApiUrl, sendToChatbot, ChatSession } from './chatApi';

export class CLIInterface {
  private rl: readline.Interface;
  private contextManager: ContextManager;
  private contextCommands: ContextCommands;
  private sourceManager: SourceManager;
  private sourceCommands: SourceCommands;
  private ingestionPipeline: IngestionPipeline;
  private chunkManager: ChunkManager;
  private retriever: Retriever;
  private llm: LLMClient;
  private groundingMode: 'strict' | 'hybrid' | 'general';
  /** Session state for chatbot API (active project) when CHATBOT_API_URL is set */
  private chatSession: ChatSession = { active_project_id: '', active_project_title: '' };

  constructor() {
    this.rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout
    });

    this.contextManager = new ContextManager();
    this.contextCommands = new ContextCommands(this.contextManager);
    
    this.sourceManager = new SourceManager();
    this.sourceCommands = new SourceCommands(this.sourceManager);
    
    this.chunkManager = new ChunkManager();
    this.ingestionPipeline = new IngestionPipeline();
    this.retriever = new Retriever(this.chunkManager);
    this.llm = createLLMClient();
    
    this.groundingMode = 'strict'; // Default mode
  }

  public async start(): Promise<void> {
    const chatbotUrl = getChatbotApiUrl();
    console.log('🤖 Knowledge + Planning Bot - MVP');
    if (chatbotUrl) {
      console.log(`Connected to chatbot backend: ${chatbotUrl}`);
    }
    console.log('Type commands starting with / or ask questions directly');
    console.log('Available commands: /project, /stage, /task, /context, /source, /mode, /help');
    console.log('Type /help for more information\n');

    this.rl.setPrompt('> ');
    this.rl.prompt();

    this.rl.on('line', async (input) => {
      const trimmedInput = input.trim();
      
      if (trimmedInput.toLowerCase() === '/exit' || trimmedInput.toLowerCase() === '/quit') {
        console.log('Goodbye!');
        this.rl.close();
        return;
      }

      // When CHATBOT_API_URL is set, forward all input to the Go chatbot backend
      if (chatbotUrl) {
        try {
          const { reply, exit, session } = await sendToChatbot(chatbotUrl, trimmedInput, this.chatSession);
          this.chatSession = session;
          if (reply) {
            console.log(reply);
          }
          if (exit) {
            console.log('Goodbye!');
            this.rl.close();
            return;
          }
        } catch (err) {
          console.error('Chatbot API error:', err instanceof Error ? err.message : err);
        }
        this.rl.prompt();
        return;
      }

      if (trimmedInput.startsWith('/')) {
        await this.handleCommand(trimmedInput);
      } else {
        // Automatically detect if user wants to edit or ask a question
        await this.handleInput(trimmedInput);
      }

      this.rl.prompt();
    }).on('close', () => {
      process.exit(0);
    });
  }

  private async handleRedactCommand(args: string[]): Promise<void> {
    if (args.length < 2) {
      console.log('Usage: /redact <source_id> <pattern_to_redact>');
      console.log('Example: /redact 1 "sensitive information"');
      return;
    }

    const sourceId = args[0];
    const redactPattern = args.slice(1).join(' ');

    try {
      // Get all chunks for this source
      const chunks = this.chunkManager.getChunksBySource(sourceId);
      if (chunks.length === 0) {
        // Try to parse as number to provide better feedback
        const number = parseInt(sourceId, 10);
        const isNumber = !isNaN(number);
        const idDisplay = isNumber ? `Number: ${number}` : `ID: ${sourceId}`;
        console.log(`No chunks found for source ${idDisplay}. Please ingest the source first.`);
        return;
      }

      // Perform redaction on each chunk
      let redactedCount = 0;
      const regex = new RegExp(redactPattern, 'gi');
      
      for (const chunk of chunks) {
        const originalText = chunk.text;
        const redactedText = originalText.replace(regex, '[REDACTED]');
        if (originalText !== redactedText) {
          chunk.text = redactedText;
          redactedCount++;
        }
      }

      console.log(`Redaction completed. ${redactedCount} chunks were modified.`);
      console.log(`Pattern "${redactPattern}" has been replaced with [REDACTED] in source ID: ${sourceId}`);
    } catch (error) {
      console.error('Error during redaction:', error);
    }
  }

  private async handleInput(input: string): Promise<void> {
    // Use more sophisticated intent detection
    const intent = this.detectIntent(input);
    
    if (intent === 'edit') {
      // Treat as an editing request
      await this.handleEditRequest(input);
    } else {
      // Treat as a question
      await this.handleQuestion(input);
    }
  }
  
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
      'place in'
    ];
    
    return editIndicators.some(indicator => inputLower.includes(indicator));
  }

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
      console.log(`   Active sources: ${activeSources.map(s => `${s.number} | ${s.id.substring(0, 4)} | ${s.title}`).join(', ')}`);
    }
  }
  
  private isCompleteRewriteInstruction(instruction: string): boolean {
    const deleteKeywords = ['delete', 'remove', 'erase', 'eliminate', 'get rid of', 'clear out'];
    const rewriteKeywords = ['rewrite', 'reformat', 'restructure', 'update completely'];
    
    const lowerInstruction = instruction.toLowerCase();
    return deleteKeywords.some(keyword => lowerInstruction.includes(keyword)) ||
           rewriteKeywords.some(keyword => lowerInstruction.includes(keyword));
  }
  
  private isCompleteDocumentContent(editedContent: string, originalContent: string): boolean {
    // Check if the edited content contains structural elements that suggest a complete document
    const hasMultipleSections = (editedContent.match(/\n\s*\w/ig) || []).length > 3;
    const hasListMarkers = (editedContent.match(/[-*•]\s|^\d+\./gm) || []).length > 1;
    const hasHeaders = (editedContent.match(/^\s*#+\s|\b[A-Z][A-Z\s]*:\s*$/gm) || []).length > 0;
    
    // If it has document-like structure and is reasonably complete compared to original
    return (hasMultipleSections || hasListMarkers || hasHeaders) &&
           editedContent.length > originalContent.length * 0.5;
  }
  
  private attemptFuzzyReplacement(originalContent: string, relevantChunks: Chunk[], editedContent: string): string | null {
    // Try to find the best matching section to replace
    for (const chunk of relevantChunks) {
      // Try different levels of matching
      const matches = this.findBestMatchingSection(originalContent, chunk.text);
      if (matches.length > 0) {
        // Replace the best matching section
        const bestMatch = matches[0];
        return originalContent.substring(0, bestMatch.start) +
               editedContent +
               originalContent.substring(bestMatch.end);
      }
    }
    return null;
  }
  
  private findBestMatchingSection(content: string, target: string): Array<{start: number, end: number, score: number}> {
    const targetLines = target.split('\n').filter(line => line.trim());
    const contentLines = content.split('\n');
    
    const matches: Array<{start: number, end: number, score: number}> = [];
    
    for (let i = 0; i < contentLines.length; i++) {
      for (const targetLine of targetLines) {
        if (contentLines[i].includes(targetLine.trim())) {
          // Find surrounding context to get a larger section
          const contextStart = Math.max(0, i - 2);
          const contextEnd = Math.min(contentLines.length, i + 3);
          const contextSection = contentLines.slice(contextStart, contextEnd).join('\n');
          
          // Calculate similarity score
          const score = this.calculateSimilarity(contextSection, target);
          matches.push({
            start: content.indexOf(contextSection),
            end: content.indexOf(contextSection) + contextSection.length,
            score
          });
        }
      }
    }
    
    // Sort by score (highest first)
    return matches.sort((a, b) => b.score - a.score);
  }
  
  private calculateSimilarity(text1: string, text2: string): number {
    const words1 = new Set(text1.toLowerCase().match(/\w+/g) || []);
    const words2 = new Set(text2.toLowerCase().match(/\w+/g) || []);
    
    const intersection = new Set([...words1].filter(x => words2.has(x)));
    const union = new Set([...words1, ...words2]);
    
    return union.size > 0 ? intersection.size / union.size : 0;
  }

  private async handleCommand(input: string): Promise<void> {
    const [command, ...args] = input.split(' ');
    const commandName = command.substring(1).toLowerCase(); // Remove the '/' and lowercase

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

  private async handleEditCommand(input: string): Promise<void> {
    // For MVP, we'll handle simple edit commands
    // Format: /edit "instruction" <<< text
    // Or: /edit instruction (then ask for text in next input)
    
    const editMatch = input.match(/^\/edit\s+"([^"]+)"\s+<<<\s*([\s\S]*)/);
    
    if (editMatch) {
      const instruction = editMatch[1];
      const text = editMatch[2].trim();
      
      if (!text) {
        console.log('Please provide text to edit after the instruction');
        return;
      }
      
      await this.processEdit(instruction, text);
    } else {
      // For the simple case: /edit instruction, then we ask for text
      const instruction = input.substring(6).trim(); // Remove '/edit '
      
      if (!instruction) {
        console.log('Usage: /edit "instruction" <<< text');
        console.log('Example: /edit "Make this shorter" <<< This is a long text that should be made shorter');
        return;
      }
      
      // For MVP, we'll ask for text in the next input
      console.log(`Instruction: ${instruction}`);
      console.log('Please provide the text to edit in the next input (or use /edit "instruction" <<< text format)');
    }
  }

  private async processEdit(instruction: string, text: string): Promise<void> {
    try {
      // Get relevant chunks based on the text or instruction
      let relevantChunks: Chunk[] = [];
      if (this.groundingMode !== 'general') {
        const activeSourceIds = this.sourceManager.getActiveSources().map(s => s.id);
        const retrievalResult = this.retriever.retrieve(instruction + ' ' + text, activeSourceIds);
        relevantChunks = retrievalResult.chunks;
      }

      const context = this.contextManager.getContext();
      const result = await this.llm.edit(text, instruction, relevantChunks, context);
      
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
    } catch (error) {
      console.error('Error during edit:', error);
    }
  }

  private async handleQuestion(question: string): Promise<void> {
    try {
      let relevantChunks: Chunk[] = [];
      
      if (this.groundingMode !== 'general') {
        // Get active sources and retrieve relevant chunks
        const activeSourceIds = this.sourceManager.getActiveSources().map(s => s.id);
        const retrievalResult = this.retriever.retrieve(question, activeSourceIds);
        relevantChunks = retrievalResult.chunks;
        
        // In strict mode, only return if no chunks found
        if (relevantChunks.length === 0 && this.groundingMode === 'strict') {
          console.log("In strict mode: No relevant information found in the provided sources to answer this question.");
          return;
        }
      }

      const context = this.contextManager.getContext();
      const answer = await this.llm.answer(question, relevantChunks, context);
      
      console.log('\n┌─────────────────────────────────────────────────────────┐');
      console.log('│                     🤖 AI RESPONSE                     │');
      console.log('├─────────────────────────────────────────────────────────┤');
      
      // Split the answer into lines and handle each line
      const allLines = answer.split('\n');
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
    } catch (error) {
      console.error('Error processing question:', error);
    }
  }

  private handleModeCommand(args: string[]): void {
    if (args.length === 0) {
      console.log(`Current mode: ${this.groundingMode}`);
      console.log('Available modes: strict, hybrid, general');
      return;
    }

    const mode = args[0].toLowerCase();
    if (mode === 'strict' || mode === 'hybrid' || mode === 'general') {
      this.groundingMode = mode;
      console.log(`Mode set to: ${mode}`);
    } else {
      console.log(`Invalid mode: ${mode}. Use strict, hybrid, or general`);
    }
  }

  private async handleIngestCommand(): Promise<void> {
    const activeSources = this.sourceManager.getActiveSources();
    
    if (activeSources.length === 0) {
      console.log('No active sources to ingest. Use /source use <id> to activate sources first.');
      return;
    }

    console.log(`Processing ${activeSources.length} active source(s)...`);
    // Add a list of the sources being processed
    if (activeSources.length > 0) {
      const sourceList = activeSources.map(source => `  - ${source.number} | ${source.title}`).join('\n');
      console.log(`Sources to be processed:\n${sourceList}\n`);
    }
    
    for (const source of activeSources) {
      console.log(`\nProcessing source: ${source.number} | ${source.title}`);
      
      // Clear old chunks from this source before adding new ones
      this.chunkManager.removeChunksBySource(source.id);
      
      // Clear old chunks from this source before adding new ones
      this.chunkManager.removeChunksBySource(source.id);
      
      const result = await this.ingestionPipeline.processSource(source);
      
      if (result.status === 'success') {
        this.chunkManager.addChunks(result.chunks);
        console.log(`✓ Successfully processed ${result.chunks.length} chunks from ${source.number} | ${source.title}`);
      } else {
        console.log(`✗ Error processing ${source.number} | ${source.title}: ${result.error}`);
      }
    }
    
    console.log('\nIngestion completed!');
  }

  private showHelp(): void {
    console.log('\n📚 Available Commands:');
    console.log('/help - Show this help message');
    console.log('/project set <name> - Set the current project');
    console.log('/stage set <name> - Set the current stage');
    console.log('/task set <title> - Set the current task');
    console.log('/context show - Show current context');
    console.log('/context reset - Reset current context');
    console.log('/source add file <path> - Add a file source');
    console.log('/source add url <url> - Add a URL source');
    console.log('/source list - List all sources');
    console.log('/source use <id1,id2> - Set active sources (supports numbers: 1,2 or IDs: abc1,def2)');
    console.log('/source use all - Set all sources as active');
    console.log('/source clear - Clear active sources');
    console.log('/source remove <id> - Completely remove a source from the system');
    console.log('/ingest - Process active sources into chunks');
    console.log('/redact <source_id> <pattern> - Redact sensitive information from a source');
    console.log('/mode strict|hybrid|general - Set grounding mode');
    console.log('/exit or /quit - Exit the bot');
    console.log('');
    console.log('💡 Tips:');
    console.log('- Add sources first with /source add');
    console.log('- Use /source use to activate sources for retrieval');
    console.log('- Run /ingest to process sources into searchable chunks');
    console.log('- Use /source remove <id> to completely remove a source from the system');
    console.log('- Sources can be referenced by number (1, 2, 3) or partial ID (abc1, def2)');
    console.log('- Ask questions directly (without /) for Q&A mode');
    console.log('');
  }
}