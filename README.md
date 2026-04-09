# gRPC Hardcore: Concepts, Benchmarks & Implementation

This repository is dedicated to mastering the core concepts of gRPC through hands-on implementation and rigorous performance analysis. It explores the architectural advantages of gRPC over traditional REST, focusing on HTTP/2 multiplexing, binary serialization, and streaming capabilities.

## Projects

### 1. `pcbook` (Golang)
A deep-dive into building a microservice with Go and gRPC. This project covers:
- **Protobuf Design**: Defining complex message structures for PC components (CPU, RAM, GPU, Storage).
- **Serialization**: Comparing Protobuf binary format with JSON.
- **Service Implementation**: Building servers and clients for PC configuration.
- **Testing**: Comprehensive unit tests for serialization and logic.

### 2. `todo_app` (Node.js)
A performance-focused experimental project comparing **gRPC vs. REST (Express)** under high concurrency and simulated network latency.

#### Performance Benchmark Results
The following tests were conducted with a simulated **20ms network latency** to reflect real-world conditions.

| Concurrency Level | REST (RPS) | gRPC (RPS) | Performance Gain |
| :--- | :--- | :--- | :--- |
| **50 Parallel** | 1,003 | 1,136 | +13% |
| **100 Parallel** | 2,071 | 2,196 | +6% |
| **1,000 Parallel** | 2,423 | **5,358** | **2.2x Faster** |

**Key Finding: The "Concurrency Wall"**
REST (HTTP/1.1) hits a performance ceiling around 2,400 RPS due to the overhead of managing thousands of simultaneous TCP connections. gRPC (HTTP/2) uses **Multiplexing** to handle 1,000+ requests over a single TCP connection, effectively doubling throughput under high load.

## Tech Stack
- **Languages**: Golang, JavaScript (Node.js)
- **Protocols**: gRPC, HTTP/2, REST (HTTP/1.1)
- **Data Format**: Protocol Buffers (Protobuf), JSON
- **Frameworks**: Express.js (for REST benchmarking)

## Core Concepts Explored
- **Binary vs. Text**: Why Protobuf's binary format is more efficient than JSON for payload transmission.
- **HTTP/1.1 vs. HTTP/2**: Understanding Head-of-Line blocking in REST and Multiplexing in gRPC.
- **Unary vs. Streaming**: Implementation of simple request/response and high-throughput server-side streaming.

## Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) (for `pcbook`)
- [Node.js](https://nodejs.org/en/download/) (for `todo_app`)
- [protoc](https://grpc.io/docs/protoc-installation/) (Protocol Buffer Compiler)

### Running `pcbook`
```bash
cd pcbook
make protoc   # Generate Go code from proto files
make run      # Run the main application
make test     # Run unit tests
```

### Running `todo_app`
```bash
cd todo_app
npm install
node server.js       # Start gRPC Server
node restServer.js   # Start REST Server
node client.js       # Run benchmarks
```

---
*Created for learning and benchmarking gRPC core capabilities.*

## License
This project is licensed under the [MIT License](LICENSE).

