# --- Stage 1: Builder for backend and goose ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install goose (cached)
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.10.0

# Copy Go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy Go code
COPY . .

# --- Stage 2: Final image ---
FROM golang:1.23-alpine

WORKDIR /app

# Copy goose binary
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Install Node & pnpm once
RUN apk add --no-cache nodejs npm && npm install -g pnpm

# Copy just package.json and pnpm-lock.yaml first for caching
WORKDIR /app/client
COPY client/pnpm-lock.yaml client/package.json ./

# If `node_modules` is already there from build cache, use it
COPY --chown=node:node client/node_modules ./node_modules

# Install deps if missing (won’t re-run unless package.json/pnpm-lock.yaml changed)
RUN [ -d "node_modules" ] || pnpm install

# Copy the rest of the frontend
COPY client ./

# Build frontend (cached if unchanged)
RUN pnpm build

# Copy Go app
WORKDIR /app
COPY . .

# Start everything
CMD goose -dir ./migrations up && go run cmd/api/main.go

