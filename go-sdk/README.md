# Actflow Agent Go SDK

A Go SDK for building agents for the Actflow workflow engine.

## Installation

```bash
go get github.com/yunis-du/actflow-agent-sdk/go-sdk@v0.1.1
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    sdk "github.com/yunis-du/actflow-agent-sdk/go-sdk"
)

type MyAgent struct{}

func (a *MyAgent) Run(ctx context.Context, nid string, execCtx *sdk.Context, inputs any, logCh chan<- string) (*sdk.Output, error) {
    logCh <- "Starting agent..."
    // Your agent logic here
    return sdk.Success(map[string]any{"result": "done"}), nil
}

func (a *MyAgent) Shutdown(ctx context.Context) error {
    log.Println("Shutting down...")
    return nil
}

func main() {
    sdk.NewServer(&MyAgent{}).Serve("0.0.0.0:50051")
}
```

## API Reference

### Agent Interface

```go
type Agent interface {
    Run(ctx context.Context, nid string, execCtx *Context, inputs any, log chan<- string) (*Output, error)
    Shutdown(ctx context.Context) error
}
```

### Context

```go
type Context struct {
    PID  string            // Process ID
    Env  map[string]string // Environment variables
    Outputs map[string]any    // Workflow outputs
}
```

### Output Helpers

| Function | Description |
|----------|-------------|
| `Success(outputs any)` | Create a successful output |
| `Failed(err string)` | Create a failed output |
| `FailedErr(err error)` | Create a failed output from error |
| `Exception(msg string)` | Create an exception output |

### ExecutionStatus

| Status | Value |
|--------|-------|
| `StatusPending` | 0 |
| `StatusSucceeded` | 1 |
| `StatusFailed` | 2 |
| `StatusException` | 3 |
| `StatusStopped` | 4 |
| `StatusPaused` | 5 |

## Examples

### Simple Echo Agent

```bash
cd examples/simple && go run main.go
```

### Calculator Agent (Error Handling)

```bash
cd examples/calculator && go run main.go
```

### Stateful Agent (Shutdown Handling)

```bash
cd examples/stateful && go run main.go
```

## License

MIT
