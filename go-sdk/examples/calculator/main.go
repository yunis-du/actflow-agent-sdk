// Calculator Agent Example
//
// This example demonstrates error handling and input validation.
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

// CalculatorAgent performs basic arithmetic operations.
type CalculatorAgent struct{}

// Run executes the calculator logic.
func (a *CalculatorAgent) Run(ctx context.Context, nid string, execCtx *sdk.Context, inputs any, logCh chan<- string) (*sdk.Output, error) {
	logCh <- "Calculator agent started"

	// Parse inputs as map
	inputMap, ok := inputs.(map[string]any)
	if !ok {
		return sdk.Failed("Invalid inputs format"), nil
	}

	// Get operation
	op, ok := inputMap["operation"].(string)
	if !ok {
		return sdk.Failed("Missing 'operation' field"), nil
	}

	// Get operands
	aVal, ok := getNumber(inputMap, "a")
	if !ok {
		return sdk.Failed("Missing or invalid 'a' field"), nil
	}

	bVal, ok := getNumber(inputMap, "b")
	if !ok {
		return sdk.Failed("Missing or invalid 'b' field"), nil
	}

	logCh <- fmt.Sprintf("Calculating: %v %s %v", aVal, op, bVal)

	// Perform operation
	var result float64
	switch op {
	case "add", "+":
		result = aVal + bVal
	case "sub", "-":
		result = aVal - bVal
	case "mul", "*":
		result = aVal * bVal
	case "div", "/":
		if bVal == 0 {
			return sdk.Exception("Division by zero"), nil
		}
		result = aVal / bVal
	default:
		return sdk.Failed(fmt.Sprintf("Unknown operation: %s", op)), nil
	}

	logCh <- fmt.Sprintf("Result: %v", result)

	return sdk.Success(map[string]any{
		"result":     result,
		"expression": fmt.Sprintf("%v %s %v = %v", aVal, op, bVal, result),
	}), nil
}

// Shutdown is called when the agent should shut down.
func (a *CalculatorAgent) Shutdown(ctx context.Context) error {
	log.Println("CalculatorAgent shutting down...")
	return nil
}

// getNumber extracts a number from the map.
func getNumber(m map[string]any, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func main() {
	addr := "[::1]:50052"
	log.Printf("Starting CalculatorAgent server on %s", addr)

	server := sdk.NewServer(&CalculatorAgent{})
	if err := server.Serve(addr); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

