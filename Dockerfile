FROM golang:1.24.4-alpine3.22 AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o sentry ./cmd/sentry

FROM alpine:latest

COPY --from=builder /app/sentry /sentry

WORKDIR /

EXPOSE 8080

ENTRYPOINT [ "/sentry" ]