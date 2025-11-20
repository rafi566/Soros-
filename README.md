# Soros

Soros is a lightweight, Airbyte-inspired data movement control plane. It pairs a Next.js experience layer with a minimal Go backend so you can explore sources, destinations, and connections without heavy infrastructure.

## Project layout
- `frontend/`: Next.js 14 app with a Shadcn-inspired dashboard.
- `backend/`: Go HTTP service exposing health, source, destination, and connection endpoints.
- `docs/`: Documentation and architecture notes.

## Backend
### Run locally
```bash
cd backend
PORT=8080 go run ./cmd/server
```

Endpoints:
- `GET /health` – service status
- `GET /sources` – example sources
- `GET /destinations` – example destinations
- `GET /connections` – example connections
- `POST /jobs` – start a long-running sync job with asynchronous progress tracking
- `GET /jobs` – list all in-flight and completed jobs
- `GET /jobs/{id}` – inspect job progress and completion timestamps

### Tests
```bash
cd backend
GOPROXY=off GOSUMDB=off go test ./...
```

## Frontend
Install dependencies and start the Next.js dev server:
```bash
cd frontend
npm install
npm run dev
```

The dashboard highlights current sources, destinations, and sync posture with compact cards and feature callouts. Styling lives in `frontend/styles/globals.css` and favors a Shadcn-inspired look without the heavier tailwind toolchain.

## Architecture notes
High-level design choices are documented in `docs/ARCHITECTURE.md`.
