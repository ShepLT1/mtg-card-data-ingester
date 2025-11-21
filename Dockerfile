# card-data-ingester/Dockerfile
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY card-data-ingester/go.mod card-data-ingester/go.sum ./
RUN go mod download

# Copy the rest of the code
COPY card-data-ingester/ ./ 

# Build binary
RUN go build -o data-ingester ./cmd/server/main.go

# Final image
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/data-ingester ./

# Set environment variables if needed
# ENV DB_HOST=postgres ...

EXPOSE 8085
CMD ["./data-ingester"]
