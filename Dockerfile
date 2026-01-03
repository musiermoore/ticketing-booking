# Stage 1: Go modules
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dev dependencies
RUN apk add --no-cache git bash

# Copy go.mod / go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# ------------------------------
# For dev: install air for live reload
# ------------------------------
RUN go install github.com/air-verse/air@latest

EXPOSE 8080

# Command for dev
CMD ["air", "-c", ".air.toml"]
