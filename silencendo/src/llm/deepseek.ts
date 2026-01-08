import { LLMClient, LLMMessage } from './types';
import { Chunk } from '../ingestion/types';
import { Context } from '../context/types';
import OpenAI from 'openai';
import 'dotenv/config';

// Type alias to match OpenAI's expected format
type OpenAIMessage = {
  role: 'user' | 'assistant' | 'system';
  content: string;
};

export class DeepSeekLLM implements LLMClient {
  private client: OpenAI;
  private model: string;

  constructor() {
    const apiKey = process.env.DEEPSEEK_API_KEY;
    if (!apiKey) {
      throw new Error(
        'DEEPSEEK_API_KEY environment variable is required. ' +
        'Please set it using: export DEEPSEEK_API_KEY=your_key_here'
      );
    }

    const baseUrl = process.env.DEEPSEEK_BASE_URL || 'https://api.deepseek.com';
    this.model = process.env.DEEPSEEK_MODEL || 'deepseek-chat';

    this.client = new OpenAI({
      apiKey,
      baseURL: baseUrl,
    });
  }

  async generate(messages: LLMMessage[]): Promise<string> {
    try {
      // Convert our LLMMessage type to OpenAI's expected format
      const openAiMessages: OpenAIMessage[] = messages.map(msg => ({
        role: msg.role,
        content: msg.content
      }));

      const response = await this.client.chat.completions.create({
        model: this.model,
        messages: openAiMessages,
        temperature: 0.7,
      });

      return response.choices[0]?.message?.content || '';
    } catch (error) {
      throw new Error(`DeepSeek API error: ${(error as Error).message}`);
    }
  }

  async answer(question: string, retrievedChunks: Chunk[], context: Context): Promise<string> {
    // If no chunks are found, let DeepSeek answer the question based on its general knowledge
    // but provide context about the project
    if (retrievedChunks.length === 0) {
      const systemMessage = `You are an AI assistant. The user asked a question but no specific context was provided from their documents.

PROJECT CONTEXT:
Project: ${context.project || 'Not set'}
Stage: ${context.stage || 'Not set'}
Task: ${context.task || 'Not set'}

If the question is related to the project context above, provide a relevant response.
If the question is general and not related to the project, you may use your general knowledge to answer.
If the question is specific to the user's documents and you have no context, state that you don't have the relevant information from the documents.`;

      const userMessage = `Question: ${question}

Please provide an appropriate response based on the project context and your general knowledge if applicable.`;

      try {
        const response = await this.client.chat.completions.create({
          model: this.model,
          messages: [
            { role: 'system', content: systemMessage },
            { role: 'user', content: userMessage }
          ],
          temperature: 0.7,
        });

        return response.choices[0]?.message?.content || 'No response generated.';
      } catch (error) {
        return `DeepSeek API error: ${(error as Error).message}`;
      }
    }

    // Format the chunks for the system prompt when chunks are available
    const formattedChunks = retrievedChunks
      .map((chunk, index) => `[SOURCE ${index + 1}]: ${chunk.text}`)
      .join('\n\n');

    const systemMessage = `You are an AI assistant that answers questions based ONLY on the provided CONTEXT CHUNKS.
If the answer is not found in the provided chunks, explicitly state that you cannot answer from the provided sources.
Do not fabricate information or rely on general knowledge when in strict mode.

CONTEXT CHUNKS:
${formattedChunks}

PROJECT CONTEXT:
Project: ${context.project || 'Not set'}
Stage: ${context.stage || 'Not set'}
Task: ${context.task || 'Not set'}

Follow these rules:
1. Answer ONLY using information from the CONTEXT CHUNKS
2. If information is not in the chunks, say you cannot answer from sources
3. Reference the source number when citing information
4. Be concise and accurate`;

    const userMessage = `Question: ${question}

Please provide a comprehensive answer based on the provided context chunks. If the information is not available in the provided sources, state that you cannot answer from the provided sources.`;

    try {
      const response = await this.client.chat.completions.create({
        model: this.model,
        messages: [
          { role: 'system', content: systemMessage },
          { role: 'user', content: userMessage }
        ],
        temperature: 0.3,
      });

      return response.choices[0]?.message?.content || 'No response generated.';
    } catch (error) {
      return `DeepSeek API error: ${(error as Error).message}`;
    }
  }

  async edit(inputText: string, instruction: string, retrievedChunks: Chunk[], context: Context): Promise<string> {
    // Format the chunks for context
    const formattedChunks = retrievedChunks
      .map((chunk, index) => `[SOURCE ${index + 1}]: ${chunk.text}`)
      .join('\n\n');

    const systemMessage = `You are an AI editor that modifies text based on specific instructions.
Use the provided CONTEXT CHUNKS as reference material when applicable to the editing task.
Apply the given instruction to transform the input text while maintaining accuracy and coherence.

CONTEXT CHUNKS:
${formattedChunks}

PROJECT CONTEXT:
Project: ${context.project || 'Not set'}
Stage: ${context.stage || 'Not set'}  
Task: ${context.task || 'Not set'}

Editing Rules:
1. Follow the instruction precisely
2. Use context chunks as reference when relevant to the edit
3. Preserve the original meaning unless instructed otherwise
4. Make minimal changes unless specifically asked for significant modifications`;

    const userMessage = `Original Text: ${inputText}

Instruction: ${instruction}

Please apply the instruction to the original text. Return only the edited text without additional commentary.`;

    try {
      const response = await this.client.chat.completions.create({
        model: this.model,
        messages: [
          { role: 'system', content: systemMessage },
          { role: 'user', content: userMessage }
        ],
        temperature: 0.5,
      });

      return response.choices[0]?.message?.content || 'No edit generated.';
    } catch (error) {
      return `DeepSeek API error: ${(error as Error).message}`;
    }
  }
}