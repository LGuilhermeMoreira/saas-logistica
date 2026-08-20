FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY proto/ ./proto/
COPY gateway/ ./gateway/

WORKDIR /app/gateway
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd

# Stage final
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/server .
RUN chmod +x /app/server

EXPOSE 8080

CMD ["/app/server"]
