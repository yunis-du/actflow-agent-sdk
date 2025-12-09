//! Simple Agent Example
//!
//! This example demonstrates how to create a simple agent using the Actflow Agent SDK.
//!
//! Run with:
//! ```bash
//! cargo run --example simple_agent
//! ```

use actflow_agent_sdk::{Agent, AgentOutput, AgentServer, Inputs, LogSender};

/// A simple echo agent that returns the inputs with some processing.
struct EchoAgent;

impl Agent for EchoAgent {
    async fn run(&self, pid: String, nid: String, inputs: Inputs, log: LogSender) -> AgentOutput {
        // Log the start of execution
        log.send(format!("Node {} starting execution", nid)).await;
        log.send(format!("Process ID: {}", pid)).await;

        // Process the inputs
        log.send(format!("Received inputs: {}", inputs)).await;

        // Simulate some work
        tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
        log.send("Processing complete").await;

        // Return success with the processed result
        AgentOutput::success(serde_json::json!({
            "node_id": nid,
            "echo": inputs,
            "message": "Hello from EchoAgent!"
        }))
    }

    async fn shutdown(&self) {
        println!("EchoAgent shutting down...");
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50051";
    println!("Starting EchoAgent server on {}", addr);

    AgentServer::new(EchoAgent).serve(addr).await?;

    Ok(())
}
