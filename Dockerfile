# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-s -w" \
    -o /app/k4a \
    ./cmd/k4a

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN addgroup -g 1000 k4a && \
    adduser -D -u 1000 -G k4a k4a

# Set working directory
WORKDIR /home/k4a

# Copy binary from builder
COPY --from=builder /app/k4a /usr/local/bin/k4a

# Change ownership
RUN chown -R k4a:k4a /home/k4a

# Switch to non-root user
USER k4a

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/k4a"]

