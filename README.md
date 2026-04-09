# PulseLounge

Basic scaffold with:
- Go web API
- React + TypeScript UI
- UI embedded into Go binary via `go:embed`
- Docker Compose with app + PostgreSQL

Frontend release assets build into `frontend/dist` and are staged into `frontend/embed/generated` before the Go binary is built.

## API
- `GET /api/health`
- `GET /api/messages`
- `POST /api/messages`

## Development

Development setup, local run, Docker workflows, and hot reload instructions are in [DEVELOPMENT.md](./DEVELOPMENT.md).
