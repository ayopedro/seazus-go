# --- Build Stage ---
FROM golang:1.23-alpine AS builder

# Install build dependencies (make is here just in case, but we'll use go build)
RUN apk add --no-cache git make

WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# We target linux and disable CGO for a static, portable binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o seazus-api ./cmd/api

# --- Run Stage ---
FROM alpine:latest

# Add certificates for any external API calls (HTTPS)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/seazus-api .

# Ensure the binary is executable
RUN chmod +x ./seazus-api

# Railway uses the PORT env var; this is just a default
EXPOSE 8080

# Run the binary
CMD ["./seazus-api"]