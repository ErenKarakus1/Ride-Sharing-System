# Ride Sharing System

Microservices-based ride sharing backend built with Go.

## Backend Services

| Service | Port | Purpose |
| --- | ---: | --- |
| API Gateway | 8088 | Public REST entrypoint and JWT-protected routing |
| Auth Service | 8081 / 9092 | Registration, login, JWT issuing, token validation over gRPC |
| User Service | 8080 / 9091 | Rider and driver profiles, user gRPC API |
| Ride Service | 8082 | Ride lifecycle and Kafka ride events |
| Location Service | 8083 | Redis-backed driver location and availability |
| Matching Service | 8084 | Nearby driver matching |
| Notification Service | 8085 | Kafka ride-event consumer and authenticated WebSockets |
| Pricing Service | 8086 | Fare estimates |
| Payment Service | 8087 | Mock payment authorization, capture, and refunds |

## Stack

- Go
- Gin
- REST
- JWT
- PostgreSQL with pgx
- Redis
- WebSockets
- gRPC
- Kafka
- Docker Compose
- Prometheus
- Grafana

## Local Checks

Run all backend module tests:

```powershell
.\scripts\test-all.ps1
```

Validate Docker Compose:

```powershell
docker compose config
```

Start the backend stack after Docker Desktop is running:

```powershell
docker compose up --build
```

Run health smoke checks:

```powershell
.\scripts\smoke.ps1
```

## Observability

Every HTTP service exposes:

```text
/health
/metrics
```

Prometheus is available at:

```text
http://localhost:9090
```

Grafana is available at:

```text
http://localhost:3000
```
