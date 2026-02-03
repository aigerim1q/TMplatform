package sources

import (
	"fmt"
	"strconv"
	"strings"
)

type SourceCommands struct {
	manager *SourceManager
}

func NewSourceCommands(manager *SourceManager) *SourceCommands {
	return &SourceCommands{
		manager: manager,
	}
}

func (sc *SourceCommands) HandleSourceCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /source add|list|use|clear|remove"
	}

	command := args[0]

	switch command {
	case "add":
		return sc.handleAddCommand(args[1:])
	case "list":
		return sc.handleListCommand()
	case "use":
		return sc.handleUseCommand(args[1:])
	case "clear":
		return sc.handleClearCommand()
	case "remove":
		return sc.handleRemoveCommand(args[1:])
	default:
		return fmt.Sprintf("Unknown source command: %s. Use add, list, use, clear, or remove", command)
	}
}

func (sc *SourceCommands) handleAddCommand(args []string) string {
	if len(args) < 2 {
		return "Usage: /source add file <path> or /source add url <url>"
	}

	sourceType := args[0]
	path := args[1]

	if sourceType != "file" && sourceType != "url" {
		return "Source type must be 'file' or 'url'"
	}

	source := sc.manager.AddSource(sourceType, path)
	return fmt.Sprintf("Added source: [Inactive] | id : %d | document name : %s", source.Number, source.Title)
}

func (sc *SourceCommands) handleListCommand() string {
	sources := sc.manager.ListSources()
	if len(sources) == 0 {
		return "No sources added yet."
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Total sources: %d\n", len(sources)))
	
	activeIds := make(map[string]bool)
	for _, id := range sc.manager.GetActiveSourceIds() {
		activeIds[id] = true
	}

	for _, source := range sources {
		status := "Inactive"
		if activeIds[source.ID] {
			status = "Active"
		}
		result.WriteString(fmt.Sprintf("[%s] | id : %d | document name : %s\n",
			status, source.Number, source.Title))
	}

	return result.String()
}

func (sc *SourceCommands) handleUseCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /source use <id1,id2> or /source use all or /source use clear"
	}

	if args[0] == "all" {
		allSources := sc.manager.GetAllSources()
		var sourceIds []string
		for _, source := range allSources {
			sourceIds = append(sourceIds, source.ID)
		}
		sc.manager.SetActiveSources(sourceIds)
		return fmt.Sprintf("Set all %d sources as active", len(allSources))
	}

	if args[0] == "clear" || args[0] == "none" {
		sc.manager.ClearActiveSources()
		return "Cleared all active sources"
	}

	// Parse comma-separated list of source IDs or numbers
	var sourceIds []string
	idStr := strings.Join(args, ",")
	ids := strings.Split(idStr, ",")

	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		// Try to find the source by number first
		if num, err := strconv.Atoi(id); err == nil {
			source := sc.manager.GetSourceByNumber(num)
			if source != nil {
				sourceIds = append(sourceIds, source.ID)
			} else {
				return fmt.Sprintf("No source found with number: %s", id)
			}
		} else {
			// Otherwise treat as ID (exact or partial)
			source := sc.manager.GetSourceByID(id)
			if source != nil {
				sourceIds = append(sourceIds, source.ID)
			} else {
				return fmt.Sprintf("No source found with ID: %s", id)
			}
		}
	}

	sc.manager.SetActiveSources(sourceIds)
	return fmt.Sprintf("Set %d sources as active", len(sourceIds))
}

func (sc *SourceCommands) handleClearCommand() string {
	sc.manager.ClearActiveSources()
	return "Cleared all active sources"
}

func (sc *SourceCommands) handleRemoveCommand(args []string) string {
	if len(args) == 0 {
		return "Usage: /source remove <id>"
	}

	id := args[0]
	success := sc.manager.RemoveSource(id)
	if success {
		return fmt.Sprintf("Removed source: %s", id)
	} else {
		return fmt.Sprintf("Failed to remove source: %s", id)
	}
}