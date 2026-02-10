# PulseLounge

Basic scaffold with:
- Go web API
- React + TypeScript UI
- UI embedded into Go binary via `go:embed`
- Docker Compose with app + PostgreSQL

## Environment

Configuration lives in `.env` (create from `.env.example`).

## API
- `GET /api/health`
- `GET /api/messages`

## Local run

1. Build frontend into embedded folder:
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```
2. Run Go app:
   ```bash
   make api-run
   ```
3. Open [http://localhost:8080](http://localhost:8080)

## Docker Compose

```bash
docker compose up --build
```

- App: [http://localhost:8080](http://localhost:8080)
- Postgres: `localhost:5433` by default (`postgres/postgres`, db `pulselounge`)
- Override host port with `POSTGRES_HOST_PORT`, for example:
  ```bash
  POSTGRES_HOST_PORT=55432 docker compose up --build
  ```

If you previously ran the stack with a different Postgres setup and see role errors, recreate the database volume once:
```bash
docker compose down -v
docker compose up --build
```

## Makefile Shortcuts

```bash
make help
make up-build
make logs
make down
```
