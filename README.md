# 🇧🇷 Brazilian Zipcode API

A fast and extensible API to retrieve Brazilian address information by ZIP code (CEP), with support for multiple storage backends, TTL caching, and third-party provider integration.

---

## ✨ Features

- 🔍 Fetch address data using a Brazilian CEP (zipcode)
- 🧠 Cache layer with configurable TTL using Redis
- 🗄️ Plug-and-play support for multiple storage backends (e.g. Redis, PostgreSQL, MongoDB)
- ♻️ Automatic cache invalidation and rehydration from third-party APIs (e.g. Correios)
- 🧱 Clean Architecture for scalability and maintainability
- 🚦 Built-in healthcheck endpoint
- 🔧 Configurable via `.env` using Viper

---

## 📦 Installation

```bash
> git clone https://github.com/felipefrizzo/brazilian-zipcode-api.git
> cd brazilian-zipcode-api
> go build -o zipcode-api ./cmd/*
> go run ./cmd/*
```

Server runs on :8080 by default.

## 🔧 Configuration

Create a .env file:

```bash
PORT=8080

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_USERNAME=
REDIS_PASSWORD=
REDIS_DB=0

CACHE_TTL_SECONDS=3600
```

## 📘 API Endpoints

#### Healthcheck

```bash
> curl http://localhost:8080/health
```

#### Get Address by Zipcode

```bash
> curl http://localhost:8080/zipcode/01001000

{
  "federative_unit": "PR",
  "city": "Cascavel",
  "neighborhood": "Centro",
  "address_name": "Avenida Brasil",
  "complement": "- de 5623 a 6869 - lado ímpar",
  "zipcode": "85801000",
}
```

## 🧱 Architecture

- cmd/ — Application entry point
- internal/ — Domain logic and interfaces
  - address/ — Domain model + repository interface
  - address/redis — Redis implementation of address repository
  - address/postgres — PostgreSQL implementation of address repository
  - zipcode/ — API handler
  - server/ — HTTP server initialization
  - config/ — Viper-based config loader

## 📈 Roadmap

1. Correios API integration
1. PostgreSQL and MongoDB support
1. OpenAPI documentation
1. Rate limiting and metrics

## 🤝 Contributing

PRs are welcome. Ensure code is clean, covered by tests, and follows idiomatic Go practices.
