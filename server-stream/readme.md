# gRPC Server Streaming Example

A simple example demonstrating server-side streaming in gRPC, where a client
requests to stream sensor data from the server.

## Running the Example
WIP

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
