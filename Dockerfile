# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o ./app

# Stage 2: Runtime stage
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /build/app ./app
COPY --from=builder /build/.env .env
COPY --from=builder /build/migrations ./migrations

# Install bash or sh (if needed) in alpine
RUN apk add --no-cache bash

CMD ["/bin/sh", "-c", "/app/app migrate-up && /app/app"]
