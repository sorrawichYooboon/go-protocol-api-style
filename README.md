# protocol-golang

A template project for demonstrating multi-protocol API in Go, supporting **REST**, **GraphQL**, and **SOAP**.  
Built with [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io/), [gqlgen](https://github.com/99designs/gqlgen), and native SOAP.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Initialization](#initialization)
- [Configuration](#configuration)
- [Running the Server](#running-the-server)
- [Testing the API](#testing-the-api)
  - [REST](#rest)
  - [GraphQL](#graphql)
  - [SOAP](#soap)
---

## Features

- **REST API**: Standard HTTP endpoints for CRUD operations.
- **GraphQL API**: Flexible querying via GraphQL.
- **SOAP API**: Legacy support for SOAP web services.

---

## Project Structure

```
protocol-golang/
├── config/                  # Configuration loader
├── graph/                   # GraphQL schema, models, resolvers
├── internal/
│   ├── domain/              # Domain models
│   ├── dto/                 # REST DTOs
│   ├── infrastructure/
│   │   ├── database/        # DB and repository implementations
│   │   ├── graphql/         # GraphQL route handler
│   │   ├── http/            # REST route handler
│   │   └── soap/            # SOAP handler and WSDL
│   ├── usecase/             # Business logic
│   └── repository/          # Repository interface
├── logger/                  # Logger setup
├── migrations/              # DB migrations
├── main.go                  # Entry point
├── docker-compose.yml       # PostgreSQL setup
└── README.md
```

---

## Requirements

- Go 1.20+
- Docker & Docker Compose
- [curl](https://curl.se/) (for API testing)
- [jq](https://stedolan.github.io/jq/) (optional, for pretty-printing JSON responses)

---

## Initialization

### 1. Clone the repository

```sh
git clone https://github.com/sorrawichYooboon/protocol-golang.git
cd protocol-golang
```

### 2. Start PostgreSQL with Docker Compose

```sh
docker-compose up -d
```

This starts a PostgreSQL 14 server with:
- **DB Name:** `protocoldb`
- **User:** `protocol_user`
- **Password:** `protocol_password`
- **Port:** `5432`

### 3. Install dependencies

```sh
go mod tidy
```

---

## Configuration

Set environment variables to match your database (or create a `.env` file):

```
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=protocol_user
DATABASE_PASSWORD=protocol_password
DATABASE_DBNAME=protocoldb
DATABASE_SSLMODE=disable
```

The server will look for these variables at startup.  
You can export them in your shell or use [direnv](https://direnv.net/).

---

## Running the Server

```sh
go run main.go
```

The server runs on port **8081** by default.  
Override with the `PORT` environment variable if needed.

---

## Testing the API

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

### SOAP

#### Get movie by ID

Create a file called `get_movie.xml` with this content:

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
