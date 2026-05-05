# PulseLounge

Basic Discord clone using:
- Go web API
- React + TypeScript UI
- UI embedded into Go binary via `go:embed` for single-executable deployment

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

## Notes

The UI here is largely developed by AI - it's here more to provide a client for the backend than to be an example of my UI work. 