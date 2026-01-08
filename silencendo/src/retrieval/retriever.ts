import { Chunk } from '../ingestion/types';
import { RetrievalResult, RetrievalOptions } from './types';
import { ChunkManager } from '../ingestion/chunkManager';

export class Retriever {
  private chunkManager: ChunkManager;
  private options: RetrievalOptions;

  constructor(chunkManager: ChunkManager, options?: Partial<RetrievalOptions>) {
    this.chunkManager = chunkManager;
    this.options = {
      topK: options?.topK || 5,
      minScore: options?.minScore || 0.01,
    };
  }

  public retrieve(query: string, sourceIds?: string[]): RetrievalResult {
    // Get all chunks or filter by source IDs
    let chunks: Chunk[];
    if (sourceIds && sourceIds.length > 0) {
      chunks = [];
      for (const sourceId of sourceIds) {
        chunks.push(...this.chunkManager.getChunksBySource(sourceId));
      }
    } else {
      chunks = this.chunkManager.getAllChunks();
    }

    if (chunks.length === 0) {
      return { chunks: [], scores: [] };
    }

    // Calculate similarity scores for the query against all chunks
    const scores = chunks.map(chunk => this.calculateSimilarity(query, chunk.text));
    
    // Create pairs of chunks and scores, then sort by score in descending order
    let chunkScorePairs = chunks.map((chunk, index) => ({
      chunk,
      score: scores[index]
    }));

    // Filter by minimum score, but if no chunks pass the filter and we have chunks,
    // return at least the highest scoring chunk to avoid empty results
    const filteredPairs = chunkScorePairs.filter(pair => pair.score >= this.options.minScore);
    
    if (filteredPairs.length === 0 && chunkScorePairs.length > 0) {
      // If no chunks meet the threshold, return the best scoring chunk anyway
      // This prevents empty results when the query is not similar to any chunk
      const bestPair = chunkScorePairs.reduce((max, current) =>
        current.score > max.score ? current : max
      );
      chunkScorePairs = [bestPair];
    } else {
      chunkScorePairs = filteredPairs;
    }

    // Sort by score and take topK
    chunkScorePairs = chunkScorePairs
      .sort((a, b) => b.score - a.score)
      .slice(0, this.options.topK);

    // Extract chunks and scores
    const resultChunks = chunkScorePairs.map(pair => pair.chunk);
    const resultScores = chunkScorePairs.map(pair => pair.score);

    return {
      chunks: resultChunks,
      scores: resultScores
    };
  }

  private calculateSimilarity(query: string, text: string): number {
    // Enhanced similarity calculation with multiple strategies
    const queryWords = this.tokenize(query.toLowerCase());
    const textWords = this.tokenize(text.toLowerCase());
    
    if (queryWords.length === 0 || textWords.length === 0) {
      return 0;
    }
    
    // Calculate Jaccard similarity (intersection over union)
    const queryWordSet = new Set(queryWords);
    const textWordSet = new Set(textWords);
    
    let overlap = 0;
    for (const word of queryWordSet) {
      if (textWordSet.has(word)) {
        overlap++;
      }
    }
    
    const jaccard = overlap / Math.max(queryWordSet.size, textWordSet.size, 1);
    
    // Also calculate overlap ratio relative to query size (for queries that are shorter than text)
    const queryOverlapRatio = overlap / Math.max(queryWordSet.size, 1);
    
    // Use the higher of the two scores to handle both long text vs short query and vice versa
    return Math.max(jaccard, queryOverlapRatio);
  }

  private tokenize(text: string): string[] {
    // Simple tokenization: split on non-word characters and filter empty strings
    return text.split(/\W+/).filter(token => token.length > 0);
  }
}