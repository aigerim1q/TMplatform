package context

import (
	"fmt"
	"strconv"
	"strings"
)

type ContextCommands struct {
	manager *ContextManager
}

func NewContextCommands(manager *ContextManager) *ContextCommands {
	return &ContextCommands{
		manager: manager,
	}
}

func (cc *ContextCommands) HandleProjectCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /project set <name>"
	}
	
	if args[0] == "set" && len(args) > 1 {
		name := strings.Join(args[1:], " ")
		cc.manager.SetProject(name)
		return fmt.Sprintf("Project set to: %s", name)
	}
	
	return "Usage: /project set <name>"
}

func (cc *ContextCommands) HandleStageCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /stage set <name>"
	}
	
	if args[0] == "set" && len(args) > 1 {
		name := strings.Join(args[1:], " ")
		cc.manager.SetStage(name)
		return fmt.Sprintf("Stage set to: %s", name)
	}
	
	return "Usage: /stage set <name>"
}

func (cc *ContextCommands) HandleTaskCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /task set <title>"
	}
	
	if args[0] == "set" && len(args) > 1 {
		title := strings.Join(args[1:], " ")
		cc.manager.SetTask(title)
		return fmt.Sprintf("Task set to: %s", title)
	}
	
	return "Usage: /task set <title>"
}

func (cc *ContextCommands) HandleContextCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /context show | reset"
	}
	
	if args[0] == "show" {
		return cc.manager.GetContextSummary()
	} else if args[0] == "reset" {
		cc.manager.ResetContext()
		return "Context reset successfully"
	}
	
	return "Usage: /context show | reset"
}