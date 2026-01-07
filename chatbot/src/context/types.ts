export interface Context {
  project: string | null;
  stage: string | null;
  task: string | null;
  goals: string[];
  constraints: string[];
  decisions: string[];
  next_steps: string[];
}

export interface ContextState {
  active: Context;
  history: Context[];
}