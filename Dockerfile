# Build stage
FROM golang:1.26 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies (if go.sum exists)
RUN go mod download || go mod tidy

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ldag .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/ldag /app/ldag

EXPOSE 8080

ENTRYPOINT ["/app/ldag"]
