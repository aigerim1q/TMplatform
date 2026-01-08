import { KnowledgeSource, SourceState } from './types';
import * as fs from 'fs';
import * as path from 'path';
import { randomUUID } from 'crypto';

const SOURCES_FILE = path.join('.bot', 'sources.json');

export class SourceManager {
  private state: SourceState;

  constructor() {
    this.state = {
      allSources: [],
      activeSourceIds: [],
      nextSourceNumber: 1,
    };
    this.loadSources();
  }

  private ensureBotDir() {
    if (!fs.existsSync('.bot')) {
      fs.mkdirSync('.bot', { recursive: true });
    }
  }

  private loadSources(): void {
    this.ensureBotDir();
    
    if (fs.existsSync(SOURCES_FILE)) {
      try {
        const data = JSON.parse(fs.readFileSync(SOURCES_FILE, 'utf8'));
        // Convert date strings back to Date objects
        const loadedSources = data.allSources.map((source: any) => ({
          ...source,
          createdAt: new Date(source.createdAt)
        }));
        
        // Calculate nextSourceNumber based on highest number in loaded sources
        const maxNumber = loadedSources.length > 0
          ? Math.max(...loadedSources.map((s: any) => s.number))
          : 0;
        const nextNumber = Math.max(maxNumber + 1, data.nextSourceNumber || 1);
        
        this.state = {
          ...data,
          allSources: loadedSources,
          nextSourceNumber: nextNumber
        };
      } catch (error) {
        console.error('Error loading sources:', error);
        // Use default state if loading fails
      }
    } else {
      // Initialize with empty state
      this.saveSources();
    }
  }

  private saveSources(): void {
    this.ensureBotDir();
    // Convert Date objects to strings for JSON serialization
    const serializableState = {
      ...this.state,
      allSources: this.state.allSources.map(source => ({
        ...source,
        createdAt: source.createdAt.toISOString()
      })),
      nextSourceNumber: this.state.nextSourceNumber
    };
    fs.writeFileSync(SOURCES_FILE, JSON.stringify(serializableState, null, 2));
  }

  public addSource(type: 'file' | 'url', path: string): KnowledgeSource {
    // Check if source already exists
    const existingSource = this.state.allSources.find(s => s.path === path);
    if (existingSource) {
      return existingSource;
    }

    // Find the lowest available number
    let availableNumber = 1;
    const usedNumbers = new Set(this.state.allSources.map(s => s.number));
    while (usedNumbers.has(availableNumber)) {
      availableNumber++;
    }
    
    const newSource: KnowledgeSource = {
      id: randomUUID(),
      number: availableNumber,
      type,
      path,
      title: path.split('/').pop() || path,
      createdAt: new Date(),
      metadata: {}
    };
    
    // Update nextSourceNumber to be the next available number
    this.state.nextSourceNumber = Math.max(this.state.nextSourceNumber, availableNumber + 1);

    this.state.allSources.push(newSource);
    this.saveSources();
    return newSource;
  }

  public listSources(): KnowledgeSource[] {
    return this.state.allSources;
  }

  public setSourceActive(sourceId: string, active: boolean): void {
    if (active && !this.state.activeSourceIds.includes(sourceId)) {
      this.state.activeSourceIds.push(sourceId);
    } else if (!active) {
      this.state.activeSourceIds = this.state.activeSourceIds.filter(id => id !== sourceId);
    }
    this.saveSources();
  }

  public setActiveSources(sourceIds: string[]): void {
    this.state.activeSourceIds = sourceIds;
    this.saveSources();
  }

  public clearActiveSources(): void {
    this.state.activeSourceIds = [];
    this.saveSources();
  }

  public getActiveSources(): KnowledgeSource[] {
    return this.state.allSources.filter(source => 
      this.state.activeSourceIds.includes(source.id)
    );
  }

  public getSourceById(id: string): KnowledgeSource | undefined {
    // First try exact UUID match
    const exactMatch = this.state.allSources.find(source => source.id === id);
    if (exactMatch) {
      return exactMatch;
    }
    
    // Then try partial UUID match (for backward compatibility with 4-char prefixes)
    const partialMatch = this.state.allSources.find(source => source.id.startsWith(id));
    if (partialMatch) {
      return partialMatch;
    }
    
    // Then try sequential number match
    const number = parseInt(id, 10);
    if (!isNaN(number)) {
      const numberMatch = this.state.allSources.find(source => source.number === number);
      if (numberMatch) {
        return numberMatch;
      }
    }
    
    return undefined;
  }
  
  public getSourceByNumber(number: number): KnowledgeSource | undefined {
    return this.state.allSources.find(source => source.number === number);
  }
  
  public getSourceIdByNumber(number: number): string | undefined {
    const source = this.state.allSources.find(source => source.number === number);
    return source?.id;
  }

  public getAllSources(): KnowledgeSource[] {
    return this.state.allSources;
  }

  public updateSourceContent(sourceId: string, newContent: string): boolean {
    const source = this.state.allSources.find(s => s.id === sourceId);
    if (!source || source.type !== 'file') {
      return false;
    }

    try {
      // Write the new content back to the file
      fs.writeFileSync(source.path, newContent, 'utf8');
      
      // After updating the file, we should re-ingest it to update the chunks
      return true;
    } catch (error) {
      console.error('Error updating source content:', error);
      return false;
    }
  }

  public getSourceContent(sourceId: string): string | null {
    const source = this.state.allSources.find(s => s.id === sourceId);
    if (!source || source.type !== 'file') {
      return null;
    }

    try {
      return fs.readFileSync(source.path, 'utf8');
    } catch (error) {
      console.error('Error reading source content:', error);
      return null;
    }
  }

  public removeSource(sourceId: string): boolean {
    // Try to find the source by UUID, partial UUID, or sequential number
    let matchingSource: KnowledgeSource | undefined = undefined;
    
    // First try exact UUID match
    const exactMatch = this.state.allSources.find(s => s.id === sourceId);
    if (exactMatch) {
      matchingSource = exactMatch;
    }
    
    // Then try partial UUID match
    if (!matchingSource) {
      const partialMatch = this.state.allSources.find(s => s.id.startsWith(sourceId));
      if (partialMatch) {
        matchingSource = partialMatch;
      }
    }
    
    // Then try sequential number match
    if (!matchingSource) {
      const number = parseInt(sourceId, 10);
      if (!isNaN(number)) {
        const numberMatch = this.state.allSources.find(s => s.number === number);
        if (numberMatch) {
          matchingSource = numberMatch;
        }
      }
    }
    
    if (!matchingSource) {
      return false;
    }
    
    const initialLength = this.state.allSources.length;
    this.state.allSources = this.state.allSources.filter(source => source.id !== matchingSource!.id);
    
    // Also remove from active sources if it was active
    this.state.activeSourceIds = this.state.activeSourceIds.filter(id => id !== matchingSource!.id);
    
    if (this.state.allSources.length < initialLength) {
      this.saveSources();
      return true;
    }
    return false;
  }
}