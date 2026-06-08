# Development Guide

## Environment

Configuration lives in the `.env` file (create your own from `.env.example`).

Logging uses plain text outside production and JSON in production. Set `APP_ENV=production` or `GO_ENV=production` to enable JSON logs, or override explicitly with `LOG_FORMAT=json` or `LOG_FORMAT=text`. Set `LOG_LEVEL=debug`, `info`, `warn`, or `error` to control verbosity.

## Local VSCode TypeScript Setup

You can open the full repo root in VSCode. The React TypeScript project resolves from `frontend/`, and the repo settings point VSCode at the frontend workspace TypeScript SDK.

Install the frontend dependencies once so VSCode has the local TypeScript SDK and React types available:
```bash
cd frontend
npm ci
```

Notes:
- `.vscode/settings.json` points VSCode TypeScript to `frontend/node_modules/typescript/lib`.
- If you see TypeScript errors in VSCode after opening the repo root, `frontend/node_modules` is likely incomplete or corrupted in some way. Re-run `npm ci` in `frontend/`.

## Docker Hot Reload (Go + JS)

Use the dev compose file to get:
- Go auto-reload via `air`
- Vite HMR for React/TypeScript

Run:
```bash
make dev-up
```

The following links can then be used to access various bits:
- App (likely what you want): [http://localhost:8080](http://localhost:8080)
- UI: [http://localhost:5173](http://localhost:5173)
- API: [http://localhost:8080/api/health](http://localhost:8080/api/health)

Notes:
- `make dev-up` starts Postgres, runs pending migrations, then attaches to the API and UI dev services.
- The wait-for-database step happens inside `make migrate-up`; it retries the Postgres connection before running `goose`.
- `http://localhost:8080` proxies frontend routes to Vite in dev mode.
- Frontend requests to `/api/*` are proxied by Vite to the Go container.
- This dev setup mounts the source tree into containers, so file edits reload automatically.
- Stop and remove containers with:
  ```bash
  make dev-down
  ```

## Makefile Commands

```bash
make help
make dev-up
make dev-down
make db-reset-dev
make migrate-up
make migrate-status
make ui-build
make seed-dev
make seed-reset-dev
make lint-go
make lint-frontend
```

`make ui-build` builds the frontend into `frontend/dist` and stages the assets into `frontend/embed/generated` for `go build`.
`make migrate-up` applies the SQL files in `db/migrations` using `goose`, which records migration history in its `goose_db_version` table.
`make dev-up` already runs `make migrate-up`, so `migrate-up` is mainly for manual retries or inspecting schema setup separately.
`make db-reset-dev` removes the dev Postgres volume and restarts Postgres with a fresh empty database.

## Dev Seed Data

Run after `make dev-up`:

```bash
make seed-dev
```

To wipe and reseed:

```bash
make seed-reset-dev
```

If you need a completely fresh database before rerunning migrations:

```bash
make db-reset-dev
make dev-up
make seed-dev
```
