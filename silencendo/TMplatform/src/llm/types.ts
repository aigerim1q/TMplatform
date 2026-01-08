import { Chunk } from '../ingestion/types';
import { Context } from '../context/types';

export interface LLMMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
}

export interface LLMClient {
  generate(messages: LLMMessage[]): Promise<string>;
  answer(question: string, retrievedChunks: Chunk[], context: Context): Promise<string>;
  edit(inputText: string, instruction: string, retrievedChunks: Chunk[], context: Context): Promise<string>;
}

export interface GenerationOptions {
  temperature?: number;
  maxTokens?: number;
}