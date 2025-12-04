// Simple Agent Example
//
// This example demonstrates how to create a simple echo agent.
//
// Run with:
//
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/yunis-du/actflow-agent-sdk/go-sdk"
)

// EchoAgent is a simple agent that echoes the inputs.
type EchoAgent struct{}

// Run executes the agent logic.
func (a *EchoAgent) Run(ctx context.Context, nid string, execCtx *sdk.Context, inputs any, logCh chan<- string) (*sdk.Output, error) {
	logCh <- fmt.Sprintf("Node %s starting execution", nid)
	logCh <- fmt.Sprintf("Process ID: %s", execCtx.PID)

	// Log environment variables
	if len(execCtx.Env) > 0 {
		logCh <- fmt.Sprintf("Environment: %v", execCtx.Env)
	}

	// Log workflow variables
	if len(execCtx.Vars) > 0 {
		logCh <- fmt.Sprintf("Variables: %v", execCtx.Vars)
	}

	// Process the inputs
	logCh <- fmt.Sprintf("Received inputs: %v", inputs)

	// Return success
	return sdk.Success(map[string]any{
		"node_id": nid,
		"echo":    inputs,
		"message": "Hello from EchoAgent!",
	}), nil
}

// Shutdown is called when the agent should shut down.
func (a *EchoAgent) Shutdown(ctx context.Context) error {
	log.Println("EchoAgent shutting down...")
	return nil
}

func main() {
	addr := "[::1]:50051"
	log.Printf("Starting EchoAgent server on %s", addr)

	server := sdk.NewServer(&EchoAgent{})
	if err := server.Serve(addr); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
