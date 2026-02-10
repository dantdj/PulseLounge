FROM node:20-alpine AS ui-build
WORKDIR /app/frontend
COPY frontend/package.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.22-alpine AS go-build
WORKDIR /app
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY --from=ui-build /app/cmd/server/web ./cmd/server/web
RUN go mod download
RUN go build -mod=readonly -o pulselounge ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=go-build /app/pulselounge ./pulselounge
EXPOSE 8080
CMD ["./pulselounge"]
