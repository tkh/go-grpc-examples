# gRPC Client Streaming Example

A simple example demonstrating client-side streaming in gRPC, where a client
streams a file to the server in chunks.

## Running the Example

1. First, start the server in a terminal:
```bash
# View available options
go run cmd/server/main.go -help

# Start the server with default settings
go run cmd/server/main.go
```

2. In another terminal, run the client:
```bash
# View available options
go run cmd/client/main.go -help

# Run with default settings (uploads a 100MB test file)
go run cmd/client/main.go
```

The server will store uploaded files in the configured uploads directory (default: `uploads`).

## Proto Generation

This project uses [buf](https://buf.build) and
[Connect](https://connectrpc.com) for protocol buffer management. Before
running the proto commands, please follow the [Connect Getting Started
guide](https://connectrpc.com/docs/go/getting-started) to set up the necessary
tools on your machine.

Once set up, the following commands are available:

```bash
# Format proto files
buf format proto -w

# Lint proto files
buf lint proto

# Generate Go code from proto files
buf generate proto
```
