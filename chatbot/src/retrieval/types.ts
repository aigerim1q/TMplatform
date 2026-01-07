import { Chunk } from '../ingestion/types';

export interface RetrievalResult {
  chunks: Chunk[];
  scores: number[];
}

export interface RetrievalOptions {
  topK: number;
  minScore: number;
}