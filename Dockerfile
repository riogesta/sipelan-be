# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies for CGO (required by SQLite)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the rest of the source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Final Image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/main .
# Copy .env file if it exists (optional, as we can use env vars in compose)
# COPY --from=builder /app/.env . 

EXPOSE 8081

CMD ["./main"]
