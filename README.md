# go-protocol-api-style

> **Note on API Styles and Protocols:**  
> This project demonstrates how multiple API styles (REST, GraphQL, SOAP) and a true protocol (gRPC) can be supported in a single Go codebase.  
> - **Protocols** (like HTTP, gRPC) define *how* data is transmitted over the network.  
> - **API styles/specifications** (like REST, SOAP, GraphQL) define *how* APIs structure requests and responses, but rely on a protocol (almost always HTTP or HTTP/2) for transport.  
>   - **REST**: API style using HTTP and JSON  
>   - **GraphQL**: Query language/specification using HTTP and JSON  
>   - **SOAP**: Messaging specification using HTTP and XML  
>   - **gRPC**: Protocol and framework using HTTP/2 and Protobuf  
>  
> In summary, **gRPC is a protocol**, while **REST, SOAP, and GraphQL are API styles/specifications** that use a protocol underneath.

---

A robust template project demonstrating **multi-protocol APIs in Go**: **REST**, **GraphQL**, **gRPC**, and **SOAP**.  
Built with [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io/), [gqlgen](https://github.com/99designs/gqlgen), [protobuf/gRPC](https://grpc.io/docs/languages/go/quickstart/), and native SOAP.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Setup & Initialization](#setup--initialization)
- [Configuration](#configuration)
- [Running the Server](#running-the-server)
- [Testing & Usage](#testing--usage)
- [Protocols Deep Dive](#protocols-deep-dive)
  - [REST](#rest)
  - [GraphQL](#graphql)
  - [gRPC](#grpc)
  - [SOAP](#soap)
- [Adding New APIs](#adding-new-apis)
- [Auto-Generation & Tools](#auto-generation--tools)
- [Comparison: API Styles vs Protocols](#comparison-api-styles-vs-protocols)
- [Why Install These Tools?](#why-install-these-tools)
- [What is WSDL and .proto?](#what-is-wsdl-and-proto)
- [Protocol Lifecycle Visualization](#protocol-lifecycle-visualization)
- [Protocols and HTTP Explained](#protocols-and-http-explained)

---

## Features

- **REST API:** Standard HTTP endpoints for CRUD and business logic.
- **GraphQL API:** Flexible and efficient querying for frontends and integrations.
- **gRPC API:** High-performance, type-safe remote procedure calls via Protocol Buffers.
- **SOAP API:** Enterprise-grade, legacy system compatibility.
- **Unified codebase:** All API styles share domain models and business logic.

---

## Project Structure

```
go-protocol-api-style/
├── config/                  # Configuration loader
├── graph/                   # GraphQL schema, models, resolvers
├── internal/
│   ├── domain/              # Domain models (shared)
│   ├── dto/                 # REST Data Transfer Objects
│   ├── infrastructure/
│   │   ├── database/        # DB and repository implementations
│   │   ├── graphql/         # GraphQL route handler
│   │   ├── grpc/            # gRPC service implementation and protos
│   │   │   └── moviepb/     # gRPC generated files
│   │   ├── http/            # REST route handler
│   │   └── soap/            # SOAP handler and WSDL
│   ├── usecase/             # Business logic/services
│   └── repository/          # Repository interfaces
├── logger/                  # Logger setup
├── migrations/              # DB migrations
├── main.go                  # Entry point (starts ALL protocols)
├── docker-compose.yml       # PostgreSQL setup for dev/test
└── README.md
```

---

## Requirements

- Go 1.20+
- Docker & Docker Compose (for running PostgreSQL)
- [curl](https://curl.se/) (API testing)
- [jq](https://stedolan.github.io/jq/) (optional - prettifying JSON)
- Protocol-specific tools and generators (see [Auto-Generation & Tools](#auto-generation--tools))

---

## Setup & Initialization

### 1. Clone the Repository

```sh
git clone https://github.com/sorrawichYooboon/go-protocol-api-style.git
cd go-protocol-api-style
```

### 2. Start PostgreSQL with Docker Compose

```sh
docker-compose up -d
```

- **DB Name:** `protocoldb`
- **User:** `protocol_user`
- **Password:** `protocol_password`
- **Port:** `5432`

### 3. Install Go Dependencies

```sh
go mod tidy
```

---

## Configuration

Set your environment variables (or use a `.env` file):

```
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=protocol_user
DATABASE_PASSWORD=protocol_password
DATABASE_DBNAME=protocoldb
DATABASE_SSLMODE=disable
```

You can export these in your shell or use [direnv](https://direnv.net/).

---

## Running the Server

```sh
go run main.go
```

The server runs on **port 8081** by default.  
Override with the `PORT` environment variable if needed.

gRPC server runs typically on **50051** (see grpc/server.go for details).

---

## Testing & Usage

### REST

#### List all movies

```sh
curl -s http://localhost:8081/movies | jq
```

#### Get a movie by ID

```sh
curl -s http://localhost:8081/movies/1 | jq
```

---

### GraphQL

#### Query all movies

```sh
curl -s -X POST http://localhost:8081/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ movies { id title description releaseDate } }"}' | jq
```

#### Query movie by ID

```sh
curl -s -X POST http://localhost:8081/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ movie(id: \"1\") { id title description releaseDate } }"}' | jq
```

#### Use Playground

Open [http://localhost:8081/playground](http://localhost:8081/playground) in your browser.

---

### gRPC

### Prerequisites

- Install [grpcurl](https://github.com/fullstorydev/grpcurl):

```sh
brew install grpcurl
```

- Install [protobuf compiler](https://grpc.io/docs/protoc-installation/):

```sh
brew install protobuf
```

- If your server does **not** support reflection, download or use your `.proto` files.

#### Install Go plugins and libraries

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

go get google.golang.org/grpc
go get google.golang.org/protobuf
```

#### gRPC: Reflection vs Importing Proto

gRPC clients like `grpcurl` can discover and interact with your gRPC server in **two ways**:  
**1. Using server reflection**  
**2. By importing proto files manually**

#### **Server Reflection**

**Reflection** is a feature where the gRPC server exposes metadata about its services and methods at runtime.

- **When enabled**, clients (like grpcurl) can dynamically discover service names, methods, and message types without needing the `.proto` files.
- This makes testing and debugging easier, since you can query the server directly.

**Example usage (with reflection enabled):**
```sh
grpcurl -plaintext localhost:50051 movie.MovieService/ListMovies | jq
grpcurl -plaintext -d '{"id": 1}' localhost:50051 movie.MovieService/GetMovie | jq
```
You do **not** need to specify `-import-path` or `-proto`—everything is discovered via reflection.

**When to use:**  
- Quick testing/debugging
- You do **not** have access to the `.proto` files
- The server supports reflection (not all do)

#### **Importing Proto Files**

If the gRPC server **does not support reflection**, you must provide the `.proto` files to grpcurl so it knows how to encode requests and decode responses.

- You specify the directory containing the proto files with `-import-path`
- You specify the proto file(s) themselves with `-proto`

**Example usage (without reflection):**
```sh
grpcurl -plaintext \
  -import-path proto \
  -proto proto/movie.proto \
  localhost:50051 movie.MovieService/ListMovies | jq

grpcurl -plaintext \
  -import-path proto \
  -proto proto/movie.proto \
  -d '{"id": 1}' localhost:50051 movie.MovieService/GetMovie | jq
```

**When to use:**  
- The server does **not** support reflection
- You have the `.proto` files locally
- For automated scripts or CI, where you want strict type definitions

---

### **Summary Table**

| Method           | Reflection Enabled | Need Proto Files | Usage                                      |
|------------------|-------------------|------------------|---------------------------------------------|
| Reflection       | Yes               | No               | `grpcurl -plaintext localhost:50051 ...`    |
| Import Proto     | No                | Yes              | `grpcurl -plaintext -import-path ...`       |

**Tip:**  
Most production gRPC servers disable reflection for security.  
For local development and testing, reflection is very convenient.

#### ListMovies example (reflection)

```sh
grpcurl -plaintext localhost:50051 movie.MovieService/ListMovies | jq
```

#### GetMovie by ID example (reflection)

```sh
grpcurl -plaintext -d '{"id": 1}' localhost:50051 movie.MovieService/GetMovie | jq
```

#### ListMovies example (import proto)

```sh
grpcurl -plaintext \
  -import-path proto \
  -proto proto/movie.proto \
  localhost:50051 movie.MovieService/ListMovies | jq
```

#### GetMovie by ID example (import proto)

```sh
grpcurl -plaintext \
  -import-path proto \
  -proto proto/movie.proto \
  -d '{"id": 1}' localhost:50051 movie.MovieService/GetMovie | jq
```

#### Generate gRPC Go code from proto

```sh
protoc --go_out=. --go-grpc_out=. proto/movie.proto
```

---

### SOAP

#### Get movie by ID

Create a file `get_movie.xml`:

```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
  <soapenv:Body>
    <GetMovieRequest>
      <id>1</id>
    </GetMovieRequest>
  </soapenv:Body>
</soapenv:Envelope>
```

Send request:

```sh
curl -s -X POST http://localhost:8081/soap/movie \
  -H "Content-Type: text/xml; charset=utf-8" \
  --data-binary @get_movie.xml
```

#### Get WSDL

```sh
curl -s http://localhost:8081/soap/movie.wsdl
```

---

## Protocols Deep Dive

### REST

**Type:** API architectural style (not a protocol).  
**Transport:** HTTP  
**Format:** JSON  
**Setup:**  
- Routes handled in `internal/infrastructure/http/`.
- Uses [Gin](https://github.com/gin-gonic/gin) for fast routing & middleware.

**How it works:**  
- Client sends HTTP request to an endpoint (e.g., `/movies/1`).
- Server parses URL/verb, decodes JSON, processes, returns JSON response.

---

### GraphQL

**Type:** Query language and specification (not a protocol).  
**Transport:** HTTP  
**Format:** JSON  
**Setup:**  
- Schema in `graph/schema.graphqls`.
- Resolvers in `graph/resolver.go`.
- Uses [gqlgen](https://github.com/99designs/gqlgen) for codegen.

**How it works:**  
- Client sends POST request with query/mutation.
- Server resolves requested fields via resolvers, returns JSON.

---

### gRPC

**Type:** Protocol and framework.  
**Transport:** HTTP/2  
**Format:** Protobuf (binary)  
**Setup:**  
- `.proto` files in `internal/infrastructure/grpc/moviepb/`
- Generated Go code (`movie.pb.go`, `movie_grpc.pb.go`)
- Server implementation in `internal/infrastructure/grpc/server.go`
- Uses [google.golang.org/grpc](https://grpc.io/docs/languages/go/quickstart/)

**How it works:**  
- Client creates and sends binary requests (protobuf) over HTTP/2.
- Server receives, deserializes, processes via service methods, returns protobuf response.

---

### SOAP

**Type:** Messaging specification (not a protocol).  
**Transport:** HTTP  
**Format:** XML  
**Setup:**  
- Handlers and WSDL in `internal/infrastructure/soap/`.
- Uses native Go XML + custom logic (or third-party packages).

**How it works:**  
- Client sends XML payload in `soap:Envelope` via HTTP POST.
- Server parses XML, processes request, returns XML response.

---

## Adding New APIs

### REST

1. Define DTO for request/response in `internal/dto/`.
2. Implement handler in `internal/infrastructure/http/handler.go`.
3. Register route in `router.go`.
4. Connect to usecase for business logic.

### GraphQL

1. Update schema in `graph/schema.graphqls`.
2. Run `gqlgen generate` for codegen.
3. Implement resolver in `graph/resolver.go`.
4. Connect to usecase/domain as needed.

### gRPC

1. Edit/create `.proto` file in `internal/infrastructure/grpc/moviepb/`.
2. Run `protoc` (see below) to re-generate Go code.
3. Implement new handler in `server.go`.
4. Wire to domain/usecase logic.

### SOAP

1. Define new XML request/response structs in `internal/infrastructure/soap/`.
2. Implement handler for SOAP action.
3. Update WSDL file as needed.
4. Wire to business logic as needed.

---

## Auto-Generation & Tools

### Why install these libraries and tools?

- **brew install grpcurl**: Installs grpcurl, a CLI for testing gRPC endpoints without writing client code.
- **brew install protobuf**: Installs `protoc`, the Protocol Buffers compiler, which converts `.proto` files into Go code.
- **go install google.golang.org/protobuf/cmd/protoc-gen-go@latest**: Installs the Go plugin for `protoc` to generate Go structs from `.proto` files.
- **go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest**: Installs the Go gRPC code generator for server/client code.
- **go get google.golang.org/grpc**/**google.golang.org/protobuf**: Installs the Go libraries for gRPC and Protobuf.

## Comparison: API Styles vs Protocols

| Name      | Type                 | Transport Protocol | Data Format | Code Generation | Description                               |
|-----------|----------------------|--------------------|-------------|-----------------|-------------------------------------------|
| REST      | API Style            | HTTP               | JSON        | None            | Simple HTTP verbs and JSON                |
| GraphQL   | Spec/Query Language  | HTTP               | JSON        | `gqlgen`        | Flexible querying, sent over HTTP         |
| SOAP      | Messaging Spec       | HTTP               | XML         | Manual/WSDL     | XML-based structure, sent over HTTP       |
| gRPC      | Protocol/Framework   | HTTP/2             | Protobuf    | `protoc`        | Binary, fast, uses protobuf, own protocol |

- **Protocols** (like HTTP, gRPC) define *how* data moves between computers.
- **Specifications/Styles** (like REST, SOAP, GraphQL) define *how* APIs structure requests/responses, but use a protocol underneath (usually HTTP).

---

## What is WSDL and .proto?

- **WSDL (Web Services Description Language):**  
  An XML file that describes a SOAP web service's operations, input/output types, and endpoint.  
  Acts as a contract for SOAP clients/servers.

- **.proto (Protocol Buffers definition):**  
  A `.proto` file defines messages and service RPC methods for gRPC.  
  Used by `protoc` to generate client/server code in many languages (Go, Python, Java, etc.).

---

## Protocol Lifecycle Visualization

### REST

1. Define handler and route.
2. Start server.
3. Client sends HTTP request (JSON).
4. Server responds with JSON.

### GraphQL

1. Define schema and resolvers.
2. Run `gqlgen generate` to update code.
3. Start server.
4. Client sends GraphQL query (JSON).
5. Server resolves and returns JSON.

### gRPC

1. Write `.proto` file with messages and services.
2. Run `protoc --go_out=. --go-grpc_out=. ...` to generate Go code.
3. Implement server methods using generated interfaces.
4. Start gRPC server on port (e.g., 50051).
5. Client (using same `.proto`) sends binary request over HTTP/2.
6. Server responds with binary data.

### SOAP

1. Write WSDL file (or update XML structs).
2. Implement handler for SOAP actions.
3. Start server.
4. Client sends XML SOAP request.
5. Server parses XML, responds with XML.

---

## Protocols and HTTP Explained

### Do REST, GraphQL, SOAP, and gRPC Use HTTP?

**Yes, they all use HTTP for transport, but with differences:**

- **REST:** Uses HTTP/1.1 (sometimes HTTP/2), sending JSON over standard HTTP methods.
- **GraphQL:** Uses HTTP/1.1 (sometimes HTTP/2 or WebSockets), sending queries/mutations as JSON.
- **SOAP:** Uses HTTP/1.1 (can use HTTP/2 or other protocols like SMTP), sending XML messages.
- **gRPC:** Uses HTTP/2 (required), sending binary encoded data (Protocol Buffers).

| API Style/Protocol | Uses HTTP?   | Typical HTTP Version |
|--------------------|--------------|---------------------|
| gRPC               | Yes (HTTP/2) | HTTP/2              |
| REST               | Yes          | HTTP/1.1 (or 2)     |
| SOAP               | Yes          | HTTP/1.1 (or 2)     |
| GraphQL            | Yes          | HTTP/1.1 (or 2)     |

---

### Why is gRPC a Protocol if It Uses HTTP?

- **gRPC** is an **application-level protocol** for remote procedure calls (RPCs).  
  It defines its own serialization (Protobuf), service definitions, streaming, and error handling.
- gRPC **uses HTTP/2 as its transport protocol**.  
  This means gRPC messages are sent over HTTP/2 connections, but gRPC defines its own message format and semantics.
- **REST, SOAP, and GraphQL** use HTTP as both their transport and application protocol (i.e., they use HTTP's verbs and structure directly).
- **gRPC** sits “on top” of HTTP/2, using its features (multiplexing, streams), but all service/method definitions and data framing are handled by gRPC itself.

**Analogy:**  
Think of HTTP/2 as a highway—REST, SOAP, and GraphQL are cars using the road’s rules directly, while gRPC is a special train running its own schedule and cargo system but traveling on the same tracks.

---

For protocol-specific code samples, see the respective subdirectories. Open issues or discussions for advanced topics or troubleshooting!