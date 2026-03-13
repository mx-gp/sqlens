# SQLens 

SQLens is a transparent TCP proxy that intercepts and analyzes SQL queries in real-time. It helps detect **N+1 query patterns**, **slow queries**, and provides a live performance dashboard—**with zero code changes to your application.**

## Features

- **Transparent Proxy:** Just point your app to SQLens instead of your DB.
- **N+1 Detection:** Automatically flags inefficient ORM patterns.
- **Slow Query Tracking:** Real-time latency measurement and visualization.
- **Live Dashboard:** Web-based interface to see what's happening under the hood.

## Quick Start (with Docker)

1. **Clone and Start:**
   ```bash
   make docker-up
   ```

2. **Access Dashboard:**
   Open [http://localhost:8080](http://localhost:8080)

3. **Connect your App:**
   Change your DB connection port from `5432` to `5433`.
   
   *Example psql connection:*
   ```bash
   psql -h localhost -p 5433 -U user -d demo
   ```

## Development & Testing

### Prerequisites
- Go 1.23+
- Docker & Docker Compose (optional for local DB)

### Commands
- `make build`: Compile the binary.
- `make test`: Run unit tests.
- `make benchmark`: Run a realistic load simulation script.

## Configuration

SQLens can be configured via environment variables:
- `SQLENS_LISTEN_ADDR`: Proxy listen address (default `:5433`)
- `SQLENS_TARGET_ADDR`: Target database address (default `localhost:5432`)
- `SQLENS_N1_THRESHOLD`: Number of repeated queries to trigger alert (default `5`)

---
Built for SQL performance observability.
