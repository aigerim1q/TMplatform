export interface Chunk {
  chunkId: string;
  sourceId: string;
  text: string;
  metadata: {
    position: number;
    length: number;
  };
}

export interface ChunkingOptions {
  chunkSize: number;
  overlap: number;
}

export interface IngestionResult {
  chunks: Chunk[];
  sourceId: string;
  status: 'success' | 'error';
  error?: string;
}