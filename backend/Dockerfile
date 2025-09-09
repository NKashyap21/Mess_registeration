FROM golang:alpine AS builder

RUN go install github.com/air-verse/air@latest

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main .

FROM alpine:latest


RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root

COPY --from=builder /app/main .
COPY --from=builder /go/bin/air /usr/local/bin/air


WORKDIR /app
COPY . .

EXPOSE 8080

CMD ["./main"]
