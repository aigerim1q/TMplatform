export interface KnowledgeSource {
  id: string;
  number: number; // Sequential number for user-friendly access
  type: 'file' | 'url';
  path: string; // file path or URL
  title: string;
  createdAt: Date;
  metadata?: {
    size?: number;
    format?: string;
    title?: string;
    description?: string;
  };
}

export interface SourceState {
  allSources: KnowledgeSource[];
  activeSourceIds: string[];
  nextSourceNumber: number; // For sequential numbering
}