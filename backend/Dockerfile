FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# install air
# RUN go install github.com/air-verse/air@latest

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# COPY source code
COPY . .

# Build
RUN go build -o main .

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

# COPY --from=builder /app/.env .

EXPOSE 8080

# RUN
CMD ["./main"]
