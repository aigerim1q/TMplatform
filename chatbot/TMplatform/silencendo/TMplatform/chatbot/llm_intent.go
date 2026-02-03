package chatbot

import (
	"context"
	"strings"

	router "silencendo/intent_router"
	"silencendo/services"
)

// llmIntentProcessor uses the LLM as the central reasoning engine
type llmIntentProcessor struct {
	router *router.Router
}

// NewLLMIntentProcessor creates a new processor that uses the LLM router
func NewLLMIntentProcessor(router *router.Router) *llmIntentProcessor {
	return &llmIntentProcessor{
		router: router,
	}
}

// DetectIntentWithLLM processes the message using the LLM as the central reasoning engine
func (p *llmIntentProcessor) DetectIntentWithLLM(ctx context.Context, message string) IntentMatch {
	msg := strings.TrimSpace(message)

	// Process the input through the LLM router
	llmResponse, err := p.router.ProcessInput(ctx, msg)
	if err != nil {
		// If LLM processing fails completely, return unknown with low confidence
		return IntentMatch{
			Intent:     IntentUnknown,
			Confidence: 0.1,
		}
	}

	// Convert the LLM response to an IntentMatch using the local conversion
	localMatch := p.router.ConvertToIntentMatch(llmResponse, msg)

	// Convert local match to chatbot.IntentMatch
	intentMatch := IntentMatch{
		Intent:       localMatch.Intent,
		ProjectTitle: localMatch.ProjectTitle,
		Description:  localMatch.Description,
		EntityType:   localMatch.EntityType,
		EntityName:   localMatch.EntityName,
		AssigneeName: localMatch.AssigneeName,
		Members:      localMatch.Members,
		Confidence:   localMatch.Confidence,
		// Note: StagePlans and Filters need special handling if needed
		StagePlans: nil,                    // This will need to be populated from parameters if needed
		Filters:    services.ListFilters{}, // This will need to be populated from parameters if needed
	}

	return intentMatch
}

// DetectIntentWithLLM is a helper function that creates an LLM-based intent detector
func DetectIntentWithLLM(ctx context.Context, router *router.Router, message string) IntentMatch {
	processor := NewLLMIntentProcessor(router)
	return processor.DetectIntentWithLLM(ctx, message)
}
