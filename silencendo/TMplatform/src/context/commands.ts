import { ContextManager } from './manager';

export class ContextCommands {
  constructor(private contextManager: ContextManager) {}

  public handleProjectCommand(args: string[]): string {
    if (args.length >= 2 && args[0] === 'set') {
      const projectName = args.slice(1).join(' ');
      this.contextManager.setProject(projectName);
      return `Project set to: ${projectName}`;
    }
    return 'Usage: /project set <name>';
  }

  public handleStageCommand(args: string[]): string {
    if (args.length >= 2 && args[0] === 'set') {
      const stageName = args.slice(1).join(' ');
      this.contextManager.setStage(stageName);
      return `Stage set to: ${stageName}`;
    }
    return 'Usage: /stage set <name>';
  }

  public handleTaskCommand(args: string[]): string {
    if (args.length >= 2 && args[0] === 'set') {
      const taskTitle = args.slice(1).join(' ');
      this.contextManager.setTask(taskTitle);
      return `Task set to: ${taskTitle}`;
    }
    return 'Usage: /task set <title>';
  }

  public handleContextCommand(args: string[]): string {
    if (args.length === 1 && args[0] === 'show') {
      return this.contextManager.getContextSummary();
    } else if (args.length === 1 && args[0] === 'reset') {
      this.contextManager.resetContext();
      return 'Context has been reset';
    }
    return 'Usage: /context show | /context reset';
  }
}