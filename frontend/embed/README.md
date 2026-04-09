This directory is the Go embed staging area for the frontend bundle.

- `frontend/dist` is the Vite build output.
- `frontend/embed/generated` is copied from `frontend/dist` by `make ui-build`.
- The Go server embeds `frontend/embed/generated`.
