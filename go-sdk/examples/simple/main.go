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
func (a *EchoAgent) Run(ctx context.Context, pid, nid string, inputs map[string]any, logCh chan<- string) (*sdk.Output, error) {
	logCh <- fmt.Sprintf("Node %s starting execution", nid)
	logCh <- fmt.Sprintf("Process ID: %s", pid)

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
	addr := "0.0.0.0:50051"
	log.Printf("Starting EchoAgent server on %s", addr)

	server := sdk.NewServer(&EchoAgent{})
	if err := server.Serve(addr); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
