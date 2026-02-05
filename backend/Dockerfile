FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary to simplify the runtime image.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.20 AS runner

WORKDIR /app

# Create non-root user for running the service.
RUN adduser -D appuser

COPY --from=builder /app/server ./server

EXPOSE 3001

USER appuser

CMD ["./server"]
