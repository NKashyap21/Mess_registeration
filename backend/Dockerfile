# stage 1: build
FROM golang:alpine AS builder

WORKDIR /app

# Install git (needed for private/go deps sometimes)
RUN apk add --no-cache git

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build binary
RUN go build -o server ./main.go


# stage 2: run
FROM alpine:latest

WORKDIR /root/

# Copy compiled binary
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
