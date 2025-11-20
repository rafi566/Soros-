# Architecture

Soros keeps the deployment footprint small while providing a pleasant UI inspired by Shadcn components.

## Frontend
- **Framework:** Next.js 14 with the App Router for file-based routing and metadata.
- **Styling:** Custom CSS utilities in `frontend/styles/globals.css` mimic Shadcn surface, card, and badge styling without introducing the Tailwind toolchain.
- **Components:** `InfoCard` offers a composable metric card for sources, destinations, and health summaries.

## Backend
- **Language:** Go 1.22.
- **Transport:** Standard library HTTP server with a lightweight `http.ServeMux` router.
- **Endpoints:** Health (`/health`), collections (`/sources`, `/destinations`, `/connections`), and async jobs (`/jobs`). Responses are typed via small structs in `internal/api` to keep payloads consistent.
- **Performance notes:** Request handlers wrap a 5-minute context timeout to allow larger transfers while guarding against runaway work. Sync jobs execute asynchronously with periodic progress updates so the API can frontload long-running operations without holding open client connections.

## Local development
- Run `go run ./cmd/server` from `backend/` to start the API on port 8080 (override with `PORT`).
- Run `npm run dev` from `frontend/` to launch the Next.js dashboard.
- Adjust styles or add cards by updating `frontend/styles/globals.css` and `frontend/components/info-card.tsx`.
