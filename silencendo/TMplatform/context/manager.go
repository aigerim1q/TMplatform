package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const CONTEXT_FILE = ".bot/context.json"

type ContextManager struct {
	state ContextState
}

func NewContextManager() *ContextManager {
	cm := &ContextManager{
		state: ContextState{
			Active: ContextType{
				Project:     nil,
				Stage:       nil,
				Task:        nil,
				Goals:       []string{},
				Constraints: []string{},
				Decisions:   []string{},
				NextSteps:   []string{},
			},
			History: []ContextType{},
		},
	}
	cm.loadContext()
	return cm
}

func (cm *ContextManager) ensureBotDir() {
	if _, err := os.Stat(".bot"); os.IsNotExist(err) {
		os.MkdirAll(".bot", 0755)
	}
}

func (cm *ContextManager) loadContext() {
	cm.ensureBotDir()
	
	if _, err := os.Stat(CONTEXT_FILE); err == nil {
		data, err := os.ReadFile(CONTEXT_FILE)
		if err == nil {
			err = json.Unmarshal(data, &cm.state)
			if err != nil {
				// Use default state if loading fails
				cm.saveContext()
			}
		} else {
			// Use default state if loading fails
			cm.saveContext()
		}
	} else {
		// Initialize with default context
		cm.saveContext()
	}
}

func (cm *ContextManager) saveContext() {
	cm.ensureBotDir()
	data, err := json.MarshalIndent(cm.state, "", "  ")
	if err == nil {
		os.WriteFile(CONTEXT_FILE, data, 0644)
	}
}

func (cm *ContextManager) GetContext() ContextType {
	return cm.state.Active
}

func (cm *ContextManager) SetProject(name string) {
	cm.state.Active.Project = &name
	cm.saveContext()
}

func (cm *ContextManager) SetStage(name string) {
	cm.state.Active.Stage = &name
	cm.saveContext()
}

func (cm *ContextManager) SetTask(title string) {
	cm.state.Active.Task = &title
	cm.saveContext()
}

func (cm *ContextManager) AddGoal(goal string) {
	cm.state.Active.Goals = append(cm.state.Active.Goals, goal)
	cm.saveContext()
}

func (cm *ContextManager) AddConstraint(constraint string) {
	cm.state.Active.Constraints = append(cm.state.Active.Constraints, constraint)
	cm.saveContext()
}

func (cm *ContextManager) AddDecision(decision string) {
	cm.state.Active.Decisions = append(cm.state.Active.Decisions, decision)
	cm.saveContext()
}

func (cm *ContextManager) AddNextStep(step string) {
	cm.state.Active.NextSteps = append(cm.state.Active.NextSteps, step)
	cm.saveContext()
}

func (cm *ContextManager) ResetContext() {
	// Add current context to history
	cm.state.History = append(cm.state.History, cm.state.Active)
	
	// Reset to default context
	cm.state.Active = ContextType{
		Project:     nil,
		Stage:       nil,
		Task:        nil,
		Goals:       []string{},
		Constraints: []string{},
		Decisions:   []string{},
		NextSteps:   []string{},
	}
	cm.saveContext()
}

func (cm *ContextManager) GetContextSummary() string {
	project := ""
	stage := ""
	task := ""
	
	if cm.state.Active.Project != nil {
		project = *cm.state.Active.Project
	}
	if cm.state.Active.Stage != nil {
		stage = *cm.state.Active.Stage
	}
	if cm.state.Active.Task != nil {
		task = *cm.state.Active.Task
	}
	
	var summaryParts []string
	
	if project != "" {
		summaryParts = append(summaryParts, "Project: "+project)
	}
	if stage != "" {
		summaryParts = append(summaryParts, "Stage: "+stage)
	}
	if task != "" {
		summaryParts = append(summaryParts, "Task: "+task)
	}
	
	if len(cm.state.Active.Goals) > 0 {
		summaryParts = append(summaryParts, "Goals: "+strconv.Itoa(len(cm.state.Active.Goals))+" items")
	}
	if len(cm.state.Active.Constraints) > 0 {
		summaryParts = append(summaryParts, "Constraints: "+strconv.Itoa(len(cm.state.Active.Constraints))+" items")
	}
	if len(cm.state.Active.Decisions) > 0 {
		summaryParts = append(summaryParts, "Decisions: "+strconv.Itoa(len(cm.state.Active.Decisions))+" items")
	}
	if len(cm.state.Active.NextSteps) > 0 {
		summaryParts = append(summaryParts, "Next Steps: "+strconv.Itoa(len(cm.state.Active.NextSteps))+" items")
	}
	
	return strings.Join(summaryParts, "\n")
}