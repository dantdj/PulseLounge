# PulseLounge

Basic scaffold with:
- Go web API
- React + TypeScript UI
- UI embedded into Go binary via `go:embed`
- Docker Compose with app + PostgreSQL

Frontend release assets build into `frontend/dist` and are staged into `frontend/embed/generated` before the Go binary is built.

## API
- `GET /api/health`
- `GET /api/channels`
- `POST /api/channels`
- `DELETE /api/channels/{channelID}`
- `GET /api/channels/{channelID}/messages`
- `POST /api/channels/{channelID}/messages`
- `PUT /api/messages/{messageID}`

## Development

Development setup, local run, Docker workflows, and hot reload instructions are in [DEVELOPMENT.md](./DEVELOPMENT.md).

Database schema changes live in `db/migrations` and are applied with `goose` via `make migrate-up`.
