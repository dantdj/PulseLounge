FROM node:20-alpine AS ui-build
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./frontend/
RUN cd frontend && npm ci
COPY frontend/ ./frontend/
RUN cd frontend && npm run build
RUN mkdir -p /app/frontend/embed/generated \
  && cp -R /app/frontend/dist/. /app/frontend/embed/generated/

FROM golang:1.22-alpine AS go-build
WORKDIR /app
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY uiassets.go ./
COPY frontend/embed ./frontend/embed
COPY --from=ui-build /app/frontend/embed/generated ./frontend/embed/generated
RUN go mod download
RUN go build -mod=readonly -o pulselounge ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=go-build /app/pulselounge ./pulselounge
EXPOSE 8080
CMD ["./pulselounge"]
