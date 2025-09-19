FROM golang:1.24-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# COPY source code
COPY . .

# Build
RUN go build -o main .

# RUN
CMD ["./main"]
