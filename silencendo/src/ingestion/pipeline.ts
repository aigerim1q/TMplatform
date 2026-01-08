import * as fs from 'fs';
import * as path from 'path';
import * as url from 'url';
import pdfParse from 'pdf-parse';
import * as cheerio from 'cheerio';
import { KnowledgeSource } from '../sources/types';
import { Chunk, ChunkingOptions, IngestionResult } from './types';
import { randomUUID } from 'crypto';

export class IngestionPipeline {
  private chunkingOptions: ChunkingOptions;

  constructor(options?: Partial<ChunkingOptions>) {
    this.chunkingOptions = {
      chunkSize: options?.chunkSize || 1000,
      overlap: options?.overlap || 150,
    };
  }

  public async processSource(source: KnowledgeSource): Promise<IngestionResult> {
    try {
      let text: string;

      switch (source.type) {
        case 'file':
          text = await this.extractTextFromFile(source.path);
          break;
        case 'url':
          text = await this.extractTextFromUrl(source.path);
          break;
        default:
          throw new Error(`Unsupported source type: ${source.type}`);
      }

      const chunks = this.chunkText(text, source.id);
      return {
        chunks,
        sourceId: source.id,
        status: 'success',
      };
    } catch (error) {
      return {
        chunks: [],
        sourceId: source.id,
        status: 'error',
        error: error instanceof Error ? error.message : String(error),
      };
    }
  }

  private async extractTextFromFile(filePath: string): Promise<string> {
    const resolvedPath = path.resolve(filePath);
    
    if (!fs.existsSync(resolvedPath)) {
      throw new Error(`File does not exist: ${resolvedPath}`);
    }

    const ext = path.extname(resolvedPath).toLowerCase();
    
    switch (ext) {
      case '.txt':
      case '.md':
        return fs.readFileSync(resolvedPath, 'utf-8');
      
      case '.pdf':
        const pdfData = await pdfParse(fs.readFileSync(resolvedPath));
        return pdfData.text;
      
      default:
        throw new Error(`Unsupported file type: ${ext}. Supported types: .txt, .md, .pdf`);
    }
  }

  private async extractTextFromUrl(url: string): Promise<string> {
    // For MVP, we'll simulate web scraping
    // In a real implementation, you would use a library like node-fetch or axios
    try {
      // Importing dynamically to avoid issues if not available
      const { default: fetch } = await import('node-fetch');
      const response = await fetch(url);
      
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      
      const html = await response.text();
      const $ = cheerio.load(html);
      
      // Remove script and style elements
      $('script, style, nav, footer, header').remove();
      
      // Extract text content
      return $('body').text().replace(/\s+/g, ' ').trim();
    } catch (error) {
      // If fetch is not available or fails, return a placeholder
      console.warn(`Could not fetch URL: ${url}. Using placeholder text.`);
      return `Content from ${url} (URL content could not be retrieved)`;
    }
  }

  private chunkText(text: string, sourceId: string): Chunk[] {
    const chunks: Chunk[] = [];
    const { chunkSize, overlap } = this.chunkingOptions;
    
    if (text.length <= chunkSize) {
      chunks.push({
        chunkId: randomUUID(),
        sourceId,
        text,
        metadata: {
          position: 0,
          length: text.length,
        },
      });
      return chunks;
    }

    // Split text into chunks with overlap
    let position = 0;
    while (position < text.length) {
      const chunkText = text.substring(position, position + chunkSize);
      chunks.push({
        chunkId: randomUUID(),
        sourceId,
        text: chunkText,
        metadata: {
          position,
          length: chunkText.length,
        },
      });

      // Move position by chunkSize minus overlap
      position += chunkSize - overlap;
      
      // If remaining text is less than or equal to overlap, we're done
      if (position >= text.length) {
        break;
      }
    }

    return chunks;
  }
}