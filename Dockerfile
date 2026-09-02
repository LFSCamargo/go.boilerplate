# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./src

FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates nodejs npm

WORKDIR /app

COPY --from=builder /server .
COPY emails ./emails
RUN cd emails && npm ci

RUN mkdir -p /app/uploads && chown -R nobody:nobody /app/uploads /app/emails

EXPOSE 5000

USER nobody:nobody

CMD ["./server"]

FROM runtime AS production

FROM golang:1.27-alpine AS development

RUN apk add --no-cache ca-certificates git nodejs npm
RUN go install github.com/air-verse/air@v1.62.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY scripts/docker/dev-entrypoint.sh /usr/local/bin/dev-entrypoint.sh
RUN chmod +x /usr/local/bin/dev-entrypoint.sh

EXPOSE 5000

ENTRYPOINT ["/usr/local/bin/dev-entrypoint.sh"]
