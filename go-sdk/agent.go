// Package sdk provides a Go SDK for building Actflow agents.
package sdk

import "context"

// Agent is the interface that users must implement to create an agent.
type Agent interface {
	// Run executes the agent logic.
	// The log channel is used to send log messages during execution.
	// Close the log channel or return to complete execution.
	Run(ctx context.Context, pid, nid string, inputs map[string]any, log chan<- string) (*Output, error)

	// Shutdown is called when the agent should shut down.
	// Use this to perform cleanup operations.
	Shutdown(ctx context.Context) error
}

// ExecutionStatus represents the status of agent execution.
type ExecutionStatus int32

const (
	StatusPending   ExecutionStatus = 0
	StatusSucceeded ExecutionStatus = 1
	StatusFailed    ExecutionStatus = 2
	StatusException ExecutionStatus = 3
	StatusStopped   ExecutionStatus = 4
	StatusPaused    ExecutionStatus = 5
)

// String returns the string representation of the status.
func (s ExecutionStatus) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusSucceeded:
		return "SUCCEEDED"
	case StatusFailed:
		return "FAILED"
	case StatusException:
		return "EXCEPTION"
	case StatusStopped:
		return "STOPPED"
	case StatusPaused:
		return "PAUSED"
	default:
		return "UNKNOWN"
	}
}

// Output represents the result of agent execution.
type Output struct {
	// Status is the execution status.
	Status ExecutionStatus
	// Outputs contains the output data.
	Outputs any
	// Error message (when Status is Failed).
	Error string
	// Exception message (when Status is Exception).
	Exception string
}

// Success creates a successful output.
func Success(outputs any) *Output {
	return &Output{
		Status:  StatusSucceeded,
		Outputs: outputs,
	}
}

// Failed creates a failed output with an error message.
func Failed(err string) *Output {
	return &Output{
		Status: StatusFailed,
		Error:  err,
	}
}

// FailedErr creates a failed output from an error.
func FailedErr(err error) *Output {
	return &Output{
		Status: StatusFailed,
		Error:  err.Error(),
	}
}

// Exception creates an exception output.
func Exception(exception string) *Output {
	return &Output{
		Status:    StatusException,
		Exception: exception,
	}
}
