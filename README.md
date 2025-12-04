# Actflow Agent SDK

SDKs for building agents for the [Actflow](https://github.com/yunis-du/actflow) workflow engine.

## Overview

Actflow Agent SDK provides libraries for multiple programming languages to help you build custom agents that integrate with the Actflow workflow engine. Agents are executed as gRPC services that receive workflow context and inputs, perform custom logic, and return results.

## Available SDKs

| Language | Directory | Status |
|----------|-----------|--------|
| Rust | [rust-sdk](./rust-sdk) | ✅ Ready |
| Go | [go-sdk](./go-sdk) | ✅ Ready |
| Python | - | 🚧 In Progress |

## Architecture

```
┌─────────────────┐         gRPC          ┌─────────────────┐
│                 │ ◄──────────────────── │                 │
│  Actflow Engine │      RunRequest       │   Your Agent    │
│                 │ ────────────────────► │  (SDK-based)    │
│                 │   Stream<AgentUpdate> │                 │
└─────────────────┘                       └─────────────────┘
```

### Agent Lifecycle

1. **Run**: Engine sends `RunRequest` with node ID, context, and inputs
2. **Execute**: Agent processes the request and streams log messages
3. **Complete**: Agent returns `AgentOutput` with status and results
4. **Shutdown**: Engine can request graceful shutdown

## Quick Start

### Rust

```rust
use actflow_agent_sdk::{Agent, AgentServer, AgentOutput, Context, Inputs, LogSender};

struct MyAgent;

impl Agent for MyAgent {
    async fn run(&self, nid: String, ctx: Context, inputs: Inputs, log: LogSender) -> AgentOutput {
        log.send("Processing...").await;
        AgentOutput::success(serde_json::json!({"result": "done"}))
    }
}

#[tokio::main]
async fn main() {
    AgentServer::new(MyAgent).serve("[::1]:50051").await.unwrap();
}
```

### Go

```go
type MyAgent struct{}

func (a *MyAgent) Run(ctx context.Context, nid string, execCtx *sdk.Context, inputs any, logCh chan<- string) (*sdk.Output, error) {
    logCh <- "Processing..."
    return sdk.Success(map[string]any{"result": "done"}), nil
}

func (a *MyAgent) Shutdown(ctx context.Context) error { return nil }

func main() {
    sdk.NewServer(&MyAgent{}).Serve("[::1]:50051")
}
```

## Proto Definition

The SDKs are built on a common gRPC protocol defined in the Actflow repository:

```protobuf
service AgentService {
  rpc Run(RunRequest) returns (stream AgentUpdate);
  rpc Shutdown(Empty) returns (Empty);
}
```

The proto file is included via git submodule from [actflow](https://github.com/yunis-du/actflow).

## Project Structure

```
actflow-agent-sdk/
├── rust-sdk/           # Rust SDK
│   ├── src/
│   ├── examples/
│   └── Cargo.toml
├── go-sdk/             # Go SDK
│   ├── pb/             # Generated protobuf code
│   ├── examples/
│   └── go.mod
└── external/           # Git submodule (actflow proto)
    └── actflow/
```

## Development

### Prerequisites

- Rust 1.85+ (for rust-sdk)
- Go 1.24+ (for go-sdk)
- protoc 3.0+ (for regenerating proto code)


### Running Examples

```bash
# Rust examples
cd rust-sdk
cargo run --example simple_agent
cargo run --example calculator_agent

# Go examples
cd go-sdk/examples/simple && go run main.go
cd go-sdk/examples/calculator && go run main.go
```

## License

MIT
