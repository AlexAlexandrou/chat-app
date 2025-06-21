FROM golang:1.24 AS builder

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY server/ ./server/

WORKDIR /app/server

# Build the Go server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/chat-server main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/chat-server .

# Port for server to listen on
EXPOSE 8080

CMD ["./chat-server"]