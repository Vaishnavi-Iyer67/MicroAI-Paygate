# Gateway Service

The Gateway is the high-performance entry point for the MicroAI Paygate architecture. Written in Go, it handles traffic orchestration, payment enforcement, and proxying to AI providers.

## Role & Responsibilities

- **Traffic Entry Point**: Listens on port 3000 and accepts all incoming API requests.
- **x402 Enforcement**: Inspects headers for `X-402-Signature` and `X-402-Nonce`. If missing, it rejects the request with a 402 status and payment context.
- **Verification Orchestration**: Communicates with the internal Rust Verifier service to validate cryptographic signatures.
- **Proxying**: Forwards authenticated requests to the OpenRouter API and returns the response to the client.

## Technology Stack

- **Language**: Go (Golang) 1.24
- **Framework**: Gin Web Framework
- **Concurrency**: Goroutines for non-blocking I/O operations.

## Key Files

- `main.go`: Contains the server initialization, route definitions, and the core `handleSummarize` logic.
- `Dockerfile`: Multi-stage build configuration for creating a lightweight Alpine Linux container.

## Development

To run the gateway locally:

```bash
go run main.go
```

Ensure the Verifier service is running on port 3002 before starting the Gateway.

## Configuration

Environment variables (via `.env`):

- `OPENROUTER_API_KEY` (required) — API key for OpenRouter
- `OPENROUTER_MODEL` (optional) — model name, default `z-ai/glm-4.5-air:free`
- `VERIFIER_URL` (optional) — override verifier endpoint, default `http://127.0.0.1:3002/verify`
- `RECIPIENT_ADDRESS` (optional) — payment recipient; falls back to hardcoded value if unset
- `CHAIN_ID` (optional) — chain id used in EIP-712 domain; default `8453`

Ports: Gateway listens on `3000` by default. Update configs if you must avoid conflicts.

## Testing

```bash
go test -v
```
