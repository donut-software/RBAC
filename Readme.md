# "RBAC" Microservice Example

## Introduction
This is a example of RBAC microservice that implements a collection of patterns and guidelines to create an enterprise microservice using go. 

## Domain Driven Design
This project uses a lot of the ideas introduced by Eric Evans in his book Domain Driven Design and also I've successfully created this project with the help of some online tuturials

## Project Structure
Talking specifically about microservices **only**, the structure I like to recommend is the following, everything using `<` and `>` depends on the domain being implemented and the bounded context being defined.

- [ ] `build/`: defines the code used for creating infrastructure as well as docker containers.
  - [ ] `<cloud-providers>/`: define concrete cloud provider.
  - [ ] `<executableN>/`: contains a Dockerfile used for building the binary.
- [ ] `cmd/`
  - [ ] `<primary-server>/`: uses primary database.
  - [ ] `<replica-server>/`: uses readonly databases.
  - [ ] `<binaryN>/`
- [] `db/`
  - [] `migrations/`: contains database migrations.
  - [ ] `seeds/`: contains file meant to populate basic database values.
- [ ] `internal/`: defines the _core domain_.
  - [ ] `<datastoreN>/`: a concrete _repository_ used by the domain, for example `postgresql`
  - [ ] `http/`: defines HTTP Handlers.
  - [ ] `service/`: orchestrates use cases and manages transactions.
- [] `pkg/` public API meant to be imported by other Go package.

## Tools
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1
go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.6.0
go install github.com/maxbrunsfeld/counterfeiter/v6@v6.3.0
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.8.2
go install goa.design/model/cmd/mdl@v1.7.6
go install goa.design/model/cmd/stz@v1.7.6
go install github.com/fdaines/spm-go@v0.11.1
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
```
## Features
- [X] Project Layout 
- [X] Dependency Injection 
- [X] Secure Configuration
  - [X] Using Hashicorp Vault
- [X] Metric, Traces and Logging using OpenTelemetry
- [X] Caching 
  - [X] Using memcached
- [X] Persistent Storage
  - [X] Using Postgresql
    - [X] sqlc
- [ ] Rest APIs
  - [X] HTTP Handlers
  - [X] Custom JSON Types
  - [X] Versioning 
  - [X] Error Handling
  - [X] OpenAPI 3 and Swagger-UI
  - [ ] Authorization
- [X] Tokens
  - [X] JWT
  - [X] PASETO
- [ ] Events and Messaging
  - [ ] Apache Kafka
  - [ ] RabbitMq
  - [ ] Redis
- [X] Testing
  - [X] Type safe mocks with [`maxbrunsfeld/counterfeiter`](https://github.com/maxbrunsfeld/counterfeiter)
  - [X] Equality with [`google/go-cmp`](https://github.com/google/go-cmp)
  - [X] Integration tests for Datastores with [`ory/dockertest`](https://github.com/ory/dockertest)
  - [X] REST APIs
- [ ] Containerization using Docker
- [X] Graceful Shutdown
- [X] Search Engine using [ElasticSearch](https://www.elastic.co/elasticsearch/)
- [ ] Documentation
  - [ ] [C4 Model](https://c4model.com/)
- [ ] Cloud Design Pattern
  - [ ] Reliability
    - [ ] [Circuit Breaker](https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern)
