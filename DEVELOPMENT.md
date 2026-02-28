# Development Guide

## Environment

Configuration lives in `.env` (create from `.env.example`).

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
```
