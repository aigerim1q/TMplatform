import * as fs from 'fs';
import * as path from 'path';
import { Chunk } from './types';

const CHUNKS_FILE = path.join('.bot', 'chunks.json');

export interface ChunkStorage {
  chunks: Chunk[];
}

export class ChunkManager {
  private storage: ChunkStorage;

  constructor() {
    this.storage = { chunks: [] };
    this.loadChunks();
  }

  private ensureBotDir() {
    if (!fs.existsSync('.bot')) {
      fs.mkdirSync('.bot', { recursive: true });
    }
  }

  private loadChunks(): void {
    this.ensureBotDir();
    
    if (fs.existsSync(CHUNKS_FILE)) {
      try {
        const data = JSON.parse(fs.readFileSync(CHUNKS_FILE, 'utf8'));
        this.storage = data;
      } catch (error) {
        console.error('Error loading chunks:', error);
        // Use default state if loading fails
      }
    } else {
      // Initialize with empty state
      this.saveChunks();
    }
  }

  public saveChunks(): void {
    this.ensureBotDir();
    fs.writeFileSync(CHUNKS_FILE, JSON.stringify(this.storage, null, 2));
  }

  public addChunks(chunks: Chunk[]): void {
    this.storage.chunks.push(...chunks);
    this.saveChunks();
  }

  public getChunksBySource(sourceId: string): Chunk[] {
    return this.storage.chunks.filter(chunk => chunk.sourceId === sourceId);
  }

  public getAllChunks(): Chunk[] {
    return this.storage.chunks;
  }

  public removeChunksBySource(sourceId: string): void {
    this.storage.chunks = this.storage.chunks.filter(chunk => chunk.sourceId !== sourceId);
    this.saveChunks();
  }
  
  public clearChunks(): void {
    this.storage.chunks = [];
    this.saveChunks();
  }

  public getChunksByIds(ids: string[]): Chunk[] {
    return this.storage.chunks.filter(chunk => ids.includes(chunk.chunkId));
  }
}