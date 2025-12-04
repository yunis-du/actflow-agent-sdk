//! Calculator Agent Example
//!
//! This example demonstrates error handling and input validation.
//!
//! Run with:
//! ```bash
//! cargo run --example calculator_agent
//! ```

use actflow_agent_sdk::{Agent, AgentOutput, AgentServer, Context, Inputs, LogSender};

/// A calculator agent that performs basic arithmetic operations.
struct CalculatorAgent;

impl Agent for CalculatorAgent {
    async fn run(
        &self,
        _nid: String,
        _ctx: Context,
        inputs: Inputs,
        log: LogSender,
    ) -> AgentOutput {
        log.send("Calculator agent started").await;

        // Extract operation and operands from inputs
        let op = match inputs.get("operation").and_then(|v| v.as_str()) {
            Some(op) => op,
            None => {
                return AgentOutput::failed("Missing 'operation' field");
            }
        };

        let a = match inputs.get("a").and_then(|v| v.as_f64()) {
            Some(a) => a,
            None => {
                return AgentOutput::failed("Missing or invalid 'a' field");
            }
        };

        let b = match inputs.get("b").and_then(|v| v.as_f64()) {
            Some(b) => b,
            None => {
                return AgentOutput::failed("Missing or invalid 'b' field");
            }
        };

        log.send(format!("Calculating: {} {} {}", a, op, b)).await;

        // Perform the operation
        let result = match op {
            "add" | "+" => a + b,
            "sub" | "-" => a - b,
            "mul" | "*" => a * b,
            "div" | "/" => {
                if b == 0.0 {
                    return AgentOutput::exception("Division by zero");
                }
                a / b
            }
            _ => {
                return AgentOutput::failed(format!("Unknown operation: {}", op));
            }
        };

        log.send(format!("Result: {}", result)).await;

        AgentOutput::success(serde_json::json!({
            "result": result,
            "expression": format!("{} {} {} = {}", a, op, b, result)
        }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50052";
    println!("Starting CalculatorAgent server on {}", addr);

    AgentServer::new(CalculatorAgent).serve(addr).await?;

    Ok(())
}
