# GeoWatch

GeoWatch is a multi-service application for real‑time wildfire awareness and geographic change analysis using satellite data. It combines a SvelteKit/CesiumJS frontend, a Go (Gin) API orchestrator, a Python (Flask) Google Earth Engine service, and a PostGIS database.

## Features
- Interactive 3D globe with CesiumJS and an optional Leaflet 2D map view
- Viewport change analysis (Before/After years) with PNG overlay rendering
- Overlay controls: opacity slider and downloadable analysis image
- Analytics dashboard (fires summary, trends, geographic metrics)
- Real‑time settings pane and toast notifications
- Docker Compose for one‑command local stack

## Architecture
- Frontend: `frontend/` (SvelteKit + CesiumJS)
  - Serves the UI, collects viewport and year inputs, calls Go API, shows overlays
- Go Backend: `backend/` (Gin)
  - Proxies analysis requests to Python service and streams PNG results
  - Exposes events API backed by PostGIS (future‑ready)
- Python GEE: `gee-service/` (Flask + Earth Engine)
  - Performs satellite change analysis and returns a rendered PNG
- Database: PostGIS

```text
frontend (SvelteKit) → go-backend (Gin) → python-gee-service (Flask/GEE)
                                     ↘  PostGIS (events)
```

## Quick Start (Docker)
1) Ensure Docker Desktop is running.

2) From repository root:
```bash
docker compose up --build
```

3) Open the app:
- Frontend: http://localhost:8082
- Go API: http://localhost:8081
- Python (direct test): http://localhost:8080

## Environment Variables
Sensitive values are provided via local env files (git‑ignored). Example templates:

- Frontend: `frontend/ENV_EXAMPLE.md`
  - `PUBLIC_GO_BACKEND_URL` – Browser URL for Go API (e.g., http://localhost:8081)
  - `PUBLIC_CESIUM_ION_TOKEN` – Your Cesium Ion token

- Backend: `backend/ENV_EXAMPLE.md`
  - `FRONTEND_ORIGIN` – CORS origin (e.g., http://localhost:8082)
  - `GO_SERVER_PORT` – Port (default 8000)
  - `PYTHON_SERVICE_URL` – Internal URL for Python service (Docker: http://python-gee-service:5000)
  - `DATABASE_URL` – Postgres/PostGIS DSN

Do not commit real `.env` files. The repo `.gitignore` excludes common env/secret paths.

## Frontend Notes
- 3D Globe (`src/lib/components/GlobalView.svelte`) provides:
  - Change Detection panel (year‑only) to analyze viewport
  - Optimal area/zoom helper popup (center top)
  - Overlay controls (opacity + download) at bottom‑left
  - MapControls (layer toggles, analytics trigger)
- 2D Map (`src/lib/components/MapContainer.svelte`) demonstrates Leaflet layers/markers

## Backend Notes
- `backend/cmd/main.go` sets up:
  - CORS via `FRONTEND_ORIGIN`
  - Port via `GO_SERVER_PORT`
  - Routes under `/api/v1`:
    - `POST /changes` – streams PNG change overlay from Python
    - `GET /events` – returns events in a bounding box (PostGIS)

## Python GEE Notes
- `gee-service/app.py` exposes `/analyze-changes` for the Go backend
- Local credentials must not be committed; see `.gitignore`

## Development (without Docker)
You can run services separately, but Docker is recommended due to Python/GEE dependencies and CORS/env alignment.

## Security & Secrets
- Never commit `.env` files or credentials; see `.gitignore`
- Use `PUBLIC_` prefix for frontend env that must be exposed to the browser (SvelteKit requirement)
- Keep GEE credentials local (`gee-service/gee-credentials.json` is git‑ignored)

## Troubleshooting
- Overlay not visible: ensure viewport is zoomed sufficiently; the popup indicates readiness. Check Python and Go logs if analysis fails.
- CORS errors: verify `FRONTEND_ORIGIN` and exposed ports in `docker-compose.yml`.
- Cesium assets: the app serves Cesium static files from `frontend/static/Cesium/`.

## License
MIT (or your preferred license). Update as needed.
