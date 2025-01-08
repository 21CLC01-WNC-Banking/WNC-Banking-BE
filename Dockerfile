# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o ./app

# Stage 2: Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /build/app ./app
COPY --from=builder /build/.env .env

CMD ["/app/app"]
