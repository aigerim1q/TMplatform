package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

// UnifiedRepl owns the single stdin reader and delegates all input to the dispatcher.
type UnifiedRepl struct {
	dispatcher *Dispatcher
}

func NewUnifiedRepl(dispatcher *Dispatcher) *UnifiedRepl {
	return &UnifiedRepl{dispatcher: dispatcher}
}

// Start begins the single shared input loop using the provided reader (typically os.Stdin).
func (r *UnifiedRepl) Start(ctx context.Context, reader io.Reader) error {
	fmt.Println("🤖 Knowledge + Planning Bot - Unified CLI")
	fmt.Println("Type commands starting with / or ask questions directly")
	fmt.Println("Available commands: /project, /stage, /task, /context, /source, /mode, /help")
	fmt.Println("Type /help for more information")
	fmt.Println()

	scanner := bufio.NewScanner(reader)

	fmt.Print("> ")
	for scanner.Scan() {
		exit, err := r.dispatcher.Dispatch(ctx, scanner.Text())
		if err != nil {
			fmt.Printf("Error handling input: %v\n", err)
		}
		if exit {
			break
		}
		fmt.Print("> ")
	}

	return scanner.Err()
}
