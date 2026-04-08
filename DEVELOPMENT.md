# Development Guide

## Environment

Configuration lives in `.env` (create from `.env.example`).

## VSCode Workspace (Open Repo Root)

If you open the full repo in VSCode, the React TypeScript project still resolves from `frontend/`.

One-time setup:
```bash
cd frontend
npm ci
```

Notes:
- `.vscode/settings.json` points VSCode TypeScript to `frontend/node_modules/typescript/lib`.
- If you see errors in the TS in VSCode, `frontend/node_modules` is likely incomplete. Re-run `npm ci` in `frontend/`.

## Docker Hot Reload (Go + JS)

Use the dev compose file to get:
- Go auto-reload via `air`
- Vite HMR for React/TypeScript

Run:
```bash
make dev-up
```

Open:
- App (Go API + Vite UI proxy): [http://localhost:8080](http://localhost:8080)
- UI (direct Vite HMR): [http://localhost:5173](http://localhost:5173)
- API (Go, auto-reload): [http://localhost:8080/api/health](http://localhost:8080/api/health)

Notes:
- `http://localhost:8080` proxies frontend routes to Vite in dev mode.
- Frontend requests to `/api/*` are proxied by Vite to the Go container.
- This dev setup mounts your source tree into containers, so file edits reload automatically.
- Stop and remove containers with:
  ```bash
  make dev-down
  ```

## Makefile Shortcuts

```bash
make help
make dev-up
make dev-down
make seed-dev
make seed-reset-dev
make lint-go
make lint-frontend
```

## Dev Seed Data

Run after `make dev-up`:

```bash
make seed-dev
```

To wipe and reseed:

```bash
make seed-reset-dev
```
