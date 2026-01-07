import { SourceManager } from './manager';

export class SourceCommands {
  constructor(private sourceManager: SourceManager) {}

  public handleSourceCommand(args: string[]): string {
    if (args.length < 1) {
      return 'Usage: /source add file <path> | /source add url <url> | /source list | /source use <id1,id2> | /source clear';
    }

    const action = args[0];
    const params = args.slice(1);

    switch (action) {
      case 'add':
        return this.handleAddCommand(params);
      case 'list':
        return this.handleListCommand();
      case 'use':
        return this.handleUseCommand(params);
      case 'clear':
        return this.handleClearCommand();
      case 'remove':
        return this.handleRemoveCommand(params);
      default:
        return 'Unknown source command. Use: add, list, use, clear, or remove';
    }
  }

  private handleAddCommand(args: string[]): string {
    if (args.length < 2) {
      return 'Usage: /source add file <path> | /source add url <url>';
    }

    const type = args[0];
    const path = args[1];

    if (type !== 'file' && type !== 'url') {
      return 'Source type must be "file" or "url"';
    }

    try {
      const source = this.sourceManager.addSource(type as 'file' | 'url', path);
      return `Added source: ${source.title} (ID: ${source.number} | ${source.id.substring(0, 4)})`;
    } catch (error) {
      return `Error adding source: ${error instanceof Error ? error.message : String(error)}`;
    }
  }

  private handleListCommand(): string {
    const sources = this.sourceManager.listSources();
    if (sources.length === 0) {
      return 'No sources added yet.';
    }

    const activeSources = this.sourceManager.getActiveSources();
    const activeIds = new Set(activeSources.map(s => s.id));

    const sourceList = sources.map(source => {
      const status = activeIds.has(source.id) ? '[ACTIVE]' : '[INACTIVE]';
      return `${status} ID: ${source.number} | ${source.id.substring(0, 4)} | Type: ${source.type} | ${source.path}`;
    }).join('\n');

    return `Sources:\n${sourceList}`;
  }

  private handleUseCommand(args: string[]): string {
    if (args.length === 0) {
      return 'Usage: /source use <id1,id2,...> or /source use all';
    }

    if (args[0] === 'all') {
      const allSources = this.sourceManager.listSources();
      const allIds = allSources.map(s => s.id);
      this.sourceManager.setActiveSources(allIds);
      return `All ${allSources.length} sources set as active`;
    }

    const ids = args[0].split(',').map(id => id.trim());
    // Validate that all IDs exist (allowing partial ID matching)
    const allSources = this.sourceManager.listSources();
    const validIds = [];
    const invalidIds = [];
    
    for (const id of ids) {
      // First try to match with the full ID, then with partial match, then with sequential number
      let matchingSource: any = null;
      
      // Exact UUID match
      const exactMatch = allSources.find(s => s.id === id);
      if (exactMatch) {
        matchingSource = exactMatch;
      }
      
      // Partial UUID match
      if (!matchingSource) {
        const partialMatch = allSources.find(s => s.id.startsWith(id));
        if (partialMatch) {
          matchingSource = partialMatch;
        }
      }
      
      // Sequential number match
      if (!matchingSource) {
        const number = parseInt(id, 10);
        if (!isNaN(number)) {
          const numberMatch = allSources.find(s => s.number === number);
          if (numberMatch) {
            matchingSource = numberMatch;
          }
        }
      }
      
      if (matchingSource) {
        validIds.push(matchingSource.id);
      } else {
        invalidIds.push(id);
      }
    }

    if (invalidIds.length > 0) {
      return `Invalid source IDs: ${invalidIds.join(', ')}`;
    }

    this.sourceManager.setActiveSources(validIds);
    return `Set ${validIds.length} sources as active`;
  }

  private handleRemoveCommand(args: string[]): string {
    if (args.length < 1) {
      return 'Usage: /source remove <id>';
    }

    const sourceId = args[0];
    const result = this.sourceManager.removeSource(sourceId);
    
    if (result) {
      // Try to parse as number to provide better feedback
      const number = parseInt(sourceId, 10);
      const isNumber = !isNaN(number);
      const idDisplay = isNumber ? `Number: ${number}` : `ID: ${sourceId}`;
      return `Source with ${idDisplay} has been completely removed from the system.`;
    } else {
      return `Source with ID ${sourceId} not found or could not be removed.`;
    }
  }

  private handleClearCommand(): string {
    this.sourceManager.clearActiveSources();
    return 'Cleared all active sources';
  }
}