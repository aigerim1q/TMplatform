import { Context, ContextState } from './types';
import * as fs from 'fs';
import * as path from 'path';

const CONTEXT_FILE = path.join('.bot', 'context.json');

export class ContextManager {
  private state: ContextState;

  constructor() {
    this.state = {
      active: {
        project: null,
        stage: null,
        task: null,
        goals: [],
        constraints: [],
        decisions: [],
        next_steps: [],
      },
      history: [],
    };
    this.loadContext();
  }

  private ensureBotDir() {
    if (!fs.existsSync('.bot')) {
      fs.mkdirSync('.bot', { recursive: true });
    }
  }

  private loadContext(): void {
    this.ensureBotDir();
    
    if (fs.existsSync(CONTEXT_FILE)) {
      try {
        const data = JSON.parse(fs.readFileSync(CONTEXT_FILE, 'utf8'));
        this.state = data;
      } catch (error) {
        console.error('Error loading context:', error);
        // Use default state if loading fails
      }
    } else {
      // Initialize with default context
      this.saveContext();
    }
  }

  public saveContext(): void {
    this.ensureBotDir();
    fs.writeFileSync(CONTEXT_FILE, JSON.stringify(this.state, null, 2));
  }

  public getContext(): Context {
    return this.state.active;
  }

  public setProject(name: string): void {
    this.state.active.project = name;
    this.saveContext();
  }

  public setStage(name: string): void {
    this.state.active.stage = name;
    this.saveContext();
  }

  public setTask(title: string): void {
    this.state.active.task = title;
    this.saveContext();
  }

  public addGoal(goal: string): void {
    this.state.active.goals.push(goal);
    this.saveContext();
  }

  public addConstraint(constraint: string): void {
    this.state.active.constraints.push(constraint);
    this.saveContext();
  }

  public addDecision(decision: string): void {
    this.state.active.decisions.push(decision);
    this.saveContext();
  }

  public addNextStep(step: string): void {
    this.state.active.next_steps.push(step);
    this.saveContext();
  }

  public resetContext(): void {
    this.state.history.push({ ...this.state.active });
    this.state.active = {
      project: null,
      stage: null,
      task: null,
      goals: [],
      constraints: [],
      decisions: [],
      next_steps: [],
    };
    this.saveContext();
  }

  public getContextSummary(): string {
    const { project, stage, task, goals, constraints, decisions, next_steps } = this.state.active;
    
    const summaryParts = [];
    
    if (project) summaryParts.push(`Project: ${project}`);
    if (stage) summaryParts.push(`Stage: ${stage}`);
    if (task) summaryParts.push(`Task: ${task}`);
    
    if (goals.length > 0) summaryParts.push(`Goals: ${goals.length} items`);
    if (constraints.length > 0) summaryParts.push(`Constraints: ${constraints.length} items`);
    if (decisions.length > 0) summaryParts.push(`Decisions: ${decisions.length} items`);
    if (next_steps.length > 0) summaryParts.push(`Next Steps: ${next_steps.length} items`);
    
    return summaryParts.join('\n');
  }
}