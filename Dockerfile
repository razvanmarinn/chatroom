FROM golang:1.22 AS builder
WORKDIR /app
COPY . .

WORKDIR /app/internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o chatroom .

FROM alpine:latest
COPY --from=builder /app/internal/chatroom .
COPY --from=builder /app/internal/frontend ./frontend
RUN chmod +x chatroom
CMD ["./chatroom"]